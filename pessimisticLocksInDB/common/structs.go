package common

type User struct {
	UserId int
	Username string
}
type Show struct {
	ShowId int
	showname string
}
type Booking struct {
	BookingId int
	showId int
	userId int
	seatId int
}
type Seat struct {
	SeatId int
	ShowId int
	SeatNumber string
}