package model

type Treino struct {
	BaseEntity
	TreNrID   int    `json:"treNrId" db:"tre_nr_id"`
	TreTxNome string `json:"treTxNome" db:"tre_tx_nome"`
}