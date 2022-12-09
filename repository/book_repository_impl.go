package repository

import (
	"context"
	"database/sql"
	"simple-open-library/model/domain"
)

type BookRepositoryImpl struct {
}
func (repo BookRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, book domain.Book) domain.Book {
	panic("not implemented")
}

func (repo BookRepositoryImpl) Find(ctx context.Context, tx *sql.Tx, book domain.Book) domain.Book {
	panic("not implemented")
}

