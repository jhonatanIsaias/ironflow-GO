package repository

import (
	"context"
	"ironflow/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsuarioRepository struct {
	DB *pgxpool.Pool
}

func NovoUsuarioRepository(db *pgxpool.Pool) *UsuarioRepository {
	return &UsuarioRepository{DB: db}
}

func (r *UsuarioRepository) Salvar(ctx context.Context, usuario *model.UsuarioRequest) error {
	sql := `INSERT INTO usuarios (usu_tx_nome, usu_tx_email, usu_tx_senha) VALUES ($1, $2, $3) RETURNING usu_tx_nome,usu_tx_email, created_at, updated_at`

	err := r.DB.QueryRow(ctx, sql, usuario.UsuTxNome, usuario.UsuTxEmail, usuario.UsuTxSenha).Scan(
		&usuario.UsuTxNome,
		&usuario.UsuTxEmail,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)
	return err
}

func (r *UsuarioRepository) Editar(ctx context.Context, usuario *model.Usuario) error {
	sql := `UPDATE usuarios SET usu_tx_nome = $1, 
	usu_tx_email = $2, updated_at = NOW() 
	WHERE usu_tx_id = $3 AND deleted_at IS NULL RETURNING usu_tx_nome,usu_tx_email, created_at, updated_at`

	err := r.DB.QueryRow(ctx, sql, usuario.UsuTxNome, usuario.UsuTxEmail, usuario.UsuTxId).Scan(
		&usuario.UsuTxNome,
		&usuario.UsuTxEmail,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)
	return err
}

func (r *UsuarioRepository) BuscarPorEmail(ctx context.Context, usuTxEmail string) (*model.Usuario, error) {
	sql := `SELECT usu_tx_id, usu_tx_nome, usu_tx_email, usu_tx_senha FROM usuarios WHERE usu_tx_email = $1 AND deleted_at IS NULL`
	var usuario model.Usuario
	err := r.DB.QueryRow(ctx, sql, usuTxEmail).Scan(
		&usuario.UsuTxId,
		&usuario.UsuTxNome,
		&usuario.UsuTxEmail,
		&usuario.UsuTxSenha,
	)
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}