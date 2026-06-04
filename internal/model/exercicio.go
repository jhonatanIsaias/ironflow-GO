package model

type Exercicio struct {
	BaseEntity 
	ExeNrID              int    `json:"exeNrId" db:"exe_nr_id"`
	ExeTxNome            string `json:"exeTxNome" db:"exe_tx_nome"`
	ExeTxGrupoMuscular   string `json:"exeTxGrupoMuscular" db:"exe_tx_grupo_muscular"`
	ExeTxGrupoMuscularSinergista string `json:"exeTxGrupoMuscularSinergista" db:"exe_tx_grupo_muscular_sinegista"`
	ExeTxTipoEquipamento string `json:"exeTxTipoEquipamento" db:"exe_tx_tipo_equipamento"`
}