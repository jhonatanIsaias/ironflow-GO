package handler

import (
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

func (h *SessaoTreinoHandler) BuscarPorFiltros(c *gin.Context) {
	treNrIdParam := c.Query("treNrId")
	if treNrIdParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "treNrId é obrigatório"})
		return
	}

	treNrId, err := strconv.Atoi(treNrIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "treNrId inválido"})
		return
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
