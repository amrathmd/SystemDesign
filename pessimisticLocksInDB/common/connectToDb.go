package common

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)
var queriesToSetupTables = []string{
	"CREATE TABLE IF NOT EXISTS users(userId INT AUTO_INCREMENT PRIMARY KEY,username VARCHAR(100))",
	"CREATE TABLE IF NOT EXISTS shows(showId INT AUTO_INCREMENT PRIMARY KEY,showname VARCHAR(100))",
	"CREATE TABLE IF NOT EXISTS seats(seatId INT AUTO_INCREMENT PRIMARY KEY,seatNumber VARCHAR(100),showId INT)",
	"CREATE TABLE IF NOT EXISTS bookings(bookingId INT AUTO_INCREMENT PRIMARY KEY,userId INT,showId INT,seatId INT)",
}
var alterQueries = []string{
	
}
var queriesToSetupShows = []string{
	"INSERT INTO shows(showId,showname) values (1,'arijitConcert')",
}
func generateSeatQueries(showId int, numSeats int) []string {
	queries := make([]string, 0, numSeats)
	for i := 1; i <= numSeats; i++ {
		seatNumber := fmt.Sprintf("Seat-%d", i) // Generate seat names like Seat-1, Seat-2, ...
		query := fmt.Sprintf("INSERT INTO seats (seatNumber, showId) VALUES ('%s', %d)", seatNumber, showId)
		queries = append(queries, query)
	}
	return queries
}
func generateUsersQueries(numUsers int) []string {
	queries := make([]string, 0, numUsers)
	for i := 1; i <= numUsers; i++ {
		username := fmt.Sprintf("User-%d", i) // Generate user names like User-1, User-2, ...
		query := fmt.Sprintf("INSERT INTO users (username) VALUES ('%s')", username)
		queries = append(queries, query)
	}
	return queries
}
func usersExist(db *sql.DB) bool{
	var count int;
	db.QueryRow("select COUNT(*) from users").Scan(&count);
	return count > 0
}
func seatsForShowExists(db *sql.DB,showId int) bool{
	var count int;
	db.QueryRow("select COUNT(*) from seats where showId = ?",showId).Scan(&count)
	return count > 0
}
func showExists(db *sql.DB,showId int) bool {
	var count int;
	db.QueryRow("select COUNT(*) from shows where id = ?",showId).Scan(&count)
	return count > 0
}
func ConnectToDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@/")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS gotest")
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}
	db.Close()
	db, err = sql.Open("mysql", "root:@/gotest")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
func SetupDatabase() (*sql.DB,error){
	db,err := ConnectToDb();
	if(err != nil){
		log.Fatalf("Error Connecting to db %s",err.Error());
		return nil,err;
	}
	log.Printf("Setting up tables in database")
	for _ , query := range queriesToSetupTables {
		_ ,err := db.Exec(query);
		if(err != nil){
			log.Printf("Error setting up database %s",err.Error())
		}
	}
	log.Println("Tables created successfully")


	if(!showExists(db,1)){
		log.Println("Inserting shows in the table")
		for _ , query := range queriesToSetupShows {
			_ ,err := db.Exec(query);
			if(err != nil){
				log.Printf("Error inserting show - show already exist")
			}
		}
	}
	if(!seatsForShowExists(db,1)){
		log.Printf("Inserting seats in seats table");
		seats := generateSeatQueries(1,100);
		for _ , query := range seats {
			_ ,err := db.Exec(query);
			if(err != nil){
				log.Printf("Error inserting seat - seat already exist")
			}
		}
	}
	if(!usersExist(db)){
		log.Printf("Inserting users in users table");
		users := generateUsersQueries(100);
		for _ , query := range users {
			_ ,err := db.Exec(query);
			if(err != nil){
				log.Printf("Error inserting user - user already exist")
			}
		}
	}
	log.Printf("Database successfully setup");
	return db,nil
}