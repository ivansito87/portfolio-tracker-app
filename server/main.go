package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "portfolio_user"
	DB_PASSWORD = "your_secure_password"
	DB_NAME     = "portfolio_db"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

type Transaction struct {
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

var db *sql.DB

func initDB() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgreSQL!")
}

func getTransactions(c *gin.Context) {
	var transactions []Transaction
	rows, err := db.Query("SELECT id, date, amount, type FROM transactions")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transactions = append(transactions, t)
	}
	c.JSON(http.StatusOK, transactions)
}

func createTransaction(c *gin.Context) {
	var t Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO transactions (date, amount, type) VALUES ($1, $2, $3)", t.Date, t.Amount, t.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func updateTransaction(c *gin.Context) {
	id := c.Param("id")

	var t Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE transactions SET date = $1, amount = $2, type = $3 WHERE id = $4", t.Date, t.Amount, t.Type, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func deleteTransaction(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM transactions WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}

func runSQLScript() {
	cmd := exec.Command("/bin/sh", "init_db.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Script executed successfully")
}

func main() {
	initDB()
	defer db.Close()
	runSQLScript()

	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/transactions", getTransactions)
	router.POST("/transactions", createTransaction)
	router.PUT("/transactions/:id", updateTransaction)
	router.DELETE("/transactions/:id", deleteTransaction)

	fmt.Println("Server started on http://localhost:8080")
	router.Run(":8080")
}
