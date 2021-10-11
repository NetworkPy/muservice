package models

import "github.com/google/uuid"

type User struct {
	UID      uuid.UUID `db:"uid" json:"uid"`
	Email    string    `db:"email" json:"email"`
	Passowrd string    `db:"password" json:"password"`
	Name     string    `db:"Name" json:"Name"`
	ImageURL string    `db:"image_url" json:"imageUrl"`
	Website  string    `db:"website" json:"website"`
}
