package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"ironflow/internal/database"
	"ironflow/internal/handler"
	"ironflow/internal/middleware"
	"ironflow/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Aviso: Arquivo .env não encontrado. Certifique-se de que ele está na raiz do projeto.")
	}

	err = database.Conectar()
	if err != nil {
		log.Fatalf("Falha crítica ao iniciar o banco: %v", err)
	}
	defer database.Fechar()

	database.RodarMigrations()

	usuarioRepository := repository.NovoUsuarioRepository(database.DB)
	exercicioRepository := repository.NovoExercicioRepository(database.DB)
	treinoRepository := repository.NovoTreinoRepository(database.DB)
	fichaTrinoRepository := repository.NovoFichaTreinoRepository(database.DB)
	sessaoTreinoRepository := repository.NovoSessaoTreinoRepository(database.DB)
	serieExecutadaRepository := repository.NovoSerieExecutadaRepository(database.DB)
	evolucaoRepository := repository.NovoEvolucaoRepository(database.DB)

	usuarioHandler := handler.NovoUsuarioHandler(usuarioRepository)
	treinoHandler := handler.NovoTreinoHandler(treinoRepository)
	exercicioHandler := handler.NovoExercicioHandler(exercicioRepository)
	fichaTreinoHandler := handler.NovoFichaTreinoHandler(fichaTrinoRepository, treinoRepository)
	sessaoTreinoHandler := handler.NovoSessaoTreinoHandler(sessaoTreinoRepository)
	serieExecutadaHandler := handler.NovoSerieExecutadaHandler(serieExecutadaRepository)
	evolucaoHandler := handler.NovoEvolucaoHandler(evolucaoRepository)

	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.POST("/login", usuarioHandler.Login)
	v1.POST("/usuarios", usuarioHandler.SalvarUsuario)
	v1.POST("/refresh", usuarioHandler.Refresh)

	protected := v1.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		// EXERCÍCIOS
		protected.POST("/exercicios", exercicioHandler.CriarExercicio)
		protected.GET("/exercicios/:exeNrId", exercicioHandler.BuscarPorID)
		protected.GET("/exercicios", exercicioHandler.BuscarTodos)
		protected.PUT("/exercicios", exercicioHandler.EditarExercicio)
		protected.DELETE("/exercicios/:exeNrId", exercicioHandler.DeletarPorID)

		// TREINOS
		protected.POST("/treinos", treinoHandler.CriarTreino)
		protected.PUT("/treinos", treinoHandler.EditarTreino)
		protected.GET("/treinos/:treNrId", treinoHandler.BuscarPorID)
		protected.GET("/treinos", treinoHandler.BuscarTodos)
		protected.DELETE("/treinos/:treNrId", treinoHandler.DeletarPorID)

		// FICHAS
		protected.POST("/fichas", fichaTreinoHandler.SalvarFichaTreino)
		protected.PUT("/fichas", fichaTreinoHandler.EditarFichaTreino)
		protected.GET("/fichas/:fitNrId", fichaTreinoHandler.BuscarPorID)
		protected.GET("/fichas/treino/:treNrId", fichaTreinoHandler.BuscarTodos)
		protected.DELETE("/fichas/:fitNrId", fichaTreinoHandler.DeletarPorID)

		// SESSÕES
		protected.GET("/sessoes", sessaoTreinoHandler.BuscarPorFiltros)
		protected.POST("/sessoes/:treNrId",sessaoTreinoHandler.CriarSessaoTreino)
		protected.PUT("/sessoes/:setNrId",sessaoTreinoHandler.FinalizarSessao)
		protected.GET("sessoes/ativas/:treNrId",sessaoTreinoHandler.VerificaSessaoAtivaHoje)
	

		// SÉRIES
		protected.POST("/series/:treNrId", serieExecutadaHandler.SalvarSerieExecutada)
		protected.PUT("/series", serieExecutadaHandler.EditarSerieExecutada)
		protected.GET("/series/sessao/:setNrId", serieExecutadaHandler.BuscarPorSessao)
		protected.DELETE("/series/:sexNrId", serieExecutadaHandler.DeletarSerieExecutada)

		// EVOLUÇÕES
		protected.POST("/evolucoes", evolucaoHandler.CriarEvolucao)
		protected.PUT("/evolucoes", evolucaoHandler.EditarEvolucao)
		protected.GET("/evolucoes/:evoNrID", evolucaoHandler.BuscarPorID)
		protected.GET("/evolucoes/recente", evolucaoHandler.BuscarMaisRecente)
		protected.GET("/evolucoes", evolucaoHandler.BuscarTodos)
		protected.DELETE("/evolucoes/:evoNrID", evolucaoHandler.DeletarPorID)

	}

	router.Run(":8080")
}
