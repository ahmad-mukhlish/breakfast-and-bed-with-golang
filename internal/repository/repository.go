package repository

import "github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"

type DatabaseRepository interface {
	GetUsers() bool
	InsertReservation(reservation model.Reservation) error
}
