package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "io"
    "sort"
    "github.com/gorilla/mux"
)

// Define Repo structure for GitHub API response
type Repo struct {
    Name        string `json:"name"`
    HtmlUrl     string `json:"html_url"`
    Description string `json:"description"`
    Stars       int    `json:"stargazers_count"`
}

// Define the GithubClient interface
type GithubClient interface {
    GetRepos(username string) ([]Repo, error)
}

// Define RealGithubClient struct to implement the GithubClient interface
type RealGithubClient struct{}

// RealGithubClient's GetRepos method for fetching GitHub repositories
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

        body, err := io.ReadAll(resp.Body)
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

// App struct contains the GithubClient
type App struct {
    GithubClient GithubClient
}

// HomeHandler handles requests to the root page
func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
    username := "patel-aum"
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
            .about-link { text-align: center; margin-top: 20px; }
            .about-link a { font-size: 16px; color: #0366d6; text-decoration: none; }
            .about-link a:hover { text-decoration: underline; }
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
        <div class="about-link">
            <a href="/about">Learn more about Aum Patel</a>
        </div>
    </body>
    </html>
    `)
}

// AboutHandler handles requests to the /about page
func (app *App) AboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `
    <html>
    <head>
        <style>
            body { font-family: Arial, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px; }
            h1 { color: #333; text-align: center; }
            .about-section { max-width: 800px; margin: auto; background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1); }
            p { font-size: 16px; color: #555; }
        </style>
    </head>
    <body>
        <div class="about-section">
            <h1>About Aum Patel</h1>
            <p>
                I am a 2024 Computer Science graduate with Certified Kubernetes Administrator (CKA) and Certified Penetration Tester (EJPTv2) certifications.
                With 6 months of hands-on DevOps experience, I specialize in containerization, cloud platforms (AWS, Azure), cybersecurity, and CI/CD practices.
                I am proficient in Docker, Kubernetes, Terraform, and more.
            </p>
            <p>
                I have a strong foundation in infrastructure automation and security best practices, with experience working on on-premises bank servers, improving security measures, and automating security checks.
                I am passionate about combining practical knowledge and fresh perspectives in cloud-native solutions.
            </p>
            <p>
                <strong>Certifications:</strong>
                <ul>
                    <li>Certified Kubernetes Administrator (CKA) - Linux Foundation</li>
                    <li>Certified Penetration Tester (EJPTv2) - INE (eLearnSecurity)</li>
                    <li>AWS Academy Cloud Security</li>
                    <li>Database Management System - NPTEL</li>
                    <li>Salesforce Developer Virtual Internship</li>
                </ul>
            </p>
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
    r.HandleFunc("/about", app.AboutHandler).Methods("GET")

    log.Println("Server is starting on port 8081...")
    log.Fatal(http.ListenAndServe(":8081", r))
}


