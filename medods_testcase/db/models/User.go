package models

import (
	uuid "github.com/gofrs/uuid"
)

type User struct {
	GUID       uuid.UUID `json:"guid" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	First_name string    `json:"first_name" gorm:"type:varchar(25);notNull"`
	Last_name  string    `json:"last_name" gorm:"type:varchar(25);notNull"`
	Email      string    `json:"email" gorm:"type:varchar(50);notNull;unique"`
}
