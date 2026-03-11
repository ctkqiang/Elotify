package structure

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type AnnouncementType string

var (
	ANNOUNCEMENT_TYPE_HOLIDAY = AnnouncementType("HOLIDAY")
	ANNOUNCEMENT_TYPE_EVENT   = AnnouncementType("EVENT")
)

type Announcement struct {
	mgm.DefaultModel `bson:",inline"`

	ID        string           `json:"id" bson:"_id"`
	Type      AnnouncementType `json:"type" bson:"type"`
	Message   string           `json:"message" bson:"message"`
	Priority  Priority         `json:"priority" bson:"priority"`
	CreatedAt time.Time        `json:"created_at" bson:"created_at"`
	StartedAt time.Time        `json:"started_at" bson:"started_at"`
	ExpiresAt time.Time        `json:"expires_at" bson:"expires_at"`
}
