package domain

import (
	"time"
)

type FileInfo struct {
	ID        string
	Name      string
	Metadata  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
