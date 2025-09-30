package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("POSTGRES_CONN")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=todoapp sslmode=disable"
	}
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE
	)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			rows, err := db.Query("SELECT id, title, completed FROM todos")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			defer rows.Close()
			todos := []map[string]interface{}{}
			for rows.Next() {
				var id int
				var title string
				var completed bool
				if err := rows.Scan(&id, &title, &completed); err != nil {
					continue
				}
				todos = append(todos, map[string]interface{}{
					"id":        id,
					"title":     title,
					"completed": completed,
				})
			}
			json, _ := toJSON(todos)
			w.Write(json)
		case http.MethodPost:
			var input struct {
				Title string `json:"title"`
			}
			if err := fromJSON(r, &input); err != nil || input.Title == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid input"}`))
				return
			}
			row := db.QueryRow("INSERT INTO todos (title) VALUES ($1) RETURNING id, title, completed", input.Title)
			var id int
			var title string
			var completed bool
			if err := row.Scan(&id, &title, &completed); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			json, _ := toJSON(map[string]interface{}{"id": id, "title": title, "completed": completed})
			w.WriteHeader(http.StatusCreated)
			w.Write(json)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/todos/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idStr := r.URL.Path[len("/api/todos/"):] // /api/todos/{id}
		var id int
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid id"}`))
			return
		}
		switch r.Method {
		case http.MethodPatch:
			row := db.QueryRow("UPDATE todos SET completed = NOT completed WHERE id = $1 RETURNING id, title, completed", id)
			var nid int
			var title string
			var completed bool
			if err := row.Scan(&nid, &title, &completed); err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"not found"}`))
				return
			}
			json, _ := toJSON(map[string]interface{}{"id": nid, "title": title, "completed": completed})
			w.Write(json)
		case http.MethodDelete:
			_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Go backend running on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

}

// 工具函数：JSON序列化
func toJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// 工具函数：JSON反序列化
func fromJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
