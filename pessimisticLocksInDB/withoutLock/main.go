package main

import (
	"log"
	"pessimisticLocksInDB/common"
	"sync"
	"time"
)



func main(){
	startTime := time.Now()
	db,err := common.SetupDatabase()
	var wg sync.WaitGroup
	if(err != nil){
		log.Fatalf("Error setting up the database %s",err.Error())
	}
	log.Print("----------Now try booking seats for 500 users concurrently--------");
	rows,err := db.Query("select * from users")
	if(err != nil){
		log.Fatalf("Error fetching users in the table")
	}
	defer rows.Close()
	var users[] common.User;
	for rows.Next(){
		var user common.User;
		err := rows.Scan(&user.UserId,&user.Username)
		if(err != nil){
			log.Fatal(err)
		}
		users = append(users, user)
	}
	log.Printf("---- Bookings started : Booking seats for every user here ------")
	db.Exec("TRUNCATE TABLE bookings")
	for _,user := range users {
		wg.Add(1)
		go func(user common.User){
			defer wg.Done()
			err := Book(db,user.UserId,1)
			if(err != nil){
				log.Printf("Error occured booking seat for user %s",user.Username)
			}
		}(user)
	}
	wg.Wait()
	CheckBookedSeats(db,1)
	endtime := time.Now()
	totalExectionTime := endtime.Sub(startTime)
	log.Printf("Total execution time is , %f",totalExectionTime.Seconds())
}
