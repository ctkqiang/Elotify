package structure

import (
	"time"
	"github.com/kamva/mgm/v3"
)

type Advertisement struct {
	mgm.DefaultModel `bson:",inline" swaggerignore:"true"`

	Title     *string    `json:"title" bson:"title"`
	Content   *string    `json:"content" bson:"content"`
	ImageURL  string     `json:"image_url" bson:"image_url"`
	StartTime time.Time  `json:"start_time" bson:"start_time"`
	EndTime   time.Time  `json:"end_time" bson:"end_time"`
	Locale    string     `json:"locale" bson:"locale"`
	ActionUrl *string    `json:"action_url,omitempty" bson:"action_url,omitempty"`
}