package handler

import (
	"ironflow/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SessaoTreinoHandler struct {
	SessaoTreinoRepository ISessaoTreinoRepository
}

func NovoSessaoTreinoHandler(repo ISessaoTreinoRepository) *SessaoTreinoHandler {
	return &SessaoTreinoHandler{SessaoTreinoRepository: repo}
}

func (h *SessaoTreinoHandler) CriarSessaoTreino(c *gin.Context) {

	var sessao model.SessaoTreino

	treNrIdParam := c.Param("treNrId")

	treNrId, err := strconv.Atoi(treNrIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "treNrId inválido"})
		return
	}

	usuTxId := c.GetString("usuTxId")

	err = h.SessaoTreinoRepository.Salvar(c, &sessao, treNrId, usuTxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar a sessão de treino: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sessao)
}

func (h *SessaoTreinoHandler) BuscarPorFiltros(c *gin.Context) {

	treNrIdQuery := c.Query("treNrId")
	var treNrId int
	var err error

	if treNrIdQuery != "" {
		treNrId, err = strconv.Atoi(treNrIdQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "treNrId inválido"})
			return
		}
	}

	usuTxId := c.GetString("usuTxId")

	dataInicio, err := parseOptionalDate(c.Query("dataInicio"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dataInicio inválida. Use YYYY-MM-DD."})
		return
	}

	dataFim, err := parseOptionalDate(c.Query("dataFim"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dataFim inválida. Use YYYY-MM-DD."})
		return
	}

	horaInicio, err := parseOptionalTime(c.Query("horaInicio"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "horaInicio inválida. Use HH:MM ou HH:MM:SS."})
		return
	}

	horaFim, err := parseOptionalTime(c.Query("horaFim"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "horaFim inválida. Use HH:MM ou HH:MM:SS."})
		return
	}

	sessoes, err := h.SessaoTreinoRepository.BuscarPorFiltros(
		c,
		treNrId,
		usuTxId,
		dataInicio,
		dataFim,
		horaInicio,
		horaFim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar sessões de treino"})
		return
	}

	c.JSON(http.StatusOK, sessoes)
}

func (h *SessaoTreinoHandler) FinalizarSessao(c *gin.Context) {

	setNrIdParam := c.Param("setNrId")

	setNrId, err := strconv.Atoi(setNrIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "setNrId inválido"})
		return
	}

	usuTxId := c.GetString("usuTxId")

	err = h.SessaoTreinoRepository.FinalizarSessao(c, setNrId, usuTxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao finalizar a sessão de treino: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Sessão finalizada com sucesso"})
}

func (h *SessaoTreinoHandler) VerificaSessaoAtivaHoje(c *gin.Context) {

	usuTxId := c.GetString("usuTxId")

	sessoes, err := h.SessaoTreinoRepository.VerificaSessaoAtivaHoje(c, usuTxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar sessões ativas de hoje: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessoes)
}

func parseOptionalDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", value)
}

func parseOptionalTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse("15:04", value)
	if err == nil {
		return parsed, nil
	}
	return time.Parse("15:04:05", value)
}
