package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	ips      map[string][]time.Time
	mutex    sync.Mutex
	limit    int
	window   time.Duration
	cleanup  time.Duration
	stopChan chan struct{}
}

func NewRateLimiter(limit int, window time.Duration, cleanup time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		ips:      make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		cleanup:  cleanup,
		stopChan: make(chan struct{}),
	}

	go limiter.cleanupLoop()

	return limiter
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanupExpired()
		case <-rl.stopChan:
			return
		}
	}
}

func (rl *RateLimiter) cleanupExpired() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for ip, times := range rl.ips {
		var validTimes []time.Time
		for _, t := range times {
			if now.Sub(t) < rl.window {
				validTimes = append(validTimes, t)
			}
		}
		if len(validTimes) == 0 {
			delete(rl.ips, ip)
		} else {
			rl.ips[ip] = validTimes
		}
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	times := rl.ips[ip]

	var validTimes []time.Time
	for _, t := range times {
		if now.Sub(t) < rl.window {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) >= rl.limit {
		rl.ips[ip] = validTimes
		return false
	}

	validTimes = append(validTimes, now)
	rl.ips[ip] = validTimes
	return true
}

func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

func isSQLInjectionPattern(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input))
	
	if input == "" || input == "string" || input == "en" || input == "via" {
		return false
	}

	sqlPatterns := []string{
		`\b(select|insert|update|delete|drop|truncate|create|alter|rename|grant|revoke)\s+`,
		`\b(union|exec|execute|declare|set)\s+`,
		`\b(sleep|benchmark)\s*\(`,
		`\b(or|and)\s+\d+\s*=\s*\d+`,
		`'\s*(or|and)\s+`,
		`"\s*(or|and)\s+`,
		`\b(xp_cmdshell|sysobjects|syscolumns)\b`,
		`\b(information_schema|table_schema)\b`,
		`\b(load_file|outfile|into\s+outfile)\b`,
		`--`,
		`#`,
		`\b1\s*=\s*1\b`,
	}

	for _, pattern := range sqlPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return true
		}
	}

	return false
}

func containsMaliciousKeyword(s string) bool {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(s), &jsonData); err == nil {
		return checkJSONForSQLInjection(jsonData)
	}

	return isSQLInjectionPattern(s)
}

func checkJSONForSQLInjection(data interface{}) bool {
	switch v := data.(type) {
	case string:
		return isSQLInjectionPattern(v)
	case map[string]interface{}:
		for _, value := range v {
			if checkJSONForSQLInjection(value) {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if checkJSONForSQLInjection(item) {
				return true
			}
		}
	}
	return false
}

func SecurityMiddleware() gin.HandlerFunc {
	rateLimiter := NewRateLimiter(60, time.Minute, time.Minute*5)

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)

		ip := c.ClientIP()

		// 暂时禁用限流测试功能
		if false && !rateLimiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsMaliciousKeyword(value) {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  "error",
						"message": "请求包含恶意关键词",
					})
					c.Abort()
					return
				}
			}
		}

		for name, values := range c.Request.Header {
			if name != "User-Agent" && name != "Content-Type" {
				for _, value := range values {
					if containsMaliciousKeyword(value) {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  "error",
							"message": "请求包含恶意关键词",
						})
						c.Abort()
						return
					}
				}
			}
		}

		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil {
				if containsMaliciousKeyword(string(body)) {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  "error",
						"message": "请求包含恶意关键词",
					})
					c.Abort()
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

func BlockSQLInjectionInParmametersAndBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}

		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsMaliciousKeyword(value) {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  "error",
						"message": "请求包含恶意关键词",
					})
					c.Abort()
					return
				}
			}
		}

		for name, values := range c.Request.Header {
			if name != "User-Agent" && name != "Content-Type" {
				for _, value := range values {
					if containsMaliciousKeyword(value) {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  "error",
							"message": "请求包含恶意关键词",
						})
						c.Abort()
						return
					}
				}
			}
		}

		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil {
				if containsMaliciousKeyword(string(body)) {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  "error",
						"message": "请求包含恶意关键词",
					})
					c.Abort()
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		c.Next()
	}
}
