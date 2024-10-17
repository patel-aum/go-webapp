package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "sort"
    "github.com/gorilla/mux"
)

type Repo struct {
    Name        string `json:"name"`
    HtmlUrl     string `json:"html_url"`
    Description string `json:"description"`
    Stars       int    `json:"stargazers_count"`
}

type GithubClient interface {
    GetRepos(username string) ([]Repo, error)
}

type RealGithubClient struct{}

func (c *RealGithubClient) GetRepos(username string) ([]Repo, error) {
    var allRepos []Repo
    page := 1
    perPage := 100 // Fetch 100 repos per page to minimize requests

    for {
        url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=%d", username, page, perPage)

        resp, err := http.Get(url)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return nil, fmt.Errorf("error: %s", resp.Status)
        }

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }

        var repos []Repo
        if err := json.Unmarshal(body, &repos); err != nil {
            return nil, err
        }

        // If no more repos are returned, stop the loop
        if len(repos) == 0 {
            break
        }

        // Append fetched repos to the main list
        allRepos = append(allRepos, repos...)
        page++ // Move to the next page
    }

    // Sort repos by stars in descending order
    sort.Slice(allRepos, func(i, j int) bool {
        return allRepos[i].Stars > allRepos[j].Stars
    })

    return allRepos, nil
}

type App struct {
    GithubClient GithubClient
}

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
    username := "patel-aum"  // Replace with your GitHub username
    repos, err := app.GithubClient.GetRepos(username)
    if err != nil {
        http.Error(w, "Unable to fetch GitHub repos", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, `
    <html>
    <head>
        <style>
            body { font-family: Arial, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px; }
            h1 { color: #333; text-align: center; }
            .repo-container { display: flex; flex-wrap: wrap; justify-content: center; }
            .repo-card {
                background-color: #fff;
                border: 1px solid #ddd;
                border-radius: 8px;
                padding: 16px;
                margin: 10px;
                width: 300px;
                box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
                transition: box-shadow 0.3s ease-in-out;
            }
            .repo-card:hover { box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2); }
            .repo-card h2 { font-size: 18px; margin-bottom: 10px; }
            .repo-card p { font-size: 14px; color: #555; }
            .repo-card a { color: #0366d6; text-decoration: none; }
            .repo-card a:hover { text-decoration: underline; }
            .stars { font-size: 14px; color: #666; display: flex; align-items: center; margin-bottom: 10px; }
            .stars svg { margin-right: 5px; }
        </style>
    </head>
    <body>
        <h1>Aum's GitHub Projects</h1>
        <div class="repo-container">
    `)

    for _, repo := range repos {
        fmt.Fprintf(w, `
            <div class="repo-card">
                <h2><a href="%s">%s</a></h2>
                <div class="stars">
                    <svg height="16" width="16" viewBox="0 0 16 16" aria-hidden="true"><path fill="#666" d="M8 .25a.75.75 0 01.673.418l1.86 3.766 4.153.603a.75.75 0 01.416 1.28l-3.003 2.927.709 4.137a.75.75 0 01-1.088.791L8 12.347l-3.71 1.95a.75.75 0 01-1.088-.79l.709-4.137L.907 6.317a.75.75 0 01.416-1.28l4.153-.603L7.327.668A.75.75 0 018 .25z"></path></svg>
                    %d stars
                </div>
                <p>%s</p>
            </div>
        `, repo.HtmlUrl, repo.Name, repo.Stars, repo.Description)
    }

    fmt.Fprintf(w, `
        </div>
    </body>
    </html>
    `)
}

func main() {
    app := &App{
        GithubClient: &RealGithubClient{},
    }

    r := mux.NewRouter()
    r.HandleFunc("/", app.HomeHandler).Methods("GET")

    log.Println("Server is starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}
