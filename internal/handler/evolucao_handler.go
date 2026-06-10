package handler

import (
	"context"
	"ironflow/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IEvolucaoRepository interface {
	Salvar(ctx context.Context, e *model.Evolucao, usuTxID string) error
	Editar(ctx context.Context, e *model.Evolucao, usuTxID string) error
	BuscarPorID(ctx context.Context, evoNrID int, usuTxID string) (*model.Evolucao, error)
	BuscarTodos(ctx context.Context, usuTxID string) ([]model.Evolucao, error)
	BuscarMaisRecente(ctx context.Context, usuTxID string) (*model.Evolucao, error)
	Deletar(ctx context.Context, evoNrID int, usuTxID string) error
}

type EvolucaoHandler struct {
	EvolucaoRepository IEvolucaoRepository
}

func NovoEvolucaoHandler(repo IEvolucaoRepository) *EvolucaoHandler {
	return &EvolucaoHandler{EvolucaoRepository: repo}
}

func (h *EvolucaoHandler) CriarEvolucao(c *gin.Context) {
	var e model.Evolucao

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}

	if e.EvoNrID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Não envie o ID para criar uma nova evolução"})
		return
	}

	usuTxID := c.GetString("usuTxId")

	err := h.EvolucaoRepository.Salvar(c, &e, usuTxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar a evolução"})
		return
	}

	c.JSON(http.StatusCreated, e)
}

func (h *EvolucaoHandler) EditarEvolucao(c *gin.Context) {
	var e model.Evolucao

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}

	if e.EvoNrID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID da evolução é obrigatório para edição"})
		return
	}

	usuTxID := c.GetString("usuTxId")

	err := h.EvolucaoRepository.Editar(c, &e, usuTxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, e)
}

func (h *EvolucaoHandler) BuscarPorID(c *gin.Context) {

	evoNrID := c.Param("evoNrId")

	id, err := strconv.Atoi(evoNrID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}

	usuTxID := c.GetString("usuTxId")

	evolucao, err := h.EvolucaoRepository.BuscarPorID(c, id, usuTxID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, evolucao)
}

func (h *EvolucaoHandler) BuscarTodos(c *gin.Context) {

	usuTxID := c.GetString("usuTxId")

	evolucoes, err := h.EvolucaoRepository.BuscarTodos(c, usuTxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar evoluções"})
		return
	}
	c.JSON(http.StatusOK, evolucoes)
}

func (h *EvolucaoHandler) BuscarMaisRecente(c *gin.Context) {

	usuTxID := c.GetString("usuTxId")

	evolucao, err := h.EvolucaoRepository.BuscarMaisRecente(c, usuTxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, evolucao)
}

func (h *EvolucaoHandler) DeletarPorID(c *gin.Context) {

	evoNrID := c.Param("evoNrID")

	id, err := strconv.Atoi(evoNrID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}

	usuTxID := c.GetString("usuTxId")

	err = h.EvolucaoRepository.Deletar(c, id, usuTxID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Evolução não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evolução deletada com sucesso"})
}
