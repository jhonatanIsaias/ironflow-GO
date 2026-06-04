package model

type FichaTreino struct {
	BaseEntity
	FitNrID             int     `json:"fitNrId" db:"fit_nr_id"`
	TreNrID             int     `json:"treNrId" db:"tre_nr_id"` 
	ExeNrID             int     `json:"exeNrId" db:"exe_nr_id"` 
	FitNrOrdem          int     `json:"fitNrOrdem" db:"fit_nr_ordem"`
	FitNrMetaSeries     int     `json:"fitNrMetaSeries" db:"fit_nr_meta_series"`
	FitNrMetaRepeticoes int     `json:"fitNrMetaRepeticoes" db:"fit_nr_meta_repeticoes"`
	FitNrMetaPeso       float64 `json:"fitNrMetaPeso" db:"fit_nr_meta_peso"`
}

type FichaTreinoResponse struct {
    FitNrID             int     `json:"fitNrId"`
    TreNrID             int     `json:"treNrId"` 
    ExeNrID             int     `json:"exeNrId"` 
    ExeTxNome           string  `json:"exeTxNome"` 
    FitNrOrdem          int     `json:"fitNrOrdem"`
    FitNrMetaSeries     int     `json:"fitNrMetaSeries"`
    FitNrMetaRepeticoes int     `json:"fitNrMetaRepeticoes"`
    FitNrMetaPeso       float64 `json:"fitNrMetaPeso"`
}
