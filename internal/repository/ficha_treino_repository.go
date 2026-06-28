package repository

import (
	"context"
	"errors"
	"ironflow/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FichaTreinoRepository struct {
	DB *pgxpool.Pool
}

func NovoFichaTreinoRepository(db *pgxpool.Pool) *FichaTreinoRepository {
	return &FichaTreinoRepository{DB: db}
}

func (r *FichaTreinoRepository) Salvar(c context.Context, fichaTreino *model.FichaTreino) error {

	sql := `INSERT INTO treino.fit_ficha_treino (tre_nr_id, exe_nr_id, fit_nr_ordem, fit_nr_meta_series, fit_tx_meta_repeticoes, fit_nr_meta_peso, fit_nr_grupo, exe_tx_tipo_equipamento, fit_bl_dropset)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING fit_nr_id,created_at,updated_at`


	  ctx := context.Background()
	err := r.DB.QueryRow(ctx, sql,
		&fichaTreino.TreNrID,
		&fichaTreino.ExeNrID,
		&fichaTreino.FitNrOrdem,
		&fichaTreino.FitNrMetaSeries,
		&fichaTreino.FitTxMetaRepeticoes,
		&fichaTreino.FitNrMetaPeso,
		&fichaTreino.FitNrGrupo,
		&fichaTreino.ExeTxTipoEquipamento,
		&fichaTreino.FitBlDropSet,
	).Scan(
		&fichaTreino.FitNrID,
		&fichaTreino.CreatedAt,
		&fichaTreino.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil

}

func (r *FichaTreinoRepository) Editar(c context.Context, fichaTreino *model.FichaTreino) error {

	sql :=
	`UPDATE treino.fit_ficha_treino 
	SET fit_nr_ordem = $2,
	fit_nr_meta_series = $3, 
	fit_tx_meta_repeticoes = $4, 
	fit_nr_meta_peso = $5, 
	fit_nr_grupo = $6,
	exe_tx_tipo_equipamento = $7,
	fit_bl_dropset = $8,
	updated_at = NOW() 
	WHERE fit_nr_id = $1 and deleted_at IS NULL
	RETURNING created_at,updated_at`

	err := r.DB.QueryRow(c, sql,
		fichaTreino.FitNrID,
		fichaTreino.FitNrOrdem,
		fichaTreino.FitNrMetaSeries,
		fichaTreino.FitTxMetaRepeticoes,
		fichaTreino.FitNrMetaPeso,
		fichaTreino.FitNrGrupo,
		fichaTreino.ExeTxTipoEquipamento,
		fichaTreino.FitBlDropSet,
	).Scan(
		&fichaTreino.CreatedAt,
		&fichaTreino.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("Não é possível editar: Ficha de treino inexistente")
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *FichaTreinoRepository) BuscarPorID(c context.Context, fitNrID int, usuTxID string) (*model.FichaTreinoResponse, error) {
	sql := `
	SELECT f.fit_nr_id, f.tre_nr_id, f.exe_nr_id, e.exe_tx_nome, f.fit_nr_ordem, f.fit_nr_meta_series, f.fit_tx_meta_repeticoes, f.fit_nr_meta_peso, f.fit_nr_grupo, f.exe_tx_tipo_equipamento, fit_bl_dropset, f.created_at, f.updated_at
	FROM treino.fit_ficha_treino f
	JOIN treino.exe_exercicio e ON f.exe_nr_id = e.exe_nr_id
	JOIN treino.tre_treino t ON f.tre_nr_id = t.tre_nr_id
	WHERE f.deleted_at IS NULL 
	AND f.fit_nr_id = $1
	AND t.usu_tx_id = $2
	`
	var ficha model.FichaTreinoResponse
	err := r.DB.QueryRow(c, sql, fitNrID, usuTxID).Scan(
		&ficha.FitNrID,
		&ficha.TreNrID,
		&ficha.ExeNrID,
		&ficha.ExeTxNome,
		&ficha.FitNrOrdem,
		&ficha.FitNrMetaSeries,
		&ficha.FitTxMetaRepeticoes,
		&ficha.FitNrMetaPeso,
		&ficha.FitNrGrupo,
		&ficha.ExeTxTipoEquipamento,
		&ficha.FitBlDropSet,
		&ficha.CreatedAt,
		&ficha.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ficha, nil
}

func (r *FichaTreinoRepository) BuscarTodos(c context.Context, treNrID int, exeTxNome string, usuTxId string) ([]model.FichaTreinoResponse, error) {
	sql := `
	SELECT f.fit_nr_id, f.tre_nr_id, f.exe_nr_id, e.exe_tx_nome, f.fit_nr_ordem, f.fit_nr_meta_series, f.fit_tx_meta_repeticoes, f.fit_nr_meta_peso, f.fit_nr_grupo, f.exe_tx_tipo_equipamento, fit_bl_dropset, f.created_at, f.updated_at
	FROM treino.fit_ficha_treino f
	JOIN treino.exe_exercicio e ON f.exe_nr_id = e.exe_nr_id
	JOIN treino.tre_treino t ON f.tre_nr_id = t.tre_nr_id
	WHERE f.deleted_at IS NULL AND f.tre_nr_id = $1 AND (e.exe_tx_nome <> '' AND e.exe_tx_nome ILIKE $2)
	AND t.usu_tx_id = $3
	ORDER BY f.fit_nr_ordem
	`
	rows, err := r.DB.Query(c, sql, treNrID, "%"+exeTxNome+"%", usuTxId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fichas []model.FichaTreinoResponse
	for rows.Next() {
		var ficha model.FichaTreinoResponse
		if err := rows.Scan(
			&ficha.FitNrID,
			&ficha.TreNrID,
			&ficha.ExeNrID,
			&ficha.ExeTxNome,
			&ficha.FitNrOrdem,
			&ficha.FitNrMetaSeries,
			&ficha.FitTxMetaRepeticoes,
			&ficha.FitNrMetaPeso,
			&ficha.FitNrGrupo,
			&ficha.ExeTxTipoEquipamento,
			&ficha.FitBlDropSet,
			&ficha.CreatedAt,
			&ficha.UpdatedAt,
		); err != nil {
			return nil, err
		}
		fichas = append(fichas, ficha)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fichas, nil

}

func (r *FichaTreinoRepository) Deletar(ctx context.Context, fitNrID int, usuTxId string) error {

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	sql := `UPDATE treino.fit_ficha_treino f
	SET deleted_at = NOW()
	FROM treino.tre_treino t
	WHERE f.tre_nr_id = t.tre_nr_id
	AND f.fit_nr_id = $1
	AND f.deleted_at IS NULL
	AND t.usu_tx_id = $2`

	result, err := tx.Exec(ctx, sql, fitNrID, usuTxId)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("Não é possível deletar: Ficha de treino inexistente, já deletada ou sem permissão")
	}

	return tx.Commit(ctx)

}

func (r *FichaTreinoRepository) ExisteExercicioNoTreino(ctx context.Context, treNrId int, exeNrId int) (bool, error) {

	sql := `
	SELECT EXISTS (
	SELECT 1	
	FROM treino.fit_ficha_treino f
	WHERE f.deleted_at IS NULL 
	AND f.tre_nr_id = $1
	AND f.exe_nr_id = $2
	)
  `
	var exist bool

	err := r.DB.QueryRow(ctx, sql, treNrId, exeNrId).Scan(
		&exist,
	)
	if err != nil {
		return false, err
	}
	return exist, nil

}
