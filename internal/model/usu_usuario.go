package model

import (
	"time"

	"github.com/google/uuid"
)

type Usuario struct {
	BaseEntity
	UsuTxId    uuid.UUID `db:"usu_tx_id"`
	UsuTxNome  string    `db:"usu_tx_nome"`
	UsuTxEmail string    `db:"usu_tx_email"`
	UsuTxSenha string    `db:"usu_tx_senha"`
	UsuTxRefreshToken *string `db:"usu_tx_refresh_token"`
	UsuDtRefreshTokenExp *time.Time `db:"usu_dt_refresh_token_exp"`
}

type UsuarioRequest struct {
	BaseEntity
	UsuTxNome  string `json:"usuTxNome" binding:"required"`
	UsuTxEmail string `json:"usuTxEmail" binding:"required,email"`
	UsuTxSenha string `json:"usuTxSenha" binding:"required min=8,max=12"`
}

type UsuarioResponse struct {
	BaseEntity
	UsuTxNome  string    `json:"usuTxNome"`
	UsuTxEmail string    `json:"usuTxEmail"`
}


type JWTRequest struct {
	UsuTxEmail string `json:"usuTxEmail" binding:"required,email"`
	UsuTxSenha string `json:"usuTxSenha" binding:"required"`
}

type JWTResponse struct {
	JWTToken   string    `json:"jwtToken"`
	UsuTxNome  string    `json:"usuTxNome"`
	UsuTxRefreshToken *string `json:"usuTxRefreshToken"`
}

type JWTRefresh struct {
	UsuTxRefreshToken *string `json:"usuTxRefreshToken" binding:"required"`
}