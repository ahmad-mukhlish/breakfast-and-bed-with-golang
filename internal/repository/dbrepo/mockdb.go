package dbrepo

import (
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
)

func (m *postgresDBRepository) InsertReservation(reservation model.Reservation) (int, error) {

	return 0, nil

}

func (m *postgresDBRepository) InsertRoomRestriction(roomRestriction model.RoomRestriction) error {

	return nil
}

func (m *postgresDBRepository) CheckAvailabilityForRoomById(startDate, endDate string, roomId int) (bool, error) {

	return false, nil
}

func (m *postgresDBRepository) GetAvailableRooms(startDate, endDate string) ([]model.Room, error) {

	var rooms []model.Room

	return rooms, nil
}

func (m *postgresDBRepository) GetRoomById(id int) (model.Room, error) {

	var room model.Room

	return room, nil
}
