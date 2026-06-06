package model

type SerieExecutada struct {
	BaseEntity
	SexNrID                   int     `json:"sexNrId" db:"sex_nr_id"`
	SetNrID                   int     `json:"setNrId" binding:"required" db:"set_nr_id"` 
	FitNrID                   int     `json:"fitNrId" binding:"required" db:"fit_nr_id"` 
	SexNrSerieNumero          int     `json:"sexNrSerieNumero" binding:"required" db:"sex_nr_serie_numero"`
	SexTxRepeticoesExecutadas string  `json:"sexTxRepeticoesExecutadas" binding:"required" db:"sex_tx_repeticoes_executadas"`
	SexNrPesoUtilizado        float64 `json:"sexNrPesoUtilizado" binding:"required" db:"sex_nr_peso_utilizado"`
}


type SerieExecutadaDetalhada struct {
	SerieExecutada          
	ExeTxNome      string `json:"exeTxNome"`      
	FitNrOrdem     int    `json:"fitNrOrdem"`     
}