package main

import (
	"database/sql"
	"fmt"
	"log"
	"pessimisticLocksInDB/common"
)

func Book(db *sql.DB, userId, showId int) error {
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("failed to start transaction: %v", err)
    }
    defer func(tx *sql.Tx) {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        } else if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }(tx)

    // Modified query to avoid the JOIN and use EXISTS instead
    var seat common.Seat
    err = tx.QueryRow(`
        SELECT s.seatId, s.seatNumber 
        FROM seats s 
        WHERE s.showId = ? 
        AND NOT EXISTS (
            SELECT 1 FROM bookings b 
            WHERE b.seatId = s.seatId
        )
        LIMIT 1 
        FOR UPDATE SKIP LOCKED`, showId).Scan(&seat.SeatId, &seat.SeatNumber)

    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("no available seats for showId %d", showId)
        }
        return fmt.Errorf("error getting seats: %v", err)
    }

    // Insert booking for the selected seat
    _, err = tx.Exec(`
        INSERT INTO bookings(userId, showId, seatId) 
        VALUES(?, ?, ?)`, userId, showId, seat.SeatId)
    if err != nil {
        return fmt.Errorf("error booking seat: %v", err)
    }

    log.Printf("Booked seat %s for user %d on showId %d", seat.SeatNumber, userId, showId)
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
	log.Printf("%d : length = %d",seatsBooked,len(seatsBooked));
}