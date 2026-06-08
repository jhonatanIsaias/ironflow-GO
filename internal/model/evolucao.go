package model

import "time"

type Evolucao struct {
	BaseEntity
	EvoNrID                  int       `json:"evoNrID" db:"evo_nr_id"`
	UsuNrID                  int       `json:"usuNrID" db:"usu_nr_id"`
	EvoDtData                time.Time `json:"evoDtData" db:"evo_dt_data"`
	EvoNrPeso                *float64  `json:"evoNrPeso" db:"evo_nr_peso"`
	EvoNrAltura              *float64  `json:"evoNrAltura" db:"evo_nr_altura"`
	EvoNrOmbro               *float64  `json:"evoNrOmbro" db:"evo_nr_ombro"`
	EvoNrBusto               *float64  `json:"evoNrBusto" db:"evo_nr_busto"`
	EvoNrAbdomen             *float64  `json:"evoNrAbdomen" db:"evo_nr_abdomen"`
	EvoNrCintura             *float64  `json:"evoNrCintura" db:"evo_nr_cintura"`
	EvoNrQuadril             *float64  `json:"evoNrQuadril" db:"evo_nr_quadril"`
	EvoNrBracoDireito        *float64  `json:"evoNrBracoDireito" db:"evo_nr_braco_direito"`
	EvoNrBracoEsquerdo       *float64  `json:"evoNrBracoEsquerdo" db:"evo_nr_braco_esquerdo"`
	EvoNrAntebracoDireito    *float64  `json:"evoNrAntebracoDireito" db:"evo_nr_antebraco_direito"`
	EvoNrAntebracoEsquerdo   *float64  `json:"evoNrAntebracoEsquerdo" db:"evo_nr_antebraco_esquerdo"`
	EvoNrCoxaDireita         *float64  `json:"evoNrCoxaDireita" db:"evo_nr_coxa_direita"`
	EvoNrCoxaEsquerda        *float64  `json:"evoNrCoxaEsquerda" db:"evo_nr_coxa_esquerda"`
	EvoNrPanturrilhaDireita  *float64  `json:"evoNrPanturrilhaDireita" db:"evo_nr_panturrilha_direita"`
	EvoNrPanturrilhaEsquerda *float64  `json:"evoNrPanturrilhaEsquerda" db:"evo_nr_panturrilha_esquerda"`
}