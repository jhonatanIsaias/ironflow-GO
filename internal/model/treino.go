package model

type Treino struct {
	BaseEntity
	TreNrID   int    `json:"treNrId" db:"tre_nr_id"`
	TreTxNome string `json:"treTxNome" binding:"required" db:"tre_tx_nome"`
	TreTxDescricao string `json:"treTxDescricao" binding:"required" db:"tre_tx_descricao"`
}