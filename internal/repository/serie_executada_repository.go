package repository

import (
	"context"
	"errors"
	"fmt"
	"ironflow/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SerieExecutadaRepository struct {
	DB *pgxpool.Pool
}

func NovoSerieExecutadaRepository(db *pgxpool.Pool) *SerieExecutadaRepository {
	return &SerieExecutadaRepository{DB: db}
}

// No arquivo serie_executada_repository.go

func (r *SerieExecutadaRepository) RegistrarSerieComSessaoAutomatica(ctx context.Context, serie *model.SerieExecutada,treNrId int) error {

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}

	defer tx.Rollback(ctx)

	var setNrID int

	sqlBusca := `
		SELECT set_nr_id FROM treino.set_sessao_treino 
		WHERE tre_nr_id = $1 AND set_dt_data = CURRENT_DATE AND deleted_at IS NULL 
		LIMIT 1
	`
	err = tx.QueryRow(ctx, sqlBusca, treNrId).Scan(&setNrID)

	if errors.Is(err, pgx.ErrNoRows) {
		sqlInsertSessao := `
			INSERT INTO treino.set_sessao_treino (tre_nr_id, set_dt_data, set_tm_hora_inicio)
			VALUES ($1, CURRENT_DATE, CURRENT_TIME) 
			RETURNING set_nr_id
		`
		err = tx.QueryRow(ctx, sqlInsertSessao, treNrId).Scan(&setNrID)
		if err != nil {
			return fmt.Errorf("erro ao criar sessão automática: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("erro ao buscar sessão existente: %w", err)
	}

	serie.SetNrID = setNrID
	sqlInsertSerie := `
		INSERT INTO treino.sex_serie_executada 
		(set_nr_id, fit_nr_id, sex_nr_serie_numero, sex_nr_repeticoes_realizadas, sex_nr_peso_utilizado)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING sex_nr_id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, sqlInsertSerie, 
		serie.SetNrID, 
		serie.FitNrID, 
		serie.SexNrSerieNumero, 
		serie.SexTxRepeticoesExecutadas, 
		serie.SexNrPesoUtilizado,
	).Scan(
		&serie.SexNrID,
		&serie.CreatedAt,
		&serie.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar a série executada: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *SerieExecutadaRepository) Editar(ctx context.Context, serie *model.SerieExecutada) error {
	
	sql := `UPDATE treino.set_serie_executada
	SET
	fit_nr_id = $3, 
	sex_nr_serie_numero = $4, 
	sex_nr_peso_utilizado = $5, 
	sex_tx_repeticoes_executadas = $6,
	updated_at = NOW()
	WHERE sex_nr_id = $1 AND deleted_at IS NULL
	RETURNING created_at, updated_at`

	err := r.DB.QueryRow(ctx, sql,
		serie.SexNrID,
		serie.SetNrID, 
		serie.FitNrID, 
		serie.SexNrSerieNumero, 
		serie.SexNrPesoUtilizado, 
		serie.SexTxRepeticoesExecutadas).Scan(
			&serie.CreatedAt,
			&serie.UpdatedAt,
		)
	if err != nil {
		return err
	}
	return nil
}

func (r *SerieExecutadaRepository) BuscarPorFichaTreino(ctx context.Context, fitNrId int) ([]model.SerieExecutada, error) {
	sql := `SELECT sex_nr_id, set_nr_id, fit_nr_id, sex_nr_serie_numero, sex_nr_peso_utilizado, sex_tx_repeticoes_executadas, created_at, updated_at 
	FROM treino.set_serie_executada WHERE fit_nr_id = $1 AND deleted_at IS NULL`

	rows, err := r.DB.Query(ctx, sql, fitNrId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var series []model.SerieExecutada
	for rows.Next() {
		var serie model.SerieExecutada
		err := rows.Scan(
			&serie.SexNrID,
			&serie.SetNrID,
			&serie.FitNrID,
			&serie.SexNrSerieNumero,
			&serie.SexNrPesoUtilizado,
			&serie.SexTxRepeticoesExecutadas,
			&serie.CreatedAt,
			&serie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		series = append(series, serie)
	}

	return series, nil
}

func (r *SerieExecutadaRepository) BuscarPorSessao(ctx context.Context, setNrId int) ([]model.SerieExecutadaDetalhada, error) {
	sql := `
		SELECT 
			sex.sex_nr_id, 
			sex.set_nr_id, 
			sex.fit_nr_id, 
			sex.sex_nr_serie_numero, 
			sex.sex_nr_peso_utilizado, 
			sex.sex_tx_repeticoes_executadas,
			f.fit_nr_ordem, 
			e.exe_tx_nome
		FROM treino.sex_serie_executada sex
		JOIN treino.fit_ficha_treino f ON sex.fit_nr_id = f.fit_nr_id
		JOIN treino.exe_exercicio e ON f.exe_nr_id = e.exe_nr_id
		WHERE sex.set_nr_id = $1 AND sex.deleted_at IS NULL
		ORDER BY f.fit_nr_ordem ASC, sex.sex_nr_serie_numero ASC
	`

	rows, err := r.DB.Query(ctx, sql, setNrId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var series []model.SerieExecutadaDetalhada

	for rows.Next() {
		var s model.SerieExecutadaDetalhada
		err := rows.Scan(
			&s.SexNrID,
			&s.SetNrID,
			&s.FitNrID,
			&s.SexNrSerieNumero,
			&s.SexNrPesoUtilizado,
			&s.SexTxRepeticoesExecutadas,
			&s.FitNrOrdem,
			&s.ExeTxNome,
		)
		if err != nil {
			return nil, err
		}
		series = append(series, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return series, nil
}

func (r *SerieExecutadaRepository) Deletar(ctx context.Context, sexNrId int) error {
	sql := `UPDATE treino.set_serie_executada SET deleted_at = NOW() WHERE sex_nr_id = $1 AND deleted_at IS NULL`

	_, err := r.DB.Exec(ctx, sql, sexNrId)
	if err != nil {
		return err
	}
	return nil
}