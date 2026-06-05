package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ironflow/internal/model"
)

type TreinoRepository struct {
	DB *pgxpool.Pool
}

func NovoTreinoRepository(db *pgxpool.Pool) *TreinoRepository {
	return &TreinoRepository{DB: db}
}

func (r *TreinoRepository) Salvar(ctx context.Context, t *model.Treino) error {
	sql := `
        INSERT INTO treino.tre_treino (tre_tx_nome)
        VALUES ($1) RETURNING tre_nr_id,created_at, updated_at`

	err := r.DB.QueryRow(ctx, sql, t.TreTxNome).Scan(
		&t.TreNrID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	return err
}

func (r *TreinoRepository) Editar(ctx context.Context, t *model.Treino) error {
	sql := `
		UPDATE treino.tre_treino
		SET tre_tx_nome = $1, updated_at = NOW()
		WHERE tre_nr_id = $2 AND deleted_at IS NULL
		RETURNING created_at, updated_at
	`
	err := r.DB.QueryRow(ctx, sql,
		t.TreTxNome,
		t.TreNrID,
	).Scan(
		&t.CreatedAt,
		&t.UpdatedAt,
	)


	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("Não é possível editar: Treino inexistente")
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *TreinoRepository) BuscarTodos(ctx context.Context,treTxNome string) ([]model.Treino, error) {
	sql := `
        SELECT tre_nr_id, tre_tx_nome
        FROM treino.tre_treino
        WHERE deleted_at IS NULL AND (tre_tx_nome <> '' AND tre_tx_nome ILIKE $1)
    `

	rows, err := r.DB.Query(ctx, sql, "%"+treTxNome+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var treinos []model.Treino
	for rows.Next() {
		var t model.Treino
		err := rows.Scan(
			&t.TreNrID,
			&t.TreTxNome,
		)
		if err != nil {
			return nil, err
		}
		treinos = append(treinos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return treinos, nil
}

func (r *TreinoRepository) BuscarPorID(ctx context.Context, treNrID int) (*model.Treino, error) {
	sql := `
        SELECT tre_nr_id, tre_tx_nome
        FROM treino.tre_treino
        WHERE tre_nr_id = $1 AND deleted_at IS NULL
		and
    `

	var t model.Treino
	err := r.DB.QueryRow(ctx, sql, treNrID).Scan(
		&t.TreNrID,
		&t.TreTxNome,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("Treino não encontrado")
	} else if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TreinoRepository) DeletarE_Fichas(ctx context.Context, id int) error {

	tx,err := r.DB.Begin(ctx);

	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}

	defer tx.Rollback(ctx)


	sql := `UPDATE treino.tre_treino SET deleted_at = NOW() WHERE tre_nr_id = $1`
	comando, err := tx.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if comando.RowsAffected() == 0 {
		return errors.New("Não é possível deletar: Treino inexistente")
	}

	sqlFichas := `UPDATE treino.fit_ficha_treino SET deleted_at = NOW() WHERE tre_nr_id = $1 AND deleted_at IS NULL`
	_, err = tx.Exec(ctx, sqlFichas, id)
	
	if err != nil {
		return err 
	}
	return tx.Commit(ctx)
}
