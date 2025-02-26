package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// const (
// 	DB_USER     = "portfolio_user"
// 	DB_PASSWORD = "your_secure_password"
// 	DB_NAME     = "portfolio_db"
// 	DB_HOST     = "localhost"
// 	DB_PORT     = "5432"
// )

const (
	DB_USER     = "portfolio_user"
	DB_PASSWORD = "your_secure_password"
	DB_HOST     = "database-2.crkkai2skkf4.us-east-2.rds.amazonaws.com"
	DB_PORT     = "5432"
	DB_NAME     = "portfolio_db"
)

type Transaction struct {
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

var (
	db *sql.DB
	// Channel for database operations
	dbChan = make(chan dbOperation, 100)
	// WaitGroup for graceful shutdown
	wg sync.WaitGroup
)

// Define database operation types
type dbOperation struct {
	opType     string
	data       interface{}
	resultChan chan dbResult
}

type dbResult struct {
	data interface{}
	err  error
}

func initDB() {
	var err error
	maxRetries := 5

	// Use a simpler connection string first
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)

	log.Printf("Attempting to connect to database at %s", DB_HOST)

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Failed to open database: %v\n", i+1, err)
			time.Sleep(time.Second * 5)
			continue
		}

		// Configure connection pool with more conservative settings
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(1 * time.Minute)

		// Test connection with longer timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err = db.PingContext(ctx)
		cancel()

		if err != nil {
			log.Printf("Attempt %d: Failed to ping database: %v\n", i+1, err)
			db.Close()
			time.Sleep(time.Second * 5)
			continue
		}

		log.Printf("Successfully connected to PostgreSQL at %s!", DB_HOST)
		return
	}
	log.Fatal("Failed to connect to database after maximum retries:", err)
}

func dbWorker() {
	defer wg.Done()
	for op := range dbChan {
		start := time.Now()
		log.Printf("Starting %s operation", op.opType)

		switch op.opType {
		case "getTransactions":
			transactions, err := executeGetTransactions()
			if err != nil {
				log.Printf("Error in getTransactions: %v", err)
			}
			op.resultChan <- dbResult{data: transactions, err: err}
		case "createTransaction":
			t := op.data.(Transaction)
			err := executeCreateTransaction(t)
			op.resultChan <- dbResult{data: t, err: err}
		case "updateTransaction":
			t := op.data.(map[string]interface{})
			err := executeUpdateTransaction(t["id"].(string), t["transaction"].(Transaction))
			op.resultChan <- dbResult{data: t["transaction"], err: err}
		case "deleteTransaction":
			id := op.data.(string)
			err := executeDeleteTransaction(id)
			op.resultChan <- dbResult{err: err}
		}

		log.Printf("Completed %s operation in %v", op.opType, time.Since(start))
	}
}

func executeGetTransactions() ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT id, date, amount, type FROM transactions ORDER BY date DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.Type); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return transactions, nil
}

func executeCreateTransaction(t Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `INSERT INTO transactions (date, amount, type) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRowContext(ctx, query, t.Date, t.Amount, t.Type).Scan(&t.ID)
	if err != nil {
		return fmt.Errorf("insert error: %v", err)
	}
	return nil
}

func executeUpdateTransaction(id string, t Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `UPDATE transactions SET date = $1, amount = $2, type = $3 WHERE id = $4`
	result, err := db.ExecContext(ctx, query, t.Date, t.Amount, t.Type, id)
	if err != nil {
		return fmt.Errorf("update error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("no transaction found with id %s", id)
	}
	return nil
}

func executeDeleteTransaction(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `DELETE FROM transactions WHERE id = $1`
	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("no transaction found with id %s", id)
	}
	return nil
}

// Updated handler functions
func getTransactions(c *gin.Context) {
	resultChan := make(chan dbResult)
	dbChan <- dbOperation{
		opType:     "getTransactions",
		resultChan: resultChan,
	}

	result := <-resultChan
	if result.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.data)
}

func createTransaction(c *gin.Context) {
	var t Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultChan := make(chan dbResult)
	dbChan <- dbOperation{
		opType:     "createTransaction",
		data:       t,
		resultChan: resultChan,
	}

	result := <-resultChan
	if result.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.err.Error()})
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

	resultChan := make(chan dbResult)
	dbChan <- dbOperation{
		opType: "updateTransaction",
		data: map[string]interface{}{
			"id":          id,
			"transaction": t,
		},
		resultChan: resultChan,
	}

	result := <-resultChan
	if result.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func deleteTransaction(c *gin.Context) {
	id := c.Param("id")

	resultChan := make(chan dbResult)
	dbChan <- dbOperation{
		opType:     "deleteTransaction",
		data:       id,
		resultChan: resultChan,
	}

	result := <-resultChan
	if result.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}

func main() {
	log.Printf("Starting server...")

	// Initialize database connection
	initDB()

	// Start the database worker
	wg.Add(1)
	go dbWorker()

	// Initialize router with release mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/transactions", getTransactions)
		api.POST("/transactions", createTransaction)
		api.PUT("/transactions/:id", updateTransaction)
		api.DELETE("/transactions/:id", deleteTransaction)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.PingContext(ctx)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database connection failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Server configuration
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown setup
	go func() {
		log.Printf("Server starting on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	close(dbChan)
	wg.Wait()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if db != nil {
		db.Close()
	}

	log.Println("Server exited properly")
}
