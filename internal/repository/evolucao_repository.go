package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ironflow/internal/model"
)

type EvolucaoRepository struct {
	DB *pgxpool.Pool
}

func NovoEvolucaoRepository(db *pgxpool.Pool) *EvolucaoRepository {
	return &EvolucaoRepository{DB: db}
}

func (r *EvolucaoRepository) Salvar(ctx context.Context, e *model.Evolucao, usuTxID string) error {
	sql := `
        INSERT INTO treino.evo_evolucao (
			usu_tx_id, evo_dt_data, evo_nr_peso, evo_nr_altura,
			evo_nr_ombro, evo_nr_busto, evo_nr_abdomen, evo_nr_cintura,
			evo_nr_quadril, evo_nr_braco_direito, evo_nr_braco_esquerdo,
			evo_nr_antebraco_direito, evo_nr_antebraco_esquerdo,
			evo_nr_coxa_direita, evo_nr_coxa_esquerda,
			evo_nr_panturrilha_direita, evo_nr_panturrilha_esquerda
		)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
        RETURNING evo_nr_id, created_at, updated_at`

	err := r.DB.QueryRow(ctx, sql,
		usuTxID,
		time.Now(),
		e.EvoNrPeso,
		e.EvoNrAltura,
		e.EvoNrOmbro,
		e.EvoNrBusto,
		e.EvoNrAbdomen,
		e.EvoNrCintura,
		e.EvoNrQuadril,
		e.EvoNrBracoDireito,
		e.EvoNrBracoEsquerdo,
		e.EvoNrAntebracoDireito,
		e.EvoNrAntebracoEsquerdo,
		e.EvoNrCoxaDireita,
		e.EvoNrCoxaEsquerda,
		e.EvoNrPanturrilhaDireita,
		e.EvoNrPanturrilhaEsquerda,
	).Scan(
		&e.EvoNrID,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	return err
}

func (r *EvolucaoRepository) Editar(ctx context.Context, e *model.Evolucao, usuTxID string) error {
	sql := `
		UPDATE treino.evo_evolucao
		SET evo_dt_data = $1, evo_nr_peso = $2, evo_nr_altura = $3,
			evo_nr_ombro = $4, evo_nr_busto = $5, evo_nr_abdomen = $6,
			evo_nr_cintura = $7, evo_nr_quadril = $8, evo_nr_braco_direito = $9,
			evo_nr_braco_esquerdo = $10, evo_nr_antebraco_direito = $11,
			evo_nr_antebraco_esquerdo = $12, evo_nr_coxa_direita = $13,
			evo_nr_coxa_esquerda = $14, evo_nr_panturrilha_direita = $15,
			evo_nr_panturrilha_esquerda = $16, updated_at = NOW()
		WHERE evo_nr_id = $17
		AND usu_tx_id = $18
		AND deleted_at IS NULL
		RETURNING created_at, updated_at
	`
	err := r.DB.QueryRow(ctx, sql,
		e.EvoDtData,
		e.EvoNrPeso,
		e.EvoNrAltura,
		e.EvoNrOmbro,
		e.EvoNrBusto,
		e.EvoNrAbdomen,
		e.EvoNrCintura,
		e.EvoNrQuadril,
		e.EvoNrBracoDireito,
		e.EvoNrBracoEsquerdo,
		e.EvoNrAntebracoDireito,
		e.EvoNrAntebracoEsquerdo,
		e.EvoNrCoxaDireita,
		e.EvoNrCoxaEsquerda,
		e.EvoNrPanturrilhaDireita,
		e.EvoNrPanturrilhaEsquerda,
		e.EvoNrID,
		usuTxID,
	).Scan(
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("Não é possível editar: Evolução inexistente")
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *EvolucaoRepository) BuscarPorID(ctx context.Context, evoNrID int, usuTxID string) (*model.Evolucao, error) {
	sql := `
        SELECT evo_nr_id, usu_tx_id, evo_dt_data, evo_nr_peso, evo_nr_altura,
			evo_nr_ombro, evo_nr_busto, evo_nr_abdomen, evo_nr_cintura,
			evo_nr_quadril, evo_nr_braco_direito, evo_nr_braco_esquerdo,
			evo_nr_antebraco_direito, evo_nr_antebraco_esquerdo,
			evo_nr_coxa_direita, evo_nr_coxa_esquerda,
			evo_nr_panturrilha_direita, evo_nr_panturrilha_esquerda
        FROM treino.evo_evolucao
        WHERE evo_nr_id = $1
		AND usu_tx_id = $2
		AND deleted_at IS NULL
    `

	var e model.Evolucao
	err := r.DB.QueryRow(ctx, sql, evoNrID, usuTxID).Scan(
		&e.EvoNrID,
		&e.UsuTxID,
		&e.EvoDtData,
		&e.EvoNrPeso,
		&e.EvoNrAltura,
		&e.EvoNrOmbro,
		&e.EvoNrBusto,
		&e.EvoNrAbdomen,
		&e.EvoNrCintura,
		&e.EvoNrQuadril,
		&e.EvoNrBracoDireito,
		&e.EvoNrBracoEsquerdo,
		&e.EvoNrAntebracoDireito,
		&e.EvoNrAntebracoEsquerdo,
		&e.EvoNrCoxaDireita,
		&e.EvoNrCoxaEsquerda,
		&e.EvoNrPanturrilhaDireita,
		&e.EvoNrPanturrilhaEsquerda,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("Evolução não encontrada")
	} else if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *EvolucaoRepository) BuscarTodos(ctx context.Context, usuTxID string) ([]model.Evolucao, error) {
	sql := `
        SELECT evo_nr_id, usu_tx_id, evo_dt_data, evo_nr_peso, evo_nr_altura,
			evo_nr_ombro, evo_nr_busto, evo_nr_abdomen, evo_nr_cintura,
			evo_nr_quadril, evo_nr_braco_direito, evo_nr_braco_esquerdo,
			evo_nr_antebraco_direito, evo_nr_antebraco_esquerdo,
			evo_nr_coxa_direita, evo_nr_coxa_esquerda,
			evo_nr_panturrilha_direita, evo_nr_panturrilha_esquerda
        FROM treino.evo_evolucao
        WHERE usu_tx_id = $1 AND deleted_at IS NULL
		ORDER BY evo_dt_data DESC
    `

	rows, err := r.DB.Query(ctx, sql, usuTxID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evolucoes []model.Evolucao
	for rows.Next() {
		var e model.Evolucao
		err := rows.Scan(
			&e.EvoNrID,
			&e.UsuTxID,
			&e.EvoDtData,
			&e.EvoNrPeso,
			&e.EvoNrAltura,
			&e.EvoNrOmbro,
			&e.EvoNrBusto,
			&e.EvoNrAbdomen,
			&e.EvoNrCintura,
			&e.EvoNrQuadril,
			&e.EvoNrBracoDireito,
			&e.EvoNrBracoEsquerdo,
			&e.EvoNrAntebracoDireito,
			&e.EvoNrAntebracoEsquerdo,
			&e.EvoNrCoxaDireita,
			&e.EvoNrCoxaEsquerda,
			&e.EvoNrPanturrilhaDireita,
			&e.EvoNrPanturrilhaEsquerda,
		)
		if err != nil {
			return nil, err
		}
		evolucoes = append(evolucoes, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return evolucoes, nil
}

func (r *EvolucaoRepository) BuscarMaisRecente(ctx context.Context, usuTxID string) (*model.Evolucao, error) {
	sql := `
        SELECT evo_nr_id, usu_tx_id, evo_dt_data, evo_nr_peso, evo_nr_altura,
			evo_nr_ombro, evo_nr_busto, evo_nr_abdomen, evo_nr_cintura,
			evo_nr_quadril, evo_nr_braco_direito, evo_nr_braco_esquerdo,
			evo_nr_antebraco_direito, evo_nr_antebraco_esquerdo,
			evo_nr_coxa_direita, evo_nr_coxa_esquerda,
			evo_nr_panturrilha_direita, evo_nr_panturrilha_esquerda
        FROM treino.evo_evolucao
        WHERE usu_tx_id = $1 AND deleted_at IS NULL
		ORDER BY evo_dt_data DESC
		LIMIT 1
    `

	var e model.Evolucao
	err := r.DB.QueryRow(ctx, sql, usuTxID).Scan(
		&e.EvoNrID,
		&e.UsuTxID,
		&e.EvoDtData,
		&e.EvoNrPeso,
		&e.EvoNrAltura,
		&e.EvoNrOmbro,
		&e.EvoNrBusto,
		&e.EvoNrAbdomen,
		&e.EvoNrCintura,
		&e.EvoNrQuadril,
		&e.EvoNrBracoDireito,
		&e.EvoNrBracoEsquerdo,
		&e.EvoNrAntebracoDireito,
		&e.EvoNrAntebracoEsquerdo,
		&e.EvoNrCoxaDireita,
		&e.EvoNrCoxaEsquerda,
		&e.EvoNrPanturrilhaDireita,
		&e.EvoNrPanturrilhaEsquerda,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *EvolucaoRepository) Deletar(ctx context.Context, evoNrID int, usuTxD string) error {
	sql := `UPDATE treino.evo_evolucao SET deleted_at = NOW() WHERE evo_nr_id = $1 AND usu_tx_id = $2::UUID AND deleted_at IS NULL`
	comando, err := r.DB.Exec(ctx, sql, evoNrID, usuTxD)
	if err != nil {
		return err
	}
	if comando.RowsAffected() == 0 {
		return errors.New("Não é possível deletar: Evolução inexistente ou não pertence ao usuário")
	}

	return nil
}
