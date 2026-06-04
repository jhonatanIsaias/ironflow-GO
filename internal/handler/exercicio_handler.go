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
	Salvar(ctx context.Context, e *model.Exercicio) error
	Editar(ctx context.Context, e *model.Exercicio) error
	BuscarTodos(ctx context.Context) ([]model.Exercicio, error)
	BuscarPorID(ctx context.Context, exeNrId int) (*model.Exercicio, error)
	Deletar(ctx context.Context, id int) error
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
	err := h.ExercicioRepository.Salvar(c, &e)
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

	err := h.ExercicioRepository.Editar(c, &e)
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
	exercicio, err := h.ExercicioRepository.BuscarPorID(c, id);

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Exercício não encontrado"})
		return
	}
	c.JSON(http.StatusOK, exercicio)
}

func (h *ExercicioHandler) DeletarPorID(c *gin.Context) {
	exeNrId := c.Param("exeNrId")

	id, err := strconv.Atoi(exeNrId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}
	err = h.ExercicioRepository.Deletar(c, id);

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Exercício não encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exercício deletado com sucesso"})
}

