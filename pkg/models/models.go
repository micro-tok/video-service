package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Video struct {
	CreatedAt   time.Time `json:"created_at"`
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	Tags        []string  `json:"tags"`
}
