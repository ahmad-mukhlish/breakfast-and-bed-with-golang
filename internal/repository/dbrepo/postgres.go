package dbrepo

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
)

func (m *postgresDBRepository) InsertReservation(reservation model.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	var newId int

	query := `INSERT
			into
    		reservations (first_name, last_name, email, phone,
                 start_date, end_date, room_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	rowContext := m.DB.QueryRowContext(ctx, query, reservation.FirstName,
		reservation.LastName, reservation.Email, reservation.Phone,
		reservation.StartDate, reservation.EndDate, reservation.RoomId, time.Now(), time.Now())

	if rowContext.Err() != nil {
		return 0, rowContext.Err()
	}

	err := rowContext.Scan(&newId)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (m *postgresDBRepository) InsertRoomRestriction(roomRestriction model.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	query := `INSERT
			into
    		room_restrictions (
                 start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, query, roomRestriction.StartDate,
		roomRestriction.EndDate, roomRestriction.RoomId, roomRestriction.ReservationId,
		roomRestriction.RestrictionId, time.Now(), time.Now())

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepository) CheckAvailabilityForRoomById(startDate, endDate string, roomId int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	// is not available if :

	// end date db >= start date params
	// AND
	// end date params >= start date db

	//examples:

	//db 			: 5 6 7 8
	//params        : 5 6 7 8

	//db 			: 5 6 7 8
	//params        :   6 7

	//db 			:   5 6 7 8
	//params        : 4 5 6 7 8 9

	//db 			: 5 6 7 8
	//params        :     7 8 9

	//db 			:   5 6 7 8
	//params        : 4 5 6

	//avails:

	//db 			:   5 6 7 8
	//params        : 4         (end date params [4] < start date db [5])

	//db 			:   5 6 7 8
	//params        :           9 10 (end date db [8] < start date params [9])

	query := `select
			  count(id)
			  from room_restrictions
			  where
			  end_date >= $1 and $2 >= start_date
			  and room_id = $3; `

	rowContext := m.DB.QueryRowContext(ctx, query, startDate, endDate, roomId)

	var count int

	if rowContext.Err() != nil {
		return false, rowContext.Err()
	}

	err := rowContext.Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, err
}

func (m *postgresDBRepository) GetAvailableRooms(startDate, endDate string) ([]model.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	query := `	select r.id, r.room_name
				from rooms r
				where r.id not in (
					select
					rr.room_id
					from room_restrictions rr
					where
					end_date >= $1 and $2 >= start_date
				);`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)

	var rooms []model.Room

	if err != nil {
		return rooms, err
	}

	//next is move the pointer for scan to next row
	for rows.Next() {

		var room model.Room
		err = rows.Scan(&room.Id, &room.RoomName)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)

	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresDBRepository) GetRoomById(id int) (model.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	query := `select r.id, r.room_name
				from rooms r
				where r.id = $1 ;`

	rows, err := m.DB.QueryContext(ctx, query, id)

	var room model.Room

	if err != nil {
		return room, err
	}

	//next is move the pointer for scan to next row
	for rows.Next() {

		err = rows.Scan(&room.Id, &room.RoomName)
		if err != nil {
			return room, err
		}

	}

	if err = rows.Err(); err != nil {
		return room, err
	}

	return room, nil
}

func (m *postgresDBRepository) GetUserById(id int) (model.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	query := `select u.id, u.first_name, u.last_name, u.email, u.access_level, u.created_at, u.updated_at 
				from users u
				where u.id = $1 ;`

	rowContext := m.DB.QueryRowContext(ctx, query, id)

	var user model.User

	if rowContext.Err() != nil {
		return user, rowContext.Err()
	}

	err := rowContext.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *postgresDBRepository) UpdateUser(user model.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	query := `UPDATE 
			users
    		set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5, `

	_, err := m.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.AccessLevel, time.Now())

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepository) Authenticate(inputEmail, inputPassword string) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	var id int
	var hashedPassword string

	query := `select u.id, u.password
				from users u
				where u.email = $1 ;`

	rowContext := m.DB.QueryRowContext(ctx, query, inputEmail)
	err := rowContext.Scan(&id, &hashedPassword)
	if err != nil {
		return 0, err
	}

	errCompareHashing := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))

	if errors.Is(errCompareHashing, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, errors.New("incorrect Password")
	} else if errCompareHashing != nil {
		return 0, errCompareHashing
	}

	return id, nil

}
