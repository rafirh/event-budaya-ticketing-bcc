package dto

import (
	"time"

	"event-budaya-ticketing-bcc/internal/model"

	"github.com/google/uuid"
)

type EventCommentUserSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	ProfilePhoto *string   `json:"profile_photo,omitempty"`
}

type EventCommentResponse struct {
	ID        uuid.UUID               `json:"id"`
	EventID   uuid.UUID               `json:"event_id"`
	User      EventCommentUserSummary `json:"user"`
	ParentID  *uuid.UUID              `json:"parent_id"`
	Comment   string                  `json:"comment"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt *time.Time              `json:"updated_at"`
	Replies   []EventCommentResponse  `json:"replies"`
}

func ToEventCommentResponse(comment model.EventComment) EventCommentResponse {
	resp := EventCommentResponse{
		ID:        comment.ID,
		EventID:   comment.EventID,
		User:      EventCommentUserSummary{ID: comment.User.ID, Name: comment.User.Name, ProfilePhoto: comment.User.ProfilePhoto},
		ParentID:  comment.ParentID,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Replies:   make([]EventCommentResponse, 0, len(comment.Replies)),
	}

	for _, reply := range comment.Replies {
		resp.Replies = append(resp.Replies, ToEventCommentResponse(reply))
	}

	return resp
}
