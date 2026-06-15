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
	RegistrarSerie(ctx context.Context, serie *model.SerieExecutada) error
	Editar(ctx context.Context, serie *model.SerieExecutada, usuTxId string) error
	BuscarPorFichaTreino(ctx context.Context, fitNrId int, usuTxId string) ([]model.SerieExecutada, error)
	BuscarPorSessao(ctx context.Context, setNrId int, usuTxId string) ([]model.SerieExecutadaDetalhada, error)
	Deletar(ctx context.Context, sexNrId int, usuTxId string) error
}

type ISessaoTreinoRepository interface {
	Salvar(ctx context.Context, sessao *model.SessaoTreino, treNrId int, usuTxId string) error
	BuscarPorFiltros(
		ctx context.Context,
		treNrId int,
		usuTxId string,
		dataInicio time.Time,
		dataFim time.Time,
		horaInicio time.Time,
		horaFim time.Time) ([]model.SessaoTreinoDetalhada, error)
	ObterSessaoHoje(ctx context.Context, treNrId int, usuTxId string) (int, bool, error)
	FinalizarSessao(ctx context.Context, setNrId int, usuTxId string) error
	VerificaSessaoAtivaHoje(ctx context.Context, usuTxId string) ([]model.SessaoTreinoDetalhada, error)
}

type SerieExecutadaHandler struct {
	serieExecutadaRepository ISerieExecutadaRepository
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

	err := h.serieExecutadaRepository.RegistrarSerie(c, &serie)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao registrar série executada"})
		return
	}

	c.JSON(http.StatusCreated, serie)
}

func (h *SerieExecutadaHandler) EditarSerieExecutada(c *gin.Context) {

	var serie model.SerieExecutada
	if err := c.ShouldBindJSON(&serie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Corpo da requisição inválido"})
		return
	}

	usuTxId, exists := c.Get("usuTxId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "Usuário não autenticado"})
		return
	}

	err := h.serieExecutadaRepository.Editar(c, &serie, usuTxId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao editar série executada"})
		return
	}

	c.JSON(http.StatusOK, serie)
}

func (h *SerieExecutadaHandler) BuscarPorSessao(c *gin.Context) {
	id := c.Param("setNrId")
	setNrId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID da sessão inválido"})
		return
	}

	usuTxId, exists := c.Get("usuTxId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "Usuário não autenticado"})
		return
	}

	series, err := h.serieExecutadaRepository.BuscarPorSessao(c, setNrId, usuTxId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar séries executadas"})
		return
	}

	c.JSON(http.StatusOK, series)
}

func (h *SerieExecutadaHandler) DeletarSerieExecutada(c *gin.Context) {
	id := c.Param("sexNrId")
	sexNrId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID da série executada inválido"})
		return
	}

	usuTxId, exists := c.Get("usuTxId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "Usuário não autenticado"})
		return
	}

	err = h.serieExecutadaRepository.Deletar(c, sexNrId, usuTxId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao deletar série executada"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Série executada deletada com sucesso"})
}
