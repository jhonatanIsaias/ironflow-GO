package model

import "time"

type SessaoTreino struct {
	BaseEntity
	SetNrID         int       `json:"setNrId" db:"set_nr_id"`
	TreNrID         int       `json:"treNrId" db:"tre_nr_id"` 
	SetDtData       time.Time `json:"setDtData" db:"set_dt_data"`
	SetTmHoraInicio time.Time `json:"setTmHoraInicio" db:"set_tm_hora_inicio"`
}