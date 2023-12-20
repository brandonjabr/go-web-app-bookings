package db_repo

import (
	"context"
	"time"

	"github.com/brandonjabr/go-web-app-bookings/internal/models"
)

func (repo *postgresDBRepo) InsertReservation(reservation models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newReservationID int

	query := `INSERT INTO RESERVATIONS (FIRST_NAME, LAST_NAME, EMAIL, PHONE_NUMBER, 
				CHECK_IN_DATE, CHECK_OUT_DATE, ROOM_ID, CREATED_AT, UPDATED_AT)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING ID`

	err := repo.DB.QueryRowContext(ctx, query,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.PhoneNumber,
		reservation.CheckInDate,
		reservation.CheckOutDate,
		reservation.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newReservationID)
	if err != nil {
		return 0, err
	}

	return newReservationID, nil
}

func (repo *postgresDBRepo) InsertRoomRestriction(roomRestriction models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO ROOM_RESTRICTIONS (START_DATE, END_DATE, ROOM_ID, RESERVATION_ID, 
				RESTRICTION_ID, CREATED_AT, UPDATED_AT)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := repo.DB.ExecContext(ctx, query,
		roomRestriction.StartDate,
		roomRestriction.EndDate,
		roomRestriction.RoomID,
		roomRestriction.ReservationID,
		roomRestriction.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgresDBRepo) SearchAvailabilityByDates(startDate, endDate time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var numRows int

	query := `SELECT COUNT(ID) FROM ROOM_RESTRICTIONS
				WHERE ROOM_ID = $1 AND START_DATE < $2 AND END_DATE > $3`

	row := repo.DB.QueryRowContext(ctx, query, roomID, endDate, startDate)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *postgresDBRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `SELECT R.ID R.NAME FROM ROOMS R
				WHERE R.ID NOT IN 
					(SELECT RR.ROOM_ID FROM ROOM_RESTRICTIONS RR
						WHERE RR.START_DATE < $1 AND RR.END_DATE > $2)`

	rows, err := repo.DB.QueryContext(ctx, query, endDate, startDate)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
