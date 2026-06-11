package model

import "github.com/google/uuid"

type Exercicio struct {
	BaseEntity
	ExeNrID                      int        `json:"exeNrId" db:"exe_nr_id"`
	ExeTxNome                    string     `json:"exeTxNome" binding:"required" db:"exe_tx_nome"`
	UsuTxID                      *uuid.UUID `db:"usu_tx_id" json:"-"`
	ExeTxGrupoMuscular           string     `json:"exeTxGrupoMuscular" db:"exe_tx_grupo_muscular"`
	ExeTxGrupoMuscularSinergista string     `json:"exeTxGrupoMuscularSinergista" db:"exe_tx_grupo_muscular_sinegista"`
	IsCustom                     bool       `json:"isCustom"`
}