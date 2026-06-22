package repository

import (
	"context"
	"errors"
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
	existsSQL := `SELECT 1 FROM auth.usu_usuario WHERE usu_tx_email = $1 AND deleted_at IS NULL`
	var exists int
	err := r.DB.QueryRow(ctx, existsSQL, usuario.UsuTxEmail).Scan(&exists)
	if err == nil {
		return errors.New("email já cadastrado")
	}

	sql := `INSERT INTO auth.usu_usuario (usu_tx_nome, usu_tx_email, usu_tx_senha) VALUES ($1, $2, $3) RETURNING usu_tx_nome,usu_tx_email, created_at, updated_at`

	err = r.DB.QueryRow(ctx, sql, usuario.UsuTxNome, usuario.UsuTxEmail, usuario.UsuTxSenha).Scan(
		&usuario.UsuTxNome,
		&usuario.UsuTxEmail,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)
	return err
}

func (r *UsuarioRepository) Editar(ctx context.Context, usuario *model.Usuario) error {
	sql := `UPDATE auth.usu_usuario SET usu_tx_nome = $1, 
	usu_tx_email = $2, 
	updated_at = NOW(),
	usu_tx_refresh_token = $4,
	usu_dt_refresh_token_exp = $5
	WHERE usu_tx_id = $3 AND deleted_at IS NULL RETURNING usu_tx_nome,usu_tx_email, created_at, updated_at`

	err := r.DB.QueryRow(ctx, sql, usuario.UsuTxNome, usuario.UsuTxEmail, usuario.UsuTxId, usuario.UsuTxRefreshToken, usuario.UsuDtRefreshTokenExp).Scan(
		&usuario.UsuTxNome,
		&usuario.UsuTxEmail,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)
	return err
}

func (r *UsuarioRepository) BuscarPorEmail(ctx context.Context, usuTxEmail string) (*model.Usuario, error) {
	sql := `SELECT usu_tx_id, usu_tx_nome, usu_tx_email, usu_tx_senha FROM auth.usu_usuario WHERE usu_tx_email = $1 AND deleted_at IS NULL`
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

func (r *UsuarioRepository) BuscarPorID(ctx context.Context, usuTxId string) (*model.Usuario, error) {
	sql := `SELECT usu_tx_id, usu_tx_nome, usu_tx_email, usu_tx_senha FROM auth.usu_usuario WHERE usu_tx_id = $1 AND deleted_at IS NULL`
	var usuario model.Usuario
	err := r.DB.QueryRow(ctx, sql, usuTxId).Scan(
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

func (r *UsuarioRepository) BuscarPorRefreshToken(ctx context.Context, refreshToken string) (*model.Usuario, error) {
	sql := `SELECT usu_tx_id, usu_tx_nome, usu_tx_email, usu_tx_senha, usu_tx_refresh_token, usu_dt_refresh_token_exp FROM auth.usu_usuario WHERE usu_tx_refresh_token = $1 AND deleted_at IS NULL`
	var usuario model.Usuario
	err := r.DB.QueryRow(ctx, sql, refreshToken).Scan(
		&usuario.UsuTxId,
		&usuario.UsuTxNome,
		&usuario.UsuTxEmail,
		&usuario.UsuTxSenha,
		&usuario.UsuTxRefreshToken,
		&usuario.UsuDtRefreshTokenExp,
	)
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}
