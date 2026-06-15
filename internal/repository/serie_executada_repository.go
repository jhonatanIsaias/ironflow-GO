package repository

import (
	"context"
	"fmt"
	"ironflow/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SerieExecutadaRepository struct {
	DB *pgxpool.Pool
}

func NovoSerieExecutadaRepository(db *pgxpool.Pool) *SerieExecutadaRepository {
	return &SerieExecutadaRepository{DB: db}
}

func (r *SerieExecutadaRepository) RegistrarSerie(ctx context.Context, serie *model.SerieExecutada) error {

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}

	defer tx.Rollback(ctx)

	sqlInsertSerie := `
		INSERT INTO treino.sex_serie_executada 
		(set_nr_id, fit_nr_id, sex_nr_serie_numero, sex_tx_repeticoes_executadas, sex_nr_peso_utilizado)
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

func (r *SerieExecutadaRepository) Editar(ctx context.Context, serie *model.SerieExecutada, usuTxId string) error {

	sql := `
	UPDATE treino.sex_serie_executada
	SET
	fit_nr_id = $2, 
	sex_nr_serie_numero = $3, 
	sex_nr_peso_utilizado = $4, 
	sex_tx_repeticoes_executadas = $5,
	updated_at = NOW()
	WHERE sex_nr_id = $1 AND deleted_at IS NULL
	AND sex_nr_id IN (
		SELECT sex.sex_nr_id
		FROM treino.sex_serie_executada sex
		JOIN treino.set_sessao_treino s ON sex.set_nr_id = s.set_nr_id
		JOIN treino.tre_treino t ON s.tre_nr_id = t.tre_nr_id
		WHERE t.usu_tx_id = $6
	)
	RETURNING created_at, updated_at`

	err := r.DB.QueryRow(ctx, sql,
		serie.SexNrID,
		serie.FitNrID,
		serie.SexNrSerieNumero,
		serie.SexNrPesoUtilizado,
		serie.SexTxRepeticoesExecutadas,
		usuTxId).Scan(
		&serie.CreatedAt,
		&serie.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *SerieExecutadaRepository) BuscarPorFichaTreino(ctx context.Context, fitNrId int, usuTxId string) ([]model.SerieExecutada, error) {
	sql := `
	SELECT sex_nr_id, set_nr_id, fit_nr_id, sex_nr_serie_numero, sex_nr_peso_utilizado, sex_tx_repeticoes_executadas, created_at, updated_at 
	FROM treino.sex_serie_executada sex
	JOIN treino.fit_ficha_treino f ON sex.fit_nr_id = f.fit_nr_id
	JOIN treino.tre_treino t ON f.tre_nr_id = t.tre_nr_id
	WHERE fit_nr_id = $1 AND deleted_at IS NULL AND tre.usu_tx_id = $2`

	rows, err := r.DB.Query(ctx, sql, fitNrId, usuTxId)
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

func (r *SerieExecutadaRepository) BuscarPorSessao(ctx context.Context, setNrId int, usuTxId string) ([]model.SerieExecutadaDetalhada, error) {
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
		JOIN treino.tre_treino t ON f.tre_nr_id = t.tre_nr_id
		JOIN treino.exe_exercicio e ON f.exe_nr_id = e.exe_nr_id
		WHERE sex.set_nr_id = $1 AND sex.deleted_at IS NULL AND tre.usu_tx_id = $2
		ORDER BY f.fit_nr_ordem ASC, sex.sex_nr_serie_numero ASC
	`

	rows, err := r.DB.Query(ctx, sql, setNrId, usuTxId)
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

func (r *SerieExecutadaRepository) Deletar(ctx context.Context, sexNrId int, usuTxId string) error {
	sql := `
	UPDATE treino.sex_serie_executada 
	SET deleted_at = NOW() 
	WHERE sex_nr_id = $1 AND deleted_at IS NULL
	AND sex_nr_id IN (
		SELECT sex.sex_nr_id
		FROM treino.sex_serie_executada sex
		JOIN treino.set_sessao_treino s ON sex.set_nr_id = s.set_nr_id
		JOIN treino.tre_treino t ON s.tre_nr_id = t.tre_nr_id
		WHERE t.usu_tx_id = $2
	)`

	_, err := r.DB.Exec(ctx, sql, sexNrId, usuTxId)
	if err != nil {
		return err
	}
	return nil
}
