package models

import "time"

type Reservation struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

type Users struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Rooms struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restrictions struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoomRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	Room          Rooms
	Reservation   Reservations
	Restriction   Restrictions
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Reservations struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	PhoneNumber  string
	CheckInDate  time.Time
	CheckOutDate time.Time
	RoomID       int
	Room         Rooms
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
