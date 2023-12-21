package file

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/virsavik/sample-azure-func-app/internal/adapter/repository/model"
	"github.com/virsavik/sample-azure-func-app/internal/core/domain"
	"github.com/virsavik/sample-azure-func-app/internal/mongotel"
)

type Repository struct {
	db mongotel.CollectionOperations
}

func New(client mongotel.CollectionOperations) Repository {
	return Repository{
		db: client,
	}
}

func (r Repository) findByID(ctx context.Context, id string) (model.FileInfo, error) {
	modelID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.FileInfo{}, pkgerrors.WithStack(err)
	}

	// Find documents by ID
	rs := r.db.FindOne(
		ctx,
		bson.D{{"_id", modelID}},
		nil,
	)
	if err := rs.Err(); err != nil {
		// Not error, just no document found
		if errors.Is(rs.Err(), mongo.ErrNoDocuments) {
			return model.FileInfo{}, nil
		}

		return model.FileInfo{}, pkgerrors.WithStack(err)
	}

	// Convert to model
	var m model.FileInfo
	if err := rs.Decode(&m); err != nil {
		return model.FileInfo{}, pkgerrors.WithStack(err)
	}

	return m, nil
}

func (r Repository) Save(ctx context.Context, info domain.FileInfo) error {
	var m model.FileInfo
	if info.ID != "" {
		var err error
		m, err = r.findByID(ctx, info.ID)
		if err != nil {
			return err
		}
	}

	// Update if document exists
	if !m.ID.IsZero() {
		// Prepare new field information
		var newData bson.D
		if info.Name != m.Name {
			newData = append(newData, bson.E{Key: model.FileFieldName.Name, Value: info.Name})
		}
		if info.Metadata != m.Metadata {
			newData = append(newData, bson.E{Key: model.FileFieldName.Metadata, Value: info.Metadata})
		}
		if info.DeletedAt != m.DeletedAt {
			newData = append(newData, bson.E{Key: model.FileFieldName.DeletedAt, Value: info.DeletedAt})
		}

		_, err := r.db.UpdateOne(
			ctx,
			bson.D{{model.FileFieldName.ID, info.ID}},
			newData,
		)
		if err != nil {
			return pkgerrors.WithStack(err)
		}

		return nil
	}

	// Insert new info
	m = model.FileInfo{
		ID:       primitive.NewObjectID(),
		Name:     info.Name,
		Metadata: info.Metadata,
	}

	_, err := r.db.InsertOne(
		ctx,
		m,
	)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}
