package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID 
	UserID    uuid.UUID 
	Content   string    
	ImageURL  string    
	CreatedAt time.Time 
}
