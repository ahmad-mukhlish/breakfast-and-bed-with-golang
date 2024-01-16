package repository

import "github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"

type DatabaseRepository interface {
	GetUsers() bool
	InsertReservation(reservation model.Reservation) (int, error)
	InsertRoomRestriction(reservation model.RoomRestriction) error
	CheckAvailabilityForRoomById(startDate, endDate string, roomId int) (bool, error)
}
