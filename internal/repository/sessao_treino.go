package repository

import (
	"context"
	"ironflow/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessaoTreinoRepository struct {
	DB *pgxpool.Pool
}

func NovoSessaoTreinoRepository(db *pgxpool.Pool) *SessaoTreinoRepository {
	return &SessaoTreinoRepository{DB: db}
}

func (r *SessaoTreinoRepository) Salvar(ctx context.Context, sessao *model.SessaoTreino) error {

	sql := `
		INSERT INTO treino.set_sessao_treino (tre_nr_id, set_dt_data, set_tm_hora_inicio)
		VALUES ($1, $2, $3) RETURNING set_nr_id, created_at, updated_at`
	
	err := r.DB.QueryRow(ctx, sql,
		sessao.TreNrID,
		sessao.SetDtData,
		sessao.SetTmHoraInicio,
	).Scan(
		&sessao.SetNrID,
		&sessao.CreatedAt,
		&sessao.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *SessaoTreinoRepository) BuscarPorFiltros(
	ctx context.Context,
	treNrId int,
	dataInicio time.Time,
	dataFim time.Time,
	horaInicio time.Time,
	horaFim time.Time) ([]model.SessaoTreino, error) {
	sql := `
		SELECT set_nr_id, tre_nr_id, set_dt_data, set_tm_hora_inicio, created_at, updated_at
		FROM treino.set_sessao_treino
		WHERE tre_nr_id = $1
		AND $2 IS NULL OR set_dt_data >= $2
		AND $3 IS NULL OR set_dt_data <= $3
		AND $4 IS NULL OR set_tm_hora_inicio >= $4
		AND $5 IS NULL OR set_tm_hora_inicio <= $5
		AND deleted_at IS NULL
	`
	rows, err := r.DB.Query(ctx, sql, treNrId, dataInicio, dataFim, horaInicio, horaFim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessoes []model.SessaoTreino
	
	for rows.Next() {
		var sessao model.SessaoTreino
		if err := rows.Scan(
			&sessao.SetNrID,
			&sessao.TreNrID,
			&sessao.SetDtData,
			&sessao.SetTmHoraInicio,
			&sessao.CreatedAt,
			&sessao.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessoes = append(sessoes, sessao)
	}

	return sessoes, nil
}