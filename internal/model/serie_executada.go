package model

type SerieExecutada struct {
	BaseEntity
	SexNrID                   int     `json:"sexNrId" db:"sex_nr_id"`
	SetNrID                   int     `json:"setNrId" db:"set_nr_id"` 
	FitNrID                   int     `json:"fitNrId" db:"fit_nr_id"` 
	SexNrSerieNumero          int     `json:"sexNrSerieNumero" db:"sex_nr_serie_numero"`
	SexNrRepeticoesRealizadas int     `json:"sexNrRepeticoesRealizadas" db:"sex_nr_repeticoes_realizadas"`
	SexNrPesoUtilizado        float64 `json:"sexNrPesoUtilizado" db:"sex_nr_peso_utilizado"`
}