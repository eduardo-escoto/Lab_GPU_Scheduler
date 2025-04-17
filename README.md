# Go Web Server with MariaDB and HTMX

This project is a simple web server built with Go that interacts with a MariaDB database. It features an HTMX front end for dynamic content loading.

## Project Structure

```
go-webserver-project
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── database
│   │   └── mariadb.go   # Database connection functions
│   ├── handlers
│   │   └── routes.go    # HTTP route definitions
│   └── models
│       └── user.go      # User model and methods
├── web
│   ├── static
│   │   ├── css
│   │   │   └── styles.css # CSS styles for the front end
│   │   └── js
│   │       └── scripts.js  # JavaScript code for HTMX interactions
│   └── templates
│       ├── index.html     # Main HTML template
│       └── layout.html     # Layout template for consistent structure
├── go.mod                  # Module definition
├── go.sum                  # Dependency checksums
└── README.md               # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd go-webserver-project
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Set up the MariaDB database:**
   - Ensure you have MariaDB installed and running.
   - Create a database for the application.
   - Update the database connection details in `internal/database/mariadb.go`.

4. **Run the application:**
   ```
   go run cmd/main.go
   ```

5. **Access the application:**
   Open your web browser and navigate to `http://localhost:8080`.

## Usage Guidelines

- The application supports user interactions through HTMX for dynamic content loading.
- Modify the HTML templates in the `web/templates` directory to customize the front end.
- Extend the user model and database functions as needed in the `internal/models/user.go` and `internal/database/mariadb.go` files.

## Contributing

Feel free to submit issues or pull requests for improvements and bug fixes.