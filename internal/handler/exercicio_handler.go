package handler

import (
	"context"
	"fmt"
	"ironflow/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IExercicioRepository interface {
	Salvar(ctx context.Context, e *model.Exercicio,usuTxID string) error
	Editar(ctx context.Context, e *model.Exercicio,usuTxID string) error
	BuscarTodos(ctx context.Context,usuTxID string) ([]model.Exercicio, error)
	BuscarPorID(ctx context.Context, exeNrId int,usuTxID string) (*model.Exercicio, error)
	Deletar(ctx context.Context, id int,usuTxID string) error
}

type ExercicioHandler struct {
	ExercicioRepository IExercicioRepository
}

func NovoExercicioHandler(repo IExercicioRepository) *ExercicioHandler {
	return &ExercicioHandler{ExercicioRepository: repo}
}

func (h *ExercicioHandler) CriarExercicio(c *gin.Context) {

	var e model.Exercicio

	if err := c.ShouldBindJSON(&e); err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}
	if e.ExeNrID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Não envie o ID para criar um novo exercício"})
		return
	}
	usuTxId := c.GetString("usuTxId")
	
	err := h.ExercicioRepository.Salvar(c, &e,usuTxId)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar o exercício"})
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *ExercicioHandler) EditarExercicio(c *gin.Context) {

	var e model.Exercicio
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}

	if e.ExeNrID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID do exercício é obrigatório para edição"})
		return
	}

	usuTxId := c.GetString("usuTxId")

	err := h.ExercicioRepository.Editar(c, &e,usuTxId)
	if err != nil {
		fmt.Printf("Erro ao editar exercício: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao editar o exercício"})
		return
	}
	c.JSON(http.StatusOK, e)
}


func (h *ExercicioHandler) BuscarPorID(c *gin.Context) {
	exeNrId := c.Param("exeNrId")

	id, err := strconv.Atoi(exeNrId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}

	usuTxId := c.GetString("usuTxId")
	exercicio, err := h.ExercicioRepository.BuscarPorID(c, id, usuTxId);

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Exercício não encontrado"})
		return
	}
	c.JSON(http.StatusOK, exercicio)
}

func (h *ExercicioHandler) BuscarTodos(c *gin.Context) {

	usuTxId := c.GetString("usuTxId")
	
	exercicios, err := h.ExercicioRepository.BuscarTodos(c,usuTxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar exercícios"})
		return
	}
	c.JSON(http.StatusOK, exercicios)
}

func (h *ExercicioHandler) DeletarPorID(c *gin.Context) {
	exeNrId := c.Param("exeNrId")

	id, err := strconv.Atoi(exeNrId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}

	usuTxId := c.GetString("usuTxId")
	err = h.ExercicioRepository.Deletar(c, id, usuTxId);

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Exercício não encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exercício deletado com sucesso"})
}

