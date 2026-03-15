package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// User model
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

// PaginatedResponse structure for API responses
type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

var db *sql.DB

func initDB() {
	connStr := "host=127.0.0.1 port=5433 user=postgres password=password dbname=assignment4 sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for DB to be ready
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Waiting for database... Error: %v", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Could not connect to database after several attempts")
	}
}

func getPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Filtering
	filterCol := r.URL.Query().Get("filter_col")
	filterVal := r.URL.Query().Get("filter_val")

	query := "SELECT id, name, email, gender, birth_date FROM users"
	countQuery := "SELECT COUNT(*) FROM users"
	args := []interface{}{}
	whereClauses := []string{}

	if filterCol != "" && filterVal != "" {
		// Basic validation for column name to prevent SQL injection
		validCols := map[string]bool{"id": true, "name": true, "email": true, "gender": true, "birth_date": true}
		if validCols[filterCol] {
			whereClauses = append(whereClauses, fmt.Sprintf("%s::text ILIKE $1", filterCol))
			args = append(args, "%"+filterVal+"%")
		}
	}

	if len(whereClauses) > 0 {
		where := " WHERE " + strings.Join(whereClauses, " AND ")
		query += where
		countQuery += where
	}

	// Total Count
	var totalCount int
	err := db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sorting
	orderBy := r.URL.Query().Get("order_by")
	if orderBy == "" {
		orderBy = "id"
	}
	// Simple validation for order_by
	validSortCols := map[string]bool{"id": true, "name": true, "email": true, "gender": true, "birth_date": true}
	if !validSortCols[orderBy] {
		orderBy = "id"
	}

	query += fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", orderBy, len(args)+1, len(args)+2)
	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	response := PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getCommonFriends(w http.ResponseWriter, r *http.Request) {
	u1 := r.URL.Query().Get("user_id1")
	u2 := r.URL.Query().Get("user_id2")

	if u1 == "" || u2 == "" {
		http.Error(w, "user_id1 and user_id2 are required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends f1 ON u.id = f1.friend_id
		JOIN user_friends f2 ON u.id = f2.friend_id
		WHERE f1.user_id = $1 AND f2.user_id = $2
	`

	rows, err := db.Query(query, u1, u2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var friends []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		friends = append(friends, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

func seedData() {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count > 0 {
		return
	}

	log.Println("Seeding database...")
	usernames := []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Heidi", "Ivan", "Judy", "Karl", "Linda", "Mike", "Nancy", "Oscar", "Peggy", "Quincy", "Rose", "Steve", "Trudy"}
	genders := []string{"male", "female"}

	userIDs := make([]uuid.UUID, 0)
	for i, name := range usernames {
		id := uuid.New()
		email := strings.ToLower(name) + "@example.com"
		gender := genders[rand.Intn(2)]
		birthDate := time.Now().AddDate(-20-rand.Intn(30), 0, 0)

		_, err := db.Exec("INSERT INTO users (id, name, email, gender, birth_date) VALUES ($1, $2, $3, $4, $5)", id, name, email, gender, birthDate)
		if err != nil {
			log.Printf("Error inserting user %d: %v", i, err)
			continue
		}
		userIDs = append(userIDs, id)
	}

	// Create a situation where two users have at least 3 common friends
	// Let's say User 0 and User 1 have common friends User 2, User 3, User 4
	commonFriends := []int{2, 3, 4}
	for _, friendIdx := range commonFriends {
		db.Exec("INSERT INTO user_friends (user_id, friend_id) VALUES ($1, $2)", userIDs[0], userIDs[friendIdx])
		db.Exec("INSERT INTO user_friends (user_id, friend_id) VALUES ($1, $2)", userIDs[1], userIDs[friendIdx])
	}

	// Add some random friendships
	for i := 0; i < 30; i++ {
		u1 := userIDs[rand.Intn(len(userIDs))]
		u2 := userIDs[rand.Intn(len(userIDs))]
		if u1 != u2 {
			db.Exec("INSERT INTO user_friends (user_id, friend_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", u1, u2)
		}
	}
	log.Println("Seeding complete.")
}

func main() {
	initDB()
	seedData()

	http.HandleFunc("/users", getPaginatedUsers)
	http.HandleFunc("/common-friends", getCommonFriends)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
