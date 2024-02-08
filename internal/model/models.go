package model

import (
	"time"
)

type User struct {
	Id          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccessLevel int
}

type Room struct {
	Id        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restriction struct {
	Id              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomId    int
	Room      Room
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoomRestriction struct {
	Id            int
	StartDate     time.Time
	EndDate       time.Time
	RoomId        int
	Room          Room
	ReservationId int
	Reservation   Reservation
	CreatedAt     time.Time
	UpdatedAt     time.Time
	RestrictionId int
	Restriction   Restriction
}

type MailData struct {
	To      string
	From    string
	Content string
	Subject string
}
