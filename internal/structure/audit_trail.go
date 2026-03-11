package structure

type AuditTrail struct {
	PushedBy string  `json:"pushed_by"`
	PushedAt string  `json:"pushed_at"`
	Via      *string `json:"via"`
}
