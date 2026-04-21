package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestStatus represents the state of an idempotent request
type RequestStatus string

const (
	StatusProcessing RequestStatus = "processing"
	StatusCompleted  RequestStatus = "completed"
)

type IdempotencyRecord struct {
	Status   RequestStatus
	Response gin.H
	Code     int
}

// In-memory storage for idempotency
var (
	idempotencyStore = make(map[string]*IdempotencyRecord)
	storeMu          sync.Mutex
)

// IdempotencyMiddleware intercepts requests and checks the Idempotency-Key header.
func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("Idempotency-Key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header is missing"})
			c.Abort()
			return
		}

		storeMu.Lock()
		record, exists := idempotencyStore[key]

		if exists {
			if record.Status == StatusProcessing {
				storeMu.Unlock()
				c.JSON(http.StatusConflict, gin.H{"error": "Request with this key is already in progress"})
				c.Abort()
				return
			}
			if record.Status == StatusCompleted {
				storeMu.Unlock()
				// Return saved result without calling business logic
				fmt.Println("[Middleware] Returning cached response for key:", key)
				c.JSON(record.Code, record.Response)
				c.Abort()
				return
			}
		}

		// New request: mark as processing
		idempotencyStore[key] = &IdempotencyRecord{Status: StatusProcessing}
		storeMu.Unlock()

		// Allow execution to continue
		c.Next()
	}
}

// PaymentHandler simulates heavy operation and returns JSON.
func PaymentHandler(c *gin.Context) {
	fmt.Println("Processing started...")
	time.Sleep(2 * time.Second) // Simulating heavy operation (2s)

	key := c.GetHeader("Idempotency-Key")
	
	response := gin.H{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": uuid.New().String(),
	}

	// Update record to completed
	storeMu.Lock()
	idempotencyStore[key] = &IdempotencyRecord{
		Status:   StatusCompleted,
		Response: response,
		Code:     http.StatusOK,
	}
	storeMu.Unlock()

	fmt.Println("Processing finished for key:", key)
	c.JSON(http.StatusOK, response)
}

func main() {
	// Set gin to release mode for cleaner logs
	gin.SetMode(gin.ReleaseMode)

	r := gin.New() // New() instead of Default() to avoid standard logs which might clutter the simulation results
	r.Use(gin.Recovery())
	
	// Logger with custom format to show our specific logs clearly
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[Gin] %s | %d | %s\n",
			param.Path,
			param.StatusCode,
			param.ErrorMessage,
		)
	}))

	r.Use(IdempotencyMiddleware())
	r.POST("/pay", PaymentHandler)

	// Simulation of "Double-Click" Attack
	go func() {
		time.Sleep(1 * time.Second) // Give server time to start
		
		key := uuid.New().String()
		fmt.Printf("\n=== Simulating Double-Click Attack with Key: %s ===\n", key)
		
		var wg sync.WaitGroup
		// Send 5 requests almost simultaneously
		for i := 1; i <= 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				client := &http.Client{}
				req, _ := http.NewRequest("POST", "http://localhost:8080/pay", nil)
				req.Header.Set("Idempotency-Key", key)
				
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Request %d error: %v\n", id, err)
					return
				}
				defer resp.Body.Close()
				fmt.Printf("Request %d: Status Code %d\n", id, resp.StatusCode)
			}(i)
			time.Sleep(100 * time.Millisecond) // Slight delay to separate them in logs, but all within 2s window
		}
		wg.Wait()
		
		// Final request after first one should have finished (it takes 2s)
		fmt.Println("\nWaiting for initial request to complete...")
		time.Sleep(1500 * time.Millisecond) // Remaining time to ensure the first one finishes (total ~2s since start)

		fmt.Println("=== Final Request after first completion ===")
		req, _ := http.NewRequest("POST", "http://localhost:8080/pay", nil)
		req.Header.Set("Idempotency-Key", key)
		resp, _ := http.DefaultClient.Do(req)
		
		fmt.Printf("Final Request: Status Code %d\n", resp.StatusCode)
		
		// Exit after simulation is done
		time.Sleep(500 * time.Millisecond)
		fmt.Println("\n=== Idempotency Simulation Complete ===")
		// In a real server we wouldn't exit, but for this task we do.
		// However, r.Run() is blocking. We can send a signal or just exit.
		// For the sake of the exercise, we can just kill the process here or let the user terminate.
		// Let's just wait a bit and exit.
		// (Actually, r.Run is blocking. I'll just print instructions to user).
	}()

	fmt.Println("Server starting on :8080...")
	r.Run(":8080")
}
