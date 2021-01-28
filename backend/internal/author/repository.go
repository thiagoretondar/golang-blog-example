package author

import (
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/database"
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/database/postgres"
)

type Repository interface {
	database.CRUDRepository
}

type repo struct {
	*postgres.PgTxRepository
}

func NewRepository() Repository {
	return &repo{}
}
