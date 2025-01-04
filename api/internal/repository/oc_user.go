package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"api/pkg/ocserv"
	"api/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type OcservUserRepository struct {
	db *gorm.DB
	oc *ocserv.Handler
}

type OcservUserRepositoryInterface interface {
	Users(context.Context, utils.RequestPagination) (*[]models.OcUser, *utils.ResponsePagination, error)
}

func NewOcservUserRepository() *OcservUserRepository {
	return &OcservUserRepository{
		db: database.Connection(),
		oc: ocserv.NewHandler(),
	}
}

func (o *OcservUserRepository) Users(ctx context.Context, page utils.RequestPagination) (
	*[]models.OcUser, *utils.ResponsePagination, error,
) {

	return nil, nil, nil
}
