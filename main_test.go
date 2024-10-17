package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Repo represents a GitHub repository
type Repo struct {
	Name        string
	HtmlUrl     string
	Description string
	Stars       int
}

// GithubClient interface for fetching GitHub repos
type GithubClient interface {
	GetRepos(username string) ([]Repo, error)
}

// App represents the application with a GitHub client
type App struct {
	GithubClient GithubClient
}

// MockGithubClient simulates fetching repositories
type MockGithubClient struct{}

func (c *MockGithubClient) GetRepos(username string) ([]Repo, error) {
	return []Repo{
		{Name: "TestRepo1", HtmlUrl: "https://github.com/user/testrepo1", Description: "Test Description 1", Stars: 10},
		{Name: "TestRepo2", HtmlUrl: "https://github.com/user/testrepo2", Description: "Test Description 2", Stars: 5},
	}, nil
}

// HomeHandler handles the home page request
func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type to HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Simulating a response with the list of repositories
	repos, err := app.GithubClient.GetRepos("someuser")
	if err != nil {
		http.Error(w, "Unable to fetch repos", http.StatusInternalServerError)
		return
	}

	// For simplicity, only returning the titles of repositories
	for _, repo := range repos {
		w.Write([]byte(repo.Name + "<br>")) // Using <br> for line breaks in HTML
	}
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
	expectedBodySnippet := "TestRepo1"
	if !strings.Contains(rr.Body.String(), expectedBodySnippet) {
		t.Errorf("handler returned unexpected body: expected to contain %q", expectedBodySnippet)
	}
}

