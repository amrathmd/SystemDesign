package main

import (
	"database/sql"
	"fmt"
	"log"
	"pessimisticLocksInDB/common"
)

func Book(db *sql.DB, userId, showId int) error {
    var seat common.Seat
    // Check for available seats before booking
	tx,err := db.Begin();
    tx.QueryRow(`
        SELECT s.seatId, s.seatNumber
        FROM seats s
        LEFT JOIN bookings b ON s.seatId = b.seatId
        WHERE b.bookingId IS NULL AND s.showId = ?
        LIMIT 1 for update`, showId).Scan(&seat.SeatId, &seat.SeatNumber)
	defer func(){
		if p := recover(); p != nil{
			tx.Rollback()
			panic(p)

		}else if err != nil{
			tx.Rollback()
		}else {
			tx.Commit()
		}
	}()
    // If no available seats, return an error
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("no available seats for showId %d", showId)
        }
        log.Println("Error getting seats:", err)
        return err
    }

    // Insert booking for the selected seat
    queryToBook := `INSERT INTO bookings(userId, showId, seatId) VALUES(?, ?, ?)`
    _, err = tx.Exec(queryToBook, userId, showId, seat.SeatId)
    if err != nil {
        log.Println("Error booking seat:", err)
        return err
    }

    log.Printf("Booked seat %d for user %d on showId %d", seat.SeatId, userId, showId)
    return nil
}

func CheckBookedSeats(db *sql.DB,showId int){
	rows,err := db.Query(`SELECT DISTINCT seatId FROM bookings where showId=?;`,showId);
	if(err != nil){
		log.Fatalf("Error getting the seats that are filled %s",err.Error())
	}
	var seatsBooked[] int;
	for rows.Next(){
		var seat int;
		rows.Scan(&seat);
		seatsBooked = append(seatsBooked, seat)
	}
	fmt.Println("Booked seats:")
	log.Printf("%d",seatsBooked);
}