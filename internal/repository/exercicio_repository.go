package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"ironflow/internal/model"
)

type ExercicioRepository struct {
	DB *pgxpool.Pool
	redis *redis.Client
}

func NovoExercicioRepository(db *pgxpool.Pool, redis *redis.Client) *ExercicioRepository {
	return &ExercicioRepository{DB: db, redis: redis}
}

func (r *ExercicioRepository) Salvar(ctx context.Context, e *model.Exercicio, usuTxID string) error {
	sql := `
		INSERT INTO treino.exe_exercicio (exe_tx_nome, exe_tx_grupo_muscular, exe_tx_grupo_muscular_sinegista,usu_tx_id)
		VALUES ($1, $2, $3, $4) RETURNING exe_nr_id, created_at, updated_at`

	err := r.DB.QueryRow(ctx, sql,
		e.ExeTxNome,
		e.ExeTxGrupoMuscular,
		e.ExeTxGrupoMuscularSinergista,
		usuTxID,
	).Scan(
		&e.ExeNrID,
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	return err

}

func (r *ExercicioRepository) Editar(ctx context.Context, e *model.Exercicio, usuTxID string) error {

	sql :=
		`
		UPDATE treino.exe_exercicio
		SET exe_tx_nome = $1, 
		exe_tx_grupo_muscular = $2, 
		exe_tx_grupo_muscular_sinegista = $3,
		updated_at = NOW()
		WHERE exe_nr_id = $4 AND deleted_at IS NULL AND usu_tx_id = $5
		RETURNING created_at, updated_at
	`
	err := r.DB.QueryRow(ctx, sql,
		e.ExeTxNome,
		e.ExeTxGrupoMuscular,
		e.ExeTxGrupoMuscularSinergista,
		e.ExeNrID,
		usuTxID,
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

func (r *ExercicioRepository) BuscarTodos(ctx context.Context, usuTxID string) ([]model.Exercicio, error) {

	exerciciosGerais,_ := r.buscarTodosRedis(ctx)
	
	sql := `
        SELECT exe_nr_id, exe_tx_nome, exe_tx_grupo_muscular, usu_tx_id, exe_tx_grupo_muscular_sinegista, created_at, updated_at
        FROM treino.exe_exercicio
        WHERE usu_tx_id = $1::UUID AND deleted_at IS NULL
    `

	rows, err := r.DB.Query(ctx, sql, usuTxID)
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
			&e.UsuTxID,
			&e.ExeTxGrupoMuscularSinergista,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		e.IsCustom = e.UsuTxID != nil
		exercicios = append(exercicios, e)
	}

	exercicios = append(exercicios, exerciciosGerais...)

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return exercicios, nil
}

func (r *ExercicioRepository) buscarTodosRedis(ctx context.Context) ([]model.Exercicio, error){

	var exerciciosGerais []model.Exercicio

	cacheKey := "exercicios_sistema_default"
	
	val, err := r.redis.Get(ctx, cacheKey).Result()

	switch err {
		case redis.Nil:
		
		sqlDefault := `
			SELECT exe_nr_id, exe_tx_nome, exe_tx_grupo_muscular, usu_tx_id, exe_tx_grupo_muscular_sinegista, created_at, updated_at
			FROM treino.exe_exercicio
			WHERE usu_tx_id IS NULL AND deleted_at IS NULL
		`
		rows, errDB := r.DB.Query(ctx, sqlDefault)
		if errDB != nil {
			return nil, errDB
		}
		
		for rows.Next() {
			var e model.Exercicio
			errScan := rows.Scan(&e.ExeNrID, &e.ExeTxNome, &e.ExeTxGrupoMuscular, &e.UsuTxID, &e.ExeTxGrupoMuscularSinergista, &e.CreatedAt, &e.UpdatedAt)
			if errScan == nil {
				e.IsCustom = false 
				exerciciosGerais = append(exerciciosGerais, e)
			}
		}
		rows.Close() 

		
		jsonBytes, _ := json.Marshal(exerciciosGerais)
		r.redis.Set(ctx, cacheKey, jsonBytes, 24*time.Hour)

	case nil:
		json.Unmarshal([]byte(val), &exerciciosGerais)
	} 

	return exerciciosGerais, nil
}

func (r *ExercicioRepository) BuscarPorID(ctx context.Context, exeNrId int, usuTxID string) (*model.Exercicio, error) {

	sql := `
		SELECT exe_nr_id, exe_tx_nome, exe_tx_grupo_muscular, exe_tx_grupo_muscular_sinegista, created_at, updated_at
		FROM treino.exe_exercicio
		WHERE exe_nr_id = $1 AND ( usu_tx_id IS NULL OR usu_tx_id = $2) AND deleted_at IS NULL
	`
	var e model.Exercicio

	err := r.DB.QueryRow(ctx, sql, exeNrId, usuTxID).Scan(
		&e.ExeNrID,
		&e.ExeTxNome,
		&e.ExeTxGrupoMuscular,
		&e.ExeTxGrupoMuscularSinergista,
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

func (r *ExercicioRepository) Deletar(ctx context.Context, id int, usuTxID string) error {

	sql := `UPDATE treino.exe_exercicio SET deleted_at = NOW() WHERE exe_nr_id = $1 AND usu_tx_id = $2`

	comando, err := r.DB.Exec(ctx, sql, id, usuTxID)
	if err != nil {
		return err
	}

	if comando.RowsAffected() == 0 {
		return errors.New("Não é possível deletar: Exercício inexistente")
	}

	return nil
}
