// Package db provides database connection and access to the Event data
package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/cbarlas/PeakAPI/model"
	"fmt"
	"time"
	"math/rand"
)

// Constant variables for creating and connecting to database instance
const (
	DB_DRIVER = "mysql"
	DB_USERNAME = "root"
	DB_PW = "root"
	CONNECTION = DB_USERNAME + ":" + DB_PW + "@tcp(127.0.0.1:3306)/"
	DB_NAME = "api_db"
	RESPONSE_TABLE = "responses"
)

// Prints the error to the console if it is not empty
func ErrorCheck(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Returns the Event list with the same API_KEY values
func GetEventsWithApiKey(s string) (evs model.Events) {
	conn, err := sql.Open(DB_DRIVER, CONNECTION)
	ErrorCheck(err)
	defer conn.Close()

	evs = model.Events{}
	CreateAndUseDB(DB_NAME, conn)
	rows, err := conn.Query("SELECT * FROM " + s)

	if err != nil {
		return
	}

	for rows.Next() {
		var apiKey string
		var userID uint64
		var timestamp int64
		
		err = rows.Scan(&apiKey, &userID, &timestamp)
		ErrorCheck(err)
		evs = append(evs, &model.Event{apiKey, userID, timestamp})
	} 
	
	return
}

// Returns the response times list
func GetResponseTimes() (responses []int) {
	conn, err := sql.Open(DB_DRIVER, CONNECTION)
	ErrorCheck(err)
	defer conn.Close()
	
	responses = []int{}
	CreateAndUseDB(DB_NAME, conn)

	rows, err := conn.Query("SELECT RESPONSE_TIME FROM responses")
	if err != nil {
		return
	}
	for rows.Next() {
		var rtime int
		err = rows.Scan(&rtime)
		ErrorCheck(err)
		responses = append(responses, rtime)
	}
	return
}

// If not created yet, creates the application database and uses it
func CreateAndUseDB(s string, conn *sql.DB) {
	_, err := conn.Exec("CREATE DATABASE IF NOT EXISTS " + s)
	ErrorCheck(err)
	conn.Exec("USE " + s)
}

// Inserts the event data into the table with the same API_KEY values
// and if the table is not created yet, creates one
func InsertNewEvent(name string, e *model.Event, conn *sql.DB) {
	_, err := conn.Exec("CREATE TABLE IF NOT EXISTS " + name + "(API_KEY text, USER_ID int, TIMESTAMP int)")
	ErrorCheck(err)

	_, err = conn.Exec("INSERT INTO " + name + " (API_KEY, USER_ID, TIMESTAMP) VALUES (?, ?, ?)", e.ApiKey, e.UserID, e.Timestamp)
	ErrorCheck(err)

	random := rand.Intn(101) // generating a random integer between 0-100
	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS " + RESPONSE_TABLE + "(RESPONSE_TIME int)")
	ErrorCheck(err)
	_, err = conn.Exec("INSERT INTO " + RESPONSE_TABLE + " (RESPONSE_TIME) VALUES (?)", random)
	ErrorCheck(err)

	time.Sleep(time.Millisecond * time.Duration(random))
}

// Connects to database and sends the Event data
func AddEvent(e *model.Event) {
	conn, err := sql.Open(DB_DRIVER, CONNECTION)
	ErrorCheck(err)
	defer conn.Close()
	CreateAndUseDB(DB_NAME, conn)
	InsertNewEvent(e.ApiKey, e, conn)
}

