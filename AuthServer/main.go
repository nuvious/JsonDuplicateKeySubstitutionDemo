package main

import (
    "database/sql"
    "log"
    "net/http"
    "io/ioutil"
    "github.com/buger/jsonparser"
    _ "github.com/mattn/go-sqlite3"
)

const databaseFile = "users.db"
var db *sql.DB

func main() {
	// Initialize the db
    var err error
    db, err = sql.Open("sqlite3", databaseFile)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()
    initDatabase()

	// Setup auth and registration handlers
    http.HandleFunc("/auth", authHandler)
    http.HandleFunc("/register", registerHandler)
    http.ListenAndServe(":8080", nil)
}

func initDatabase() {
    // Initialize a basic username and password database
    createTableSQL := `CREATE TABLE IF NOT EXISTS users (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "username" TEXT UNIQUE,
        "password" TEXT
    );`
    _, err := db.Exec(createTableSQL)
    if err != nil {
        log.Fatalf("Failed to create table: %v", err)
    }
}

func authHandler(w http.ResponseWriter, r *http.Request) {
    var err error

    // Ensure it's a POST request
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST is supported", http.StatusBadRequest)
        return
    }

    // Ensure it has a request body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    // Get the username and password
    username, _, _, err := jsonparser.Get(body, "username")
    if err != nil {
        http.Error(w, "Failed to parse username", http.StatusBadRequest)
        return
    }
    password, _, _, err := jsonparser.Get(body, "password")
    if err != nil {
        http.Error(w, "Failed to parse password", http.StatusBadRequest)
        return
    }
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Don't let admin login through the front UI
    if string(username) == "admin" {
        http.Error(w, "Admin user cannot authenticate through this endpoint",
                   http.StatusForbidden)
        return
    }

    // Check the password and if it's correct, return the original json,
    // otherwise return an empty object
    isValid := checkPassword(string(username), string(password))
    if isValid {
        w.Write(body)
    } else {
        w.Write([]byte("{}"))
    }
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    var err error

    // Ensure it's a POST request
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST is supported", http.StatusBadRequest)
        return
    }

    // Ensure it has a request body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
	// Get username and password
    username, _, _, err := jsonparser.Get(body, "username")
    if err != nil {
        http.Error(w, "Failed to parse username", http.StatusBadRequest)
        return
    }
    password, _, _, err := jsonparser.Get(body, "password")
    if err != nil {
        http.Error(w, "Failed to parse password", http.StatusBadRequest)
        return
    }
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

	// Add the new user in a transaction
    tx, err := db.Begin()
    if err != nil {
        log.Fatalf("Failed to begin transaction: %v", err)
        return
    }
    _, err = tx.Exec("INSERT INTO users(username, password) VALUES(?, ?)",
					 username, password)
    if err != nil {
        tx.Rollback()
        http.Error(w, "Failed to register", http.StatusInternalServerError)
        return
    }
    err = tx.Commit()
    if err != nil {
        log.Fatalf("Failed to commit transaction: %v", err)
        return
    }

    w.Write([]byte("Registered successfully"))
}

func checkPassword(username, password string) bool {
    rows, err := db.Query("SELECT id, username, password FROM users")
    if err != nil {
        log.Fatalf("Failed to query users: %v", err)
        return false
    }
    defer rows.Close()

    log.Println("Users in database:")
    for rows.Next() {
        var id int
        var storedUsername, storedPassword string

        err := rows.Scan(&id, &storedUsername, &storedPassword)
        if err != nil {
            log.Printf("Failed to scan user row: %v", err)
            continue
        }
        log.Printf("ID: %d, Username: %s, Password: %s",
		           id, storedUsername, storedPassword)

        if username == storedUsername && password == storedPassword {
            return true
        }
    }

    if err := rows.Err(); err != nil {
        log.Printf("Row reading error: %v", err)
    }

    return false
}

