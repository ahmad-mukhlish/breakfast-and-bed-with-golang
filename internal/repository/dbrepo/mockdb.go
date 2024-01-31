package dbrepo

import (
	"errors"
	"time"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
)

func (m *mockDBRepository) InsertReservation(reservation model.Reservation) (int, error) {

	if reservation.RoomId == 1000 {
		return 0, errors.New("cannot Insert Into Reservation")
	}

	return 0, nil

}

func (m *mockDBRepository) InsertRoomRestriction(roomRestriction model.RoomRestriction) error {

	if roomRestriction.RoomId == 10000 {
		return errors.New("cannot Insert Into RoomRestriction")
	}

	return nil
}

func (m *mockDBRepository) CheckAvailabilityForRoomById(startDate, endDate string, roomId int) (bool, error) {

	if roomId == 1 {
		return true, nil
	}
	if roomId == 1000 {
		return true, errors.New("mocking error")
	}

	return false, nil
}

func (m *mockDBRepository) GetAvailableRooms(startDate, endDate string) ([]model.Room, error) {

	var rooms []model.Room

	//case db error
	if startDate == "2006-01-02" {
		return rooms, errors.New("error get available Rooms")
	}

	//case empty room
	if startDate == "2006-01-03" {
		return rooms, nil
	}

	room := model.Room{
		Id:        1,
		RoomName:  "abajadun",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	rooms =
		append(
			rooms,
			room,
		)

	return rooms, nil
}

func (m *mockDBRepository) GetRoomById(id int) (model.Room, error) {

	var room model.Room

	if id > 2 {
		return room, errors.New("Room Does Not Exist")
	}

	return room, nil
}
