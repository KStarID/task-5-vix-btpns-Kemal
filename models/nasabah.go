package models

type Nasabah struct {
	Nim string `json:"nim" gorm:"primary_key"`
	Nama string `json:"nama"`
}