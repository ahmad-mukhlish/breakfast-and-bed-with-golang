package dbrepo

import (
	"errors"

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

	return false, nil
}

func (m *mockDBRepository) GetAvailableRooms(startDate, endDate string) ([]model.Room, error) {

	var rooms []model.Room

	return rooms, nil
}

func (m *mockDBRepository) GetRoomById(id int) (model.Room, error) {

	var room model.Room

	if id > 2 {
		return room, errors.New("Room Does Not Exist")
	}

	return room, nil
}
