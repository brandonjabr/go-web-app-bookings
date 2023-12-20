package models

import "time"

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Room struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Reservation struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	PhoneNumber  string
	CheckInDate  time.Time
	CheckOutDate time.Time
	RoomID       int
	Room         Room
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Restriction struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
