package database

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies          *MovieModel
	Users           *UserModel
	Tokens          *TokenModel
	PermissionModel *PermissionModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Movies:          &MovieModel{DB: db},
		Users:           &UserModel{DB: db},
		Tokens:          &TokenModel{DB: db},
		PermissionModel: &PermissionModel{DB: db},
	}
}

// used with listMoviesHandler
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func NewMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
