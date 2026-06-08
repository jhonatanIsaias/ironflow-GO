package handler

import (
	"context"
	"fmt"
	"ironflow/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IFichaTreinoRepository interface {
	Salvar(ctx context.Context, t *model.FichaTreino) error
	Editar(ctx context.Context, t *model.FichaTreino) error
	BuscarPorID(ctx context.Context, id int) (*model.FichaTreinoResponse, error)
	BuscarTodos(ctx context.Context, treNrID int, exeTxNome string) ([]model.FichaTreinoResponse, error)
	Deletar(ctx context.Context, id int) error
	ExisteExercicioNoTreino(ctx context.Context,treNrId int, exeNrId int) (bool,error)
}

type FichaTreinoHandler struct {
	FichaTreinoRepository IFichaTreinoRepository
	TreinoRepository ITreinoRepository
}

func NovoFichaTreinoHandler(repo IFichaTreinoRepository, treinoRepository ITreinoRepository) *FichaTreinoHandler{
	return &FichaTreinoHandler{FichaTreinoRepository: repo, TreinoRepository: treinoRepository}
}

func (fit *FichaTreinoHandler) SalvarFichaTreino(c *gin.Context){
	
	var ficha model.FichaTreino;
	
	if err := c.ShouldBindJSON(&ficha); err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}

	if(ficha.FitNrID != 0){
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Não envie o ID para criar um nova ficha de treino"})
		return
	}

	usuTxId := c.GetString("usuTxId")

	_,err := fit.TreinoRepository.BuscarPorID(c,ficha.TreNrID,usuTxId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	 exist,err := fit.FichaTreinoRepository.ExisteExercicioNoTreino(c,ficha.TreNrID,ficha.ExeNrID)
		
	 if err != nil {
    	fmt.Printf("Erro ao checar exercício na ficha no banco: %v\n", err) 
    	c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha interna ao validar o exercício"})
    	return
}

	if exist {
		c.JSON(http.StatusConflict, gin.H{"erro": "Este exercício já está cadastrado neste treino"})
		return
	}

	err = fit.FichaTreinoRepository.Salvar(c,&ficha)

	if(err != nil){
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar a ficha"})
		return
	}

	c.JSON(http.StatusOK,ficha)

}

func (fit *FichaTreinoHandler) EditarFichaTreino(c *gin.Context){
	
	var ficha model.FichaTreino;
	
	if err := c.ShouldBindJSON(&ficha); err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}

	if(ficha.FitNrID == 0){
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Envie o ID para editar um nova ficha de treino"})
		return
	}

	usuTxId := c.GetString("usuTxId")

	_,err := fit.TreinoRepository.BuscarPorID(c,ficha.TreNrID,usuTxId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	
	err = fit.FichaTreinoRepository.Editar(c,&ficha)

	if(err != nil){
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao Editar a ficha"})
		return
	}
	c.JSON(http.StatusOK,ficha)

}

func (fit *FichaTreinoHandler) BuscarPorID(c *gin.Context){
	
	id := c.Param("fitNrId");

	fitNrId,err  := strconv.Atoi(id);

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}

	fichaTreno,err := fit.FichaTreinoRepository.BuscarPorID(c,fitNrId);

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK,fichaTreno)

}


func (fit *FichaTreinoHandler) BuscarTodos(c *gin.Context){

	idTreino := c.Param("treNrId");
	exeTxNome := c.Query("exeTxNome");
	
	treNrId,err := strconv.Atoi(idTreino);

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}

	fichas,err := fit.FichaTreinoRepository.BuscarTodos(c,treNrId,exeTxNome);

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar treinos"})
		return
	}

	c.JSON(http.StatusOK,fichas)

}

func (fit *FichaTreinoHandler) DeletarPorID(c *gin.Context) {

	id := c.Param("fitNrId")

	fitNrId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}
	err = fit.FichaTreinoRepository.Deletar(c, fitNrId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Ficha não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ficha deletada com sucesso"})
}