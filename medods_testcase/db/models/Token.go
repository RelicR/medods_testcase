package models

import (
	"github.com/gofrs/uuid"
)

type Token struct {
	Id               uint      `json:"id" gorm:"primaryKey"`
	UserGuid         uuid.UUID `gorm:"unique;notNull"`
	Refresh_token    string    `json:"refresh_token" gorm:"type:varchar(255);unique"`
	Last_ip          string    `json:"last_ip" gorm:"type:varchar(255)"`
	Last_fingerprint string    `json:"last_fp" gorm:"type:varchar(255)"`
	PairId           string    `json:"pair_id" gorm:"type:varchar(255);unique"`
	User             User      `gorm:"foreignKey:UserGuid"`
}
