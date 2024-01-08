package dbrepo

import (
	"context"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"time"
)

func (m *postgresDBRepository) GetUsers() bool {
	return true
}

func (m *postgresDBRepository) InsertReservation(reservation model.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `INSERT 
			into 
    		reservation (first_name, last_name, email, phone, 
                 start_date, end_date, room_id, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := m.DB.ExecContext(ctx, query, reservation.FirstName,
		reservation.LastName, reservation.Email, reservation.Phone,
		reservation.StartDate, reservation.EndDate, reservation.RoomId, time.Now(), time.Now())

	if err != nil {
		return err
	}
	return nil
}
