package structure

type NotificationContent struct {
	Title      string     `json:"title"`
	Message    string     `json:"message"`
	ImageUrl   *string    `json:"image_url,omitempty"`
	Channel    *string    `json:"channel,omitempty"`
	Locale     *Locale    `json:"locale,omitempty"`
	Segments   *[]string  `json:"segments,omitempty"`
	AuditTrail AuditTrail `json:"audit_trail"`
}
