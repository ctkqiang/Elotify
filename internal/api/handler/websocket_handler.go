package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"pushnotification_services/internal/structure"
)

// NewWebSocketManager 创建一个新的 WebSocket 管理器
func NewWebSocketManager() *structure.WebSocketManager {
	return &structure.WebSocketManager{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

// Run 启动 WebSocket 管理器
func Run(manager *structure.WebSocketManager) {
	for {
		select {
		case client := <-manager.Register:
			manager.Mutex.Lock()
			manager.Clients[client] = true
			manager.Mutex.Unlock()
			log.Println("新的 WebSocket 连接建立")

		case client := <-manager.Unregister:
			manager.Mutex.Lock()
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				client.Close()
				log.Println("WebSocket 连接已关闭")
			}
			manager.Mutex.Unlock()

		case message := <-manager.Broadcast:
			manager.Mutex.Lock()
			for client := range manager.Clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("发送消息失败:", err)
					client.Close()
					delete(manager.Clients, client)
				}
			}
			manager.Mutex.Unlock()
		}
	}
}

func BroadcastMessage(manager *structure.WebSocketManager, message []byte) {
	manager.Broadcast <- message
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(manager *structure.WebSocketManager, c *gin.Context) {
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("升级 WebSocket 连接失败:", err)
		return
	}

	manager.Register <- connection

	go func() {
		defer func() {
			manager.Unregister <- connection
		}()

		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket 错误: %v", err)
				}
				break
			}

			BroadcastMessage(manager, message)
		}
	}()
}
