package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"simple-open-library/helper"
	"simple-open-library/lib/model"
)

const (
	openLibUrl = "https://openlibrary.org"
)

type OpenLibraryLib interface {
	BrowseSubjects(ctx context.Context, subject string, page int) model.OpenLibrarySubjectsResponse
}

type OpenLibraryImpl struct {
	BaseUrl string
}

func NewOpenLibraryLib() OpenLibraryLib {
	return &OpenLibraryImpl{BaseUrl: openLibUrl}
}


func (openLibrary OpenLibraryImpl) BrowseSubjects(ctx context.Context, subject string, page int) model.OpenLibrarySubjectsResponse {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/subjects/%s.json?limit=10&offset=%d",openLibrary.BaseUrl, subject, (page-1)*10),
		nil,
	)
	// request.Header.Add("")
	helper.PanicIfError(err)

	client := http.Client{}
	result, err := client.Do(request)
	helper.PanicIfError(err)
	defer result.Body.Close()

	var response model.OpenLibrarySubjectsResponse
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(&response)
	return response
}