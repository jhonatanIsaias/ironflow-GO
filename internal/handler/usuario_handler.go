package handler

import (
	"context"
	"ironflow/internal/model"
	"ironflow/internal/security"
	"net/http"

	"github.com/gin-gonic/gin"
	
)

type IUsuarioRepository interface {
	Salvar(ctx context.Context, usuario *model.UsuarioRequest) error
	Editar(ctx context.Context, usuario *model.Usuario) error
	BuscarPorEmail(ctx context.Context, usuTxEmail string) (*model.Usuario, error)
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
		c.JSON(http.DefaultMaxHeaderBytes, gin.H{"error": "Failed to bind request"})
		return
	}

	usuTxSenhaHash,err := security.HashPassword(usuarioRequest.UsuTxSenha)
	
	if err != nil {
		c.JSON(http.DefaultMaxHeaderBytes, gin.H{"error": "Failed to hash password"})
		return
	}

	usuarioRequest.UsuTxSenha = usuTxSenhaHash

	err = h.usuarioRepository.Salvar(c, &usuarioRequest)
	if err != nil {
		 c.JSON(http.DefaultMaxHeaderBytes, gin.H{"error": "Failed to save user"})
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
		c.JSON(http.DefaultMaxHeaderBytes, gin.H{"error": "Failed to bind request"})
		return
	}

	 usuarioToEdit, err := h.usuarioRepository.BuscarPorEmail(c, usuario.UsuTxEmail)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	err = h.usuarioRepository.Editar(c,usuarioToEdit)

	if err != nil {
		c.JSON(http.DefaultMaxHeaderBytes, gin.H{"error": "Failed to edit user"})
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
		c.JSON(http.DefaultMaxHeaderBytes, gin.H{"error": "Failed to bind request"})
		return
	}
	usuario, err := h.usuarioRepository.BuscarPorEmail(c, JWTRequest.UsuTxEmail)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !security.CheckPasswordHash(JWTRequest.UsuTxSenha, usuario.UsuTxSenha) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := security.GenerateJWT(usuario.UsuTxId, usuario.UsuTxEmail)
	if err != nil {
		c.JSON(http.DefaultMaxHeaderBytes, gin.H{"error": "Failed to generate JWT"})
		return
	}	

	c.JSON(http.StatusOK, model.JWTResponse{
		JWTToken: token,
		UsuTxNome: usuario.UsuTxNome,
	})

}