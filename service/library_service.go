package service

import (
	"context"
	"simple-open-library/model/web"
)

type LibraryService interface {
	BrowseBySubject(ctx context.Context, request web.SubjectRequest) web.SubjectResponse
}