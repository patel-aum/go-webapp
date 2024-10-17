package main

import (
	"net/http"
	"net/http/httptest"
	"strings"  
	"testing"
)

type MockGithubClient struct{}
func (c *MockGithubClient) GetRepos(username string) ([]Repo, error) {
	return []Repo{
		{Name: "TestRepo1", HtmlUrl: "https://github.com/user/testrepo1", Description: "Test Description 1", Stars: 10},
		{Name: "TestRepo2", HtmlUrl: "https://github.com/user/testrepo2", Description: "Test Description 2", Stars: 5},
	}, nil
}

func TestHomeHandler(t *testing.T) {
	app := &App{
		GithubClient: &MockGithubClient{},
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.HomeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expectedContentType := "text/html; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v, want %v", contentType, expectedContentType)
	}
	expectedBodySnippet := "Aum's GitHub Projects"
	if !strings.Contains(rr.Body.String(), expectedBodySnippet) {
		t.Errorf("handler returned unexpected body: expected to contain %q", expectedBodySnippet)
	}
}

