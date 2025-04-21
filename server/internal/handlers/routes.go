package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/eduardo-escoto/gpu_request/server/internal/database"
)

// HomePageData defines the structure for dynamic content passed to the template
type HomePageData struct {
	Title   string
	Heading string
	Message string
	Usage   []database.RealTimeUsage
}

// RegisterRoutesWithDB registers routes and passes the database connection to handlers
func RegisterRoutesWithDB(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/", HomeHandlerFactory(db))
	mux.HandleFunc("/gpu-usage", GPUUsageHandler(db))
	// mux.HandleFunc("/update-title", UpdateTitleHandlerFactory(db))
	// Add other handlers here, passing the db connection
}

func HomeHandlerFactory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Example query using the database connection

		rtusage, err := database.QueryRealTimeUsage(db)
		if err != nil {
			log.Fatal("Error reading query")
		}

		// log.Printf("Extracted Usage: %+v", rtusage)

		// Define the dynamic content
		data := HomePageData{
			Title:   "GPU Scheduler Home",
			Heading: "Welcome to GPU Scheduler",
			Message: fmt.Sprintf("Found %d GPUs in db.", len(rtusage)),
			Usage:   rtusage,
		}
		// Parse the template file
		tmpl, err := template.ParseFiles("web/templates/index.html")
		if err != nil {
			http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute the template with the dynamic data
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func GPUUsageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the database for real-time usage
		rtusage, err := database.QueryRealTimeUsage(db)
		if err != nil {
			http.Error(w, "Error querying GPU usage: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Define the dynamic content
		data := struct {
			Usage []database.RealTimeUsage
		}{
			Usage: rtusage,
		}

		// Parse the table template
		tmpl, err := template.New("gpu-usage").Parse(`
            <table hx-get="/gpu-usage" hx-trigger="every 5s" hx-swap="outerHTML">
                <thead>
                    <tr>
                        <th>Server Name</th>
                        <th>GPU Number</th>
                        <th>Utilization (%)</th>
                        <th>Memory Utilization (%)</th>
                        <th>Memory Used (MB)</th>
                        <th>Memory Available (MB)</th>
                        <th>Power Usage (Watts)</th>
                        <th>Temperature (°C)</th>
                        <th>Last Updated</th>
                    </tr>
                </thead>
                <tbody>
                    {{ range .Usage }}
                    <tr>
                        <td>{{ .ServerName }}</td>
                        <td>{{ .GPUNumber }}</td>
                        <td>{{ printf "%.2f" .Utilization }}</td>
                        <td>{{ printf "%.2f" .MemoryUtilization }}</td>
                        <td>{{ .MemoryUsedMB }}</td>
                        <td>{{ .MemoryAvailableMB }}</td>
                        <td>{{ printf "%.2f" .PowerUsageWatts }}</td>
                        <td>{{ printf "%.2f" .TemperatureCelsius }}</td>
                        <td>{{ .UpdatedAt.Format "2006-01-02 15:04:05" }}</td>
                    </tr>
                    {{ end }}
                </tbody>
            </table>
        `)
		if err != nil {
			http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute the template with the dynamic data
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
