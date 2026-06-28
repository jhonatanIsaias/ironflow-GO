package handler

import (
	"context"
	"fmt"
	"ironflow/internal/model"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type IFichaTreinoRepository interface {
	Salvar(ctx context.Context, t *model.FichaTreino) error
	Editar(ctx context.Context, t *model.FichaTreino) error
	BuscarPorID(ctx context.Context, id int, usuTxId string) (*model.FichaTreinoResponse, error)
	BuscarTodos(ctx context.Context, treNrID int, exeTxNome string, usuTxId string) ([]model.FichaTreinoResponse, error)
	Deletar(ctx context.Context, id int, usuTxId string) error
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

	

	c.JSON(http.StatusCreated, gin.H{
		"status":   "success",
		"mensagem": "Ficha salva com sucesso",
		"data":     ficha,
	})

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
	
	usuTxId := c.GetString("usuTxId")

	fichaTreno,err := fit.FichaTreinoRepository.BuscarPorID(c,fitNrId,usuTxId);

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

	usuTxId := c.GetString("usuTxId")

	fichas,err := fit.FichaTreinoRepository.BuscarTodos(c,treNrId,exeTxNome,usuTxId);

	if err != nil {
		log.Fatal(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"erro":err.Error()})
		return
	}

	fichaEstruturada := fit.montarFichaTreinoEstruturada(fichas)

	fmt.Printf("ficha: %+v\n", fichaEstruturada)

	c.JSON(http.StatusOK,fichaEstruturada)

}

func (fit *FichaTreinoHandler) DeletarPorID(c *gin.Context) {

	id := c.Param("fitNrId")

	fitNrId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido. Deve ser um número."})
		return
	}
	
	usuTxId := c.GetString("usuTxId")

	err = fit.FichaTreinoRepository.Deletar(c, fitNrId, usuTxId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Ficha não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ficha deletada com sucesso"})
}

func (fit *FichaTreinoHandler) montarFichaTreinoEstruturada(fichas []model.FichaTreinoResponse) []model.FichaTreinoEstruturada {
	
	 fichaEstruturada := make([]model.FichaTreinoEstruturada, 0)

	 mapGrupos := make(map[int]int) 

	for _, f := range fichas {
		repsArray := extrairRepeticoes(f.FitTxMetaRepeticoes, f.FitNrMetaSeries,f.FitBlDropSet)

		exDetail := model.ExercicioTreino{
			FitNrID:              f.FitNrID,
			ExeNrID:              f.ExeNrID ,					
			ExeTxNome:            f.ExeTxNome,
			FitNrMetaPeso:        f.FitNrMetaPeso,
			FitTxMetaRepeticoes: repsArray,
			FitBlDropSet: f.FitBlDropSet,
		}

		if f.FitNrGrupo != nil && *f.FitNrGrupo > 0 {
			
			if idx, exists := mapGrupos[*f.FitNrGrupo]; exists {
			
				fichaEstruturada[idx].Exercicios = append(fichaEstruturada[idx].Exercicios, exDetail)
			} else {
				
				novaFicha := model.FichaTreinoEstruturada{
					IsConjugado: true,
					FitNrMetaSeries:   f.FitNrMetaSeries,
					Exercicios: []model.ExercicioTreino{exDetail},
				}
				fichaEstruturada = append(fichaEstruturada, novaFicha)
				
				mapGrupos[*f.FitNrGrupo] = len(fichaEstruturada) - 1
			}
		} else {
			
			novoBloco := model.FichaTreinoEstruturada{
				IsConjugado: false,
				FitNrMetaSeries:   f.FitNrMetaSeries,
				Exercicios:  []model.ExercicioTreino{exDetail},
			}
			fichaEstruturada = append(fichaEstruturada, novoBloco)
		}
	}

	return fichaEstruturada

}


func extrairRepeticoes(repsRaw string, totalSeries int, isDropSet bool) []string {
	
	if isDropSet {
		return []string {repsRaw}
	} 
	
	parts := strings.Split(repsRaw, "-")
	
	if len(parts) == totalSeries {
		return parts
	}

	result := make([]string, totalSeries)
	for i := 0; i < totalSeries; i++ {
		if len(parts) > 0 {
			result[i] = parts[0] 
		} else {
			result[i] = "0"
		}
	}
	return result
}