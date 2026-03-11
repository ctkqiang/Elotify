package config


var (
	INDEX = "/"
	WEBSCOKET_CHANNEL = "/ws/announcement"

	SWAGGER_DOCS = "/swagger/*any"
	ENV = "/.env"

	HEALTH = "/health"
)

var (
	SEND_TEXT_PUSH_NOTIFICATION_HEAD      = "/push"
	SEND_TEXT_PUSH_NOTIFICATION           = "/text"
	SEND_TEXT_AND_IMAGE_PUSH_NOTIFICATION = "/text-image"
)

var (
	SEGMENTATION     = "/segment"
	SEGMENT_LIST_ALL = "/all"
	SEGMENT_CREATE   = "/create"
	SEGMENT_DELETE   = "/delete"
	SEGMENT_UPDATE   = "/update"
)

var (
	ANNOUNCEMENT        = "/announcement"
	ANNOUNCEMENT_CREATE = "/create"
	ANNOUNCEMENT_DELETE = "/delete"
	ANNOUNCEMENT_LATEST   = "/latest"
	ANNOUNCEMENT_UPDATE = "/update"
	ANNOUNCEMENT_LIST_ALL = "/all"
)

var (
	ADVERTISEMENT        = "/advertisement"
	ADVERTISEMENT_CREATE = "/create"
	ADVERTISEMENT_DELETE = "/delete"
	ADVERTISEMENT_LATEST   = "/latest"
	ADVERTISEMENT_UPDATE = "/update"
	ADVERTISEMENT_LIST_ALL = "/all"
)