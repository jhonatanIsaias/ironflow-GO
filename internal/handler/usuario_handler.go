package handler

import (
	"context"
	"ironflow/internal/model"
	"ironflow/internal/security"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IUsuarioRepository interface {
	Salvar(ctx context.Context, usuario *model.UsuarioRequest) error
	Editar(ctx context.Context, usuario *model.Usuario) error
	BuscarPorEmail(ctx context.Context, usuTxEmail string) (*model.Usuario, error)
	BuscarPorID(ctx context.Context, usuTxId string) (*model.Usuario, error)
	BuscarPorRefreshToken(ctx context.Context, refreshToken string) (*model.Usuario, error)
}

type UsuarioHandler struct {
	usuarioRepository IUsuarioRepository
}

func NovoUsuarioHandler(usuarioRepo IUsuarioRepository) *UsuarioHandler {
	return &UsuarioHandler{usuarioRepository: usuarioRepo}
}

func (h *UsuarioHandler) SalvarUsuario(c *gin.Context) {

	var usuarioRequest model.UsuarioRequest
	if err := c.ShouldBindJSON(&usuarioRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido"})
		return
	}

	usuTxSenhaHash, err := security.HashPassword(usuarioRequest.UsuTxSenha)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criptograr senha"})
		return
	}

	usuarioRequest.UsuTxSenha = usuTxSenhaHash

	err = h.usuarioRepository.Salvar(c, &usuarioRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var usuarioResponse model.UsuarioResponse
	usuarioResponse.UsuTxNome = usuarioRequest.UsuTxNome
	usuarioResponse.UsuTxEmail = usuarioRequest.UsuTxEmail
	usuarioResponse.CreatedAt = usuarioRequest.CreatedAt
	usuarioResponse.UpdatedAt = usuarioRequest.UpdatedAt

	c.JSON(http.StatusCreated, usuarioResponse)

}

func (h *UsuarioHandler) EditarUsuario(c *gin.Context) {

	var usuario model.UsuarioResponse

	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido"})
		return
	}

	usuTxId := c.GetString("usuTxId")
	usuarioToEdit, err := h.usuarioRepository.BuscarPorID(c, usuTxId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	err = h.usuarioRepository.Editar(c, usuarioToEdit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao editar usuário"})
		return
	}

	usuario.UsuTxNome = usuarioToEdit.UsuTxNome
	usuario.UsuTxEmail = usuarioToEdit.UsuTxEmail
	usuario.CreatedAt = usuarioToEdit.CreatedAt
	usuario.UpdatedAt = usuarioToEdit.UpdatedAt

	c.JSON(http.StatusOK, usuario)
}

func (h *UsuarioHandler) Login(c *gin.Context) {

	var JWTRequest model.JWTRequest

	if err := c.ShouldBindJSON(&JWTRequest); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Corpo da requisição inválido"})
		return
	}
	usuario, err := h.usuarioRepository.BuscarPorEmail(c, JWTRequest.UsuTxEmail)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email ou senha Inválidos"})
		return
	}

	if !security.CheckPasswordHash(JWTRequest.UsuTxSenha, usuario.UsuTxSenha) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email ou senha Inválidos"})
		return
	}

	token, err := security.GenerateJWT(usuario.UsuTxId, usuario.UsuTxEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao gerar JWT"})
		return
	}

	refreshToken, err := security.GerarRefreshToken()

	if err != nil {
		log.Print("error:" + err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao gerar refresh token"})
		return
	}

	refreshTokenExp := time.Now().Add(time.Hour * 24 * 7)

	usuario.UsuTxRefreshToken = &refreshToken
	usuario.UsuDtRefreshTokenExp = &refreshTokenExp

	err = h.usuarioRepository.Editar(c, usuario)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao gerar refresh token"})
		return
	}

	c.JSON(http.StatusOK, model.JWTResponse{
		JWTToken:          token,
		UsuTxRefreshToken: usuario.UsuTxRefreshToken,
		UsuTxNome:         usuario.UsuTxNome,
	})

}

func (h *UsuarioHandler) Refresh(c *gin.Context) {

	var JWTRefresh model.JWTRefresh

	if err := c.ShouldBindJSON(&JWTRefresh); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Corpo da requisição inválido"})
		return
	}

	if JWTRefresh.UsuTxRefreshToken == nil || *JWTRefresh.UsuTxRefreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token não fornecido"})
		return
	}

	ctx := c.Request.Context()

	usuario, err := h.usuarioRepository.BuscarPorRefreshToken(ctx, *JWTRefresh.UsuTxRefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não encontrado"})
		return
	}

	refreshExpirado := usuario.UsuDtRefreshTokenExp.Before(time.Now())

	if refreshExpirado {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expirado"})
		return
		
	}

	token, err := security.GenerateJWT(usuario.UsuTxId, usuario.UsuTxEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar JWT"})
		return
	}

	c.JSON(http.StatusOK, model.JWTResponse{
		JWTToken:          token,
		UsuTxRefreshToken: usuario.UsuTxRefreshToken,
		UsuTxNome:         usuario.UsuTxNome,
	})

}
