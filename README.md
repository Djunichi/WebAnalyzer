# WebAnalyzer

A web application written in Go that analyzes webpages by URL and provides structured information about the HTML document. This project was developed as part of a technical assignment.

---

## Features

- Analyze any publicly accessible webpage by URL
- Extract the following metadata:
    - HTML version (based on the `<!DOCTYPE>`)
    - Page `<title>`
    - Count of heading tags (`h1` through `h6`)
    - Number of internal, external, and inaccessible links
    - Presence of a login form (`<input type="password">`)
- View history of all past analyzed pages
- Reopen any previous analysis by ID

---

##  Technologies Used

- **Backend:** Golang (Gin framework)
- **Frontend:** HTML + JavaScript
- **Database:** PostgreSQL (via GORM)
- **Containerization:** Docker & docker-compose
- **HTML Parsing:** `golang.org/x/net/html`

---

## How to Run

### Prerequisites
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### ▶ Run the application

```bash
docker-compose up --build
```

### Then visit: http://localhost:8080