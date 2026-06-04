package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	
	"ironflow/internal/model" 
)

type IExercicioRepository interface {
    Salvar(ctx context.Context, e *model.Exercicio) error
    Editar(ctx context.Context, e *model.Exercicio) error
    BuscarPorID(ctx context.Context, id int) (*model.Exercicio, error)
    Deletar(ctx context.Context, id int) error
}

type ExercicioRepository struct {
	DB *pgxpool.Pool
}

func NovoExercicioRepository(db *pgxpool.Pool) *ExercicioRepository {
	return &ExercicioRepository{DB: db}
}


func (r *ExercicioRepository) Salvar (ctx context.Context, e *model.Exercicio) error{
	sql := `
		INSERT INTO treino.exe_exercicio (exe_tx_nome, exe_tx_grupo_muscular, exe_tx_grupo_muscular_sinegista, exe_tx_tipo_equipamento)
		VALUES ($1, $2, $3, $4) RETURNING exe_nr_id, created_at, updated_at`

		err := r.DB.QueryRow(ctx, sql,
		e.ExeTxNome,
		e.ExeTxGrupoMuscular,
		e.ExeTxGrupoMuscularSinergista,
		e.ExeTxTipoEquipamento,
	).Scan(
		   &e.ExeNrID,
		   &e.CreatedAt,
		   &e.UpdatedAt,
		)

	return err;

}

func (r *ExercicioRepository) Editar(ctx context.Context, e *model.Exercicio) error {
	
	sql := 
	`
		UPDATE treino.exe_exercicio
		SET exe_tx_nome = $1, 
		exe_tx_grupo_muscular = $2, 
		exe_tx_grupo_muscular_sinegista = $3,
		exe_tx_tipo_equipamento = $4,
		updated_at = NOW()
		WHERE exe_nr_id = $5 AND deleted_at IS NULL
		RETURNING created_at, updated_at
	`
	 err := r.DB.QueryRow(ctx, sql,
		e.ExeTxNome,
		e.ExeTxGrupoMuscular,
		e.ExeTxGrupoMuscularSinergista,
		e.ExeTxTipoEquipamento,
		e.ExeNrID,
	).
	Scan(
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("Não é possível editar: Exercício inexistente")
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *ExercicioRepository) BuscarTodos(ctx context.Context)([]model.Exercicio,error){
	sql := `
		SELECT exe_nr_id, exe_tx_nome, exe_tx_grupo_muscular, exe_tx_grupo_muscular_sinegista, exe_tx_tipo_equipamento, created_at, updated_at
		FROM treino.exe_exercicio
		WHERE deleted_at IS NULL
	`

	rows, err := r.DB.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var exercicios []model.Exercicio
	for rows.Next() {
		var e model.Exercicio
		err := rows.Scan(
			&e.ExeNrID,
			&e.ExeTxNome,
			&e.ExeTxGrupoMuscular,
			&e.ExeTxGrupoMuscularSinergista,
			&e.ExeTxTipoEquipamento,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercicios = append(exercicios, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return exercicios, nil
}

func (r *ExercicioRepository) BuscarPorID(ctx context.Context, exeNrId int) (*model.Exercicio, error) {

	sql := `
		SELECT exe_nr_id, exe_tx_nome, exe_tx_grupo_muscular, exe_tx_grupo_muscular_sinegista, exe_tx_tipo_equipamento, created_at, updated_at
		FROM treino.exe_exercicio
		WHERE exe_nr_id = $1 AND deleted_at IS NULL
	`
	var e model.Exercicio

	err := r.DB.QueryRow(ctx, sql, exeNrId).Scan(
		&e.ExeNrID,
		&e.ExeTxNome,
		&e.ExeTxGrupoMuscular,
		&e.ExeTxGrupoMuscularSinergista,
		&e.ExeTxTipoEquipamento,
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("Exercício não encontrado")
	} else if err != nil {
		return nil, err 
	}

	return &e, nil
}


func (r *ExercicioRepository) Deletar(ctx context.Context, id int) error {
	
	sql := `UPDATE treino.exe_exercicio SET deleted_at = NOW() WHERE exe_nr_id = $1`
	
	comando, err := r.DB.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	
	if comando.RowsAffected() == 0 {
		return errors.New("Não é possível deletar: Exercício inexistente")
	}

	return nil
}