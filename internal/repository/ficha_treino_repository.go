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
	
	sql := `INSERT INTO treino.fit_ficha_treino (tre_nr_id, exe_nr_id, fit_nr_ordem, fit_nr_meta_series, fit_nr_meta_repeticoes, fit_nr_meta_peso)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING fit_nr_id,created_at,updated_at`

	err := r.DB.QueryRow(c, sql, 
		fichaTreino.TreNrID, 
		fichaTreino.ExeNrID, 
		fichaTreino.FitNrOrdem, 
		fichaTreino.FitNrMetaSeries, 
		fichaTreino.FitNrMetaRepeticoes, 
		fichaTreino.FitNrMetaPeso,
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
	SET tre_nr_id = $1, 
	exe_nr_id = $2, 
	fit_nr_ordem = $3,
	fit_nr_meta_series = $4, 
	fit_nr_meta_repeticoes = $5, 
	fit_nr_meta_peso = $6, 
	updated_at = NOW() 
	WHERE fit_nr_id = $7 and deleted_at IS NULL
	RETURNING created_at,updated_at`

	err := r.DB.QueryRow(c, sql,
		fichaTreino.TreNrID, 
		fichaTreino.ExeNrID,
		fichaTreino.FitNrOrdem,
		fichaTreino.FitNrMetaSeries,
		fichaTreino.FitNrMetaRepeticoes,
		fichaTreino.FitNrMetaPeso,
		fichaTreino.FitNrID,
	).Scan(
		&fichaTreino.CreatedAt,
		&fichaTreino.UpdatedAt,
	)

	if(errors.Is(err, pgx.ErrNoRows)) {
		return errors.New("Não é possível editar: Ficha de treino inexistente")
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *FichaTreinoRepository) BuscarPorID(c context.Context, fitNrID int) (*model.FichaTreinoResponse, error) {
	sql := `
	SELECT f.fit_nr_id, f.tre_nr_id, f.exe_nr_id, e.exe_tx_nome, f.fit_nr_ordem, f.fit_nr_meta_series, f.fit_nr_meta_repeticoes, f.fit_nr_meta_peso
	FROM treino.fit_ficha_treino f
	JOIN treino.exe_exercicio e ON f.exe_nr_id = e.exe_nr_id
	WHERE f.deleted_at IS NULL AND f.fit_nr_id = $1
	`
	var ficha model.FichaTreinoResponse
	err := r.DB.QueryRow(c, sql, fitNrID).Scan(
		&ficha.FitNrID,
		&ficha.TreNrID,
		&ficha.ExeNrID,
		&ficha.ExeTxNome,
		&ficha.FitNrOrdem,
		&ficha.FitNrMetaSeries,
		&ficha.FitNrMetaRepeticoes,
		&ficha.FitNrMetaPeso,
	)
	if err != nil {
		return nil, err
	}
	return &ficha, nil
}


func (r *FichaTreinoRepository) BuscarTodos(c context.Context, treNrID int,exeTxNome string) ([]model.FichaTreinoResponse, error) {
	sql := `
	SELECT f.fit_nr_id, f.tre_nr_id, f.exe_nr_id, e.exe_tx_nome, f.fit_nr_ordem, f.fit_nr_meta_series, f.fit_nr_meta_repeticoes, f.fit_nr_meta_peso
	FROM treino.fit_ficha_treino f
	JOIN treino.exe_exercicio e ON f.exe_nr_id = e.exe_nr_id
	AND f.deleted_at IS NULL AND f.tre_nr_id = $1 AND (e.exe_tx_nome <> '' AND e.exe_tx_nome ILIKE $2)
	ORDER BY f.fit_nr_ordem
	`
	rows, err := r.DB.Query(c, sql, treNrID, "%"+exeTxNome+"%")
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
			&ficha.FitNrMetaRepeticoes,
			&ficha.FitNrMetaPeso,
		); err != nil {
			return nil, err
		}
		fichas = append(fichas, ficha)
	}

	return fichas, nil

}

func (r *FichaTreinoRepository) Deletar(c context.Context, fitNrID int) error {
	
	sql := `UPDATE treino.fit_ficha_treino SET deleted_at = NOW() WHERE fit_nr_id = $1 AND deleted_at IS NULL`
	result, err := r.DB.Exec(c, sql, fitNrID)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("Não é possível deletar: Ficha de treino inexistente ou já deletada")
	}
	return nil
}

func (r *FichaTreinoRepository) ExisteExercicioNoTreino(ctx context.Context,treNrId int, exeNrId int) (bool,error){

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