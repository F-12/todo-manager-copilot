package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	connStr := os.Getenv("POSTGRES_CONN")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=todoapp sslmode=disable"
	}

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

	r := gin.Default()

	r.GET("/api/todos", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, completed FROM todos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		defer rows.Close()
		todos := []Todo{}
		for rows.Next() {
			var t Todo
			if err := rows.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
				continue
			}
			todos = append(todos, t)
		}
		c.JSON(http.StatusOK, todos)
	})

	r.POST("/api/todos", func(c *gin.Context) {
		var input struct{
			Title string `json:"title"`
		}
		if err := c.BindJSON(&input); err != nil || input.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		var t Todo
		err := db.QueryRow("INSERT INTO todos (title) VALUES ($1) RETURNING id, title, completed", input.Title).Scan(&t.ID, &t.Title, &t.Completed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusCreated, t)
	})

	r.PATCH("/api/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		var t Todo
		err := db.QueryRow("UPDATE todos SET completed = NOT completed WHERE id = $1 RETURNING id, title, completed", id).Scan(&t.ID, &t.Title, &t.Completed)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, t)
	})

	r.DELETE("/api/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	log.Println("Gin backend running on :8080")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}
}
