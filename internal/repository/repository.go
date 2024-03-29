package repository

import "github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"

type DatabaseRepository interface {
	InsertReservation(reservation model.Reservation) (int, error)
	InsertRoomRestriction(reservation model.RoomRestriction) error
	CheckAvailabilityForRoomById(startDate, endDate string, roomId int) (bool, error)
	GetAvailableRooms(startDate, endDate string) ([]model.Room, error)
	GetRoomById(id int) (model.Room, error)
	GetUserById(id int) (model.User, error)
	UpdateUser(user model.User) error
	Authenticate(inputEmail, inputPassword string) (int, error)
}
