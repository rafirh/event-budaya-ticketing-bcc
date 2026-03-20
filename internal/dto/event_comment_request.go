package dto

type CreateEventCommentRequest struct {
	Comment  string  `json:"comment" validate:"required"`
	ParentID *string `json:"parent_id" validate:"omitempty,uuid"`
}
