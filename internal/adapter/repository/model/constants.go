package model

var FileFieldName = struct {
	ID        string
	Name      string
	Metadata  string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}{
	ID:        "_id",
	Name:      "name",
	Metadata:  "metadata",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}
