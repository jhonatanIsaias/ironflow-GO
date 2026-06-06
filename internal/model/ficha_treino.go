package model

type FichaTreino struct {
	BaseEntity
	FitNrID             int     `json:"fitNrId" db:"fit_nr_id"`
	TreNrID             int     `json:"treNrId" binding:"required" db:"tre_nr_id"` 
	ExeNrID             int     `json:"exeNrId" binding:"required" db:"exe_nr_id"` 
	FitNrOrdem          int     `json:"fitNrOrdem" binding:"required" db:"fit_nr_ordem"`
	FitNrMetaSeries     int     `json:"fitNrMetaSeries" db:"fit_nr_meta_series"`
	FitTxMetaRepeticoes string  `json:"fitTxMetaRepeticoes" binding:"required" db:"fit_tx_meta_repeticoes"`
	FitNrMetaPeso       float64 `json:"fitNrMetaPeso" binding:"required" db:"fit_nr_meta_peso"`
	FitNrGrupo          int     `json:"fitNrGrupo" db:"fit_nr_grupo"`
}

type FichaTreinoResponse struct {
    FichaTreino
    ExeTxNome           string  `json:"exeTxNome"` 
    
}
