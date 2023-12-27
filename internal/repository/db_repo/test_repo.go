package db_repo

import (
	"errors"
	"time"

	"github.com/brandonjabr/go-web-app-bookings/internal/models"
)

func (repo *testDBRepo) InsertReservation(reservation models.Reservation) (int, error) {
	return 0, nil
}

func (repo *testDBRepo) InsertRoomRestriction(roomRestriction models.RoomRestriction) error {
	return nil
}

func (repo *testDBRepo) SearchAvailabilityByDates(startDate, endDate time.Time, roomID int) (bool, error) {
	return false, nil
}

func (repo *testDBRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}

func (repo *testDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	var room models.Room

	if roomID > 2 {
		return room, errors.New("invalid room id")
	}

	return room, nil
}
