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
	Users       *UserModel
	Tokens      *TokenModel
	Permissions *PermissionModel
	Movies      *MovieModel
	People      *PeopleModel
	Categories  *CategoriesModel
	Casts       *CastsModel
	ImageIDs    *ImageIDsModel
	MovieLinks  *MovieLinkModel
	PeopleLinks *PeopleLinkModel
	Trailer     *TrailersModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:       &UserModel{DB: db},
		Tokens:      &TokenModel{DB: db},
		Permissions: &PermissionModel{DB: db},
		Movies:      &MovieModel{DB: db},
		People:      &PeopleModel{DB: db},
		Categories:  &CategoriesModel{DB: db},
		Casts:       &CastsModel{DB: db},
		ImageIDs:    &ImageIDsModel{DB: db},
		MovieLinks:  &MovieLinkModel{DB: db},
		PeopleLinks: &PeopleLinkModel{DB: db},
		Trailer:     &TrailersModel{DB: db},
	}
}
