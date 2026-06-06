package handler

import (
	"context"
	"ironflow/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ISerieExecutadaRepository interface {
	RegistrarSerieComSessaoAutomatica(ctx context.Context, serie *model.SerieExecutada,treNrId int) error
	Editar(ctx context.Context, serie *model.SerieExecutada) error
	BuscarPorFichaTreino(ctx context.Context, fitNrId int) ([]model.SerieExecutada, error)
	BuscarPorSessao(ctx context.Context, setNrId int) ([]model.SerieExecutadaDetalhada, error)
	Deletar(ctx context.Context, sexNrId int) error
}

type ISessaoTreinoRepository interface {
	Salvar(ctx context.Context, sessao *model.SessaoTreino) error
	BuscarPorFiltros(ctx context.Context, treNrId int, dataInicio time.Time, dataFim time.Time, horaInicio time.Time, horaFim time.Time) ([]model.SessaoTreinoDetalhada, error)
	ObterSessaoHoje(ctx context.Context, treNrId int) (int, bool, error)
}

type SerieExecutadaHandler struct {
	serieExecutadaRepository ISerieExecutadaRepository
	sessaoTreinoRepository ISessaoTreinoRepository
}

func NovoSerieExecutadaHandler(repo ISerieExecutadaRepository) *SerieExecutadaHandler {
	return &SerieExecutadaHandler{serieExecutadaRepository: repo}
}

func (h *SerieExecutadaHandler) SalvarSerieExecutada(c *gin.Context) {

	var serie model.SerieExecutada

	if err := c.ShouldBindJSON(&serie); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}

	id := c.Param("treNrId")
	treNrId, err := strconv.Atoi(id)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID do treino inválido"})
		return
	}

	err = h.serieExecutadaRepository.RegistrarSerieComSessaoAutomatica(c, &serie, treNrId)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao registrar série executada"})
		return
	}

	c.JSON(http.StatusCreated, serie)
}
