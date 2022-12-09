package repository

import (
	"context"
	"simple-open-library/model/domain"
)

type LibraryRepository interface {
	Subjects(ctx context.Context, subject string, page int) ([]domain.LibraryBook, error)
}