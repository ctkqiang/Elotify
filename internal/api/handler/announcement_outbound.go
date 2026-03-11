package handler

import (
	"pushnotification_services/internal/repositories"
	"pushnotification_services/internal/structure"
	"time"

	"github.com/google/uuid"
)

func BroadCastAnnouncementToAllApp(message string, priority structure.Priority) {
	announcement := structure.Announcement{
		ID:        uuid.New().String(),
		Type:      structure.ANNOUNCEMENT_TYPE_EVENT,
		Message:   message,
		Priority:  priority,
		CreatedAt: time.Now(),
	}

	err := repositories.WriteAnnouncement(announcement)
	if err != nil {
		return
	}
}
