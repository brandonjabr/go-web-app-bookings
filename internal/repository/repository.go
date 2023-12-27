package repository

import (
	"time"

	"github.com/brandonjabr/go-web-app-bookings/internal/models"
)

type DatabaseRepo interface {
	InsertReservation(reservation models.Reservation) (int, error)
	InsertRoomRestriction(roomRestriction models.RoomRestriction) error
	SearchAvailabilityByDates(startDate, endDate time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error)
	GetRoomByID(roomID int) (models.Room, error)
}
