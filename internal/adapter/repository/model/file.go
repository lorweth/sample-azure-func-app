package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileInfo struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Metadata  string             `bson:"metadata"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at"`
}
