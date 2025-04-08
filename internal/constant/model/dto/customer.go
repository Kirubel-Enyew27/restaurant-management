package dto

import "github.com/google/uuid"

type UpdateRequest struct {
	UserID         uuid.UUID `json:"user_id"`                   // Required field
	Username       string    `json:"username,omitempty"`        // Optional field, use *string for nullability
	Email          string    `json:"email,omitempty"`           // Optional field, use *string for nullability
	Password       string    `json:"password,omitempty"`        // Optional field, use *string for nullability
	ProfilePicture string    `json:"profile_picture,omitempty"` // Optional field, use *string for nullability

}
