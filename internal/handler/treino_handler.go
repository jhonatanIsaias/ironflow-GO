package handler

import (
	"context"
	"ironflow/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ITreinoRepository interface {
	Salvar(ctx context.Context, t *model.Treino, usuTxId string) error
	Editar(ctx context.Context, t *model.Treino, usuTxId string) error
	BuscarPorID(ctx context.Context, id int, usuTxId string) (*model.Treino, error)
	BuscarTodos(ctx context.Context, treTxNome string, usuTxId string) ([]model.Treino, error)
	DeletarE_Fichas(ctx context.Context, id int, usuTxId string) error
}

type TreinoHandler struct {
	TreinoRepository ITreinoRepository
}

func NovoTreinoHandler(repo ITreinoRepository) *TreinoHandler {
	return &TreinoHandler{TreinoRepository: repo}
}

func (h *TreinoHandler) CriarTreino(c *gin.Context) {
	var t model.Treino

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}
	if t.TreNrID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Não envie o ID para criar um novo treino"})
		return
	}

	usuTxId := c.GetString("usuTxId")

	err := h.TreinoRepository.Salvar(c, &t, usuTxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar o treino"})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func (h *TreinoHandler) EditarTreino(c *gin.Context) {
	var t model.Treino

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}

	if( t.TreNrID == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID do treino é obrigatório para edição"})
		return
	}

	usuTxId := c.GetString("usuTxId")

	err := h.TreinoRepository.Editar(c, &t, usuTxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao editar o treino"})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *TreinoHandler) BuscarPorID(c *gin.Context) {

	treNrId := c.Param("treNrId")

	id, err := strconv.Atoi(treNrId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}

	usuTxId := c.GetString("usuTxId")

	treino, err := h.TreinoRepository.BuscarPorID(c, id, usuTxId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, treino)
}

func (h *TreinoHandler) BuscarTodos(c *gin.Context) {

	treTxNome := c.Query("treTxNome")
	
	usuTxId := c.GetString("usuTxId")

	treinos, err := h.TreinoRepository.BuscarTodos(c, treTxNome, usuTxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar treinos"})
		return
	}
	c.JSON(http.StatusOK, treinos)
}

func (h *TreinoHandler) DeletarPorID(c *gin.Context) {

	treNrId := c.Param("treNrId")

	id, err := strconv.Atoi(treNrId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}
	
	usuTxId := c.GetString("usuTxId")

	err = h.TreinoRepository.DeletarE_Fichas(c, id, usuTxId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Treino não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Treino deletado com sucesso"})
}
