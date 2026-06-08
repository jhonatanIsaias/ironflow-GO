package repository

import (
	"context"
	"errors"
	"fmt"
	"ironflow/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
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
	usuTxId string, 
	dataInicio time.Time,
	dataFim time.Time,
	horaInicio time.Time,
	horaFim time.Time) ([]model.SessaoTreinoDetalhada, error) {

	sql := `
		SELECT set_nr_id, tre_nr_id, set_dt_data, set_tm_hora_inicio, tre.tre_tx_nome, created_at, updated_at
		FROM treino.set_sessao_treino set
		INNER JOIN treino.tre_treino tre ON set.tre_nr_id = tre.tre_nr_id
		WHERE tre_nr_id = $1 AND deleted_at IS NULL
		AND tre.usu_tx_id = $2
	`

	args := []any{treNrId,usuTxId}

	if !dataInicio.IsZero() {
		args = append(args, dataInicio)
	
		sql += fmt.Sprintf(" AND set_dt_data >= $%d", len(args))
	}

	if !dataFim.IsZero() {
		args = append(args, dataFim)
		sql += fmt.Sprintf(" AND set_dt_data <= $%d", len(args))
	}

	if !horaInicio.IsZero() {
		args = append(args, horaInicio)
		sql += fmt.Sprintf(" AND set_tm_hora_inicio >= $%d", len(args))
	}

	if !horaFim.IsZero() {
		args = append(args, horaFim)
		sql += fmt.Sprintf(" AND set_tm_hora_inicio <= $%d", len(args))
	}

	sql += " ORDER BY set_dt_data DESC, set_tm_hora_inicio DESC"


	rows, err := r.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessoes []model.SessaoTreinoDetalhada
	for rows.Next() {
		var sessao model.SessaoTreinoDetalhada
		if err := rows.Scan(
			&sessao.SetNrID,
			&sessao.TreNrID,
			&sessao.SetDtData,
			&sessao.SetTmHoraInicio,
			&sessao.TreTxNome,
			&sessao.CreatedAt,
			&sessao.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessoes = append(sessoes, sessao)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessoes, nil
}

func (r *SessaoTreinoRepository) ObterSessaoHoje(ctx context.Context, treNrId int, usuTxId string) (int, bool, error) {
	sql := `
		SELECT set_nr_id 
		FROM treino.set_sessao_treino set
		INNER JOIN treino.tre_treino tre ON set.tre_nr_id = tre.tre_nr_id
		WHERE set.tre_nr_id = $1 
		  AND set.set_dt_data = CURRENT_DATE 
		  AND set.deleted_at IS NULL
		  AND tre.usu_tx_id = @2
		LIMIT 1
	`

	var setNrID int
	err := r.DB.QueryRow(ctx, sql, treNrId,usuTxId).Scan(&setNrID)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, err
	}
	return setNrID, true, nil
}