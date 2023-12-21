package repository

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/virsavik/sample-azure-func-app/internal/adapter/repository/file"
	"github.com/virsavik/sample-azure-func-app/internal/mongotel"
)

type Registry interface {
	Files() file.Repository
}

func New(db *mongo.Database) Registry {
	return registry{
		file: file.New(mongotel.TraceCollection(db.Collection(FileCollectionName))),
	}
}

type registry struct {
	file file.Repository
}

func (r registry) Files() file.Repository {
	return r.file
}
