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
	serieExecutadaRepository := repository.NovoSerieExecutadaRepository(database.DB)

	
	usuarioHandler := handler.NovoUsuarioHandler(usuarioRepository)
	treinoHandler := handler.NovoTreinoHandler(treinoRepository)
	exercicioHandler := handler.NovoExercicioHandler(exercicioRepository)
	fichaTreinoHandler := handler.NovoFichaTreinoHandler(fichaTrinoRepository)
	serieExecutadaHandler := handler.NovoSerieExecutadaHandler(serieExecutadaRepository)

	router := gin.Default()

    v1 := router.Group("/api/v1")

    v1.POST("/login", usuarioHandler.Login)
    v1.POST("/usuarios", usuarioHandler.SalvarUsuario)

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
        protected.GET("/fichas", fichaTreinoHandler.BuscarTodos)
        protected.DELETE("/fichas/:fitNrId", fichaTreinoHandler.DeletarPorID)

        // SÉRIES
        protected.POST("/series", serieExecutadaHandler.SalvarSerieExecutada)
        protected.PUT("/series", serieExecutadaHandler.EditarSerieExecutada)
        protected.GET("/series/sessao/:setNrId", serieExecutadaHandler.BuscarPorSessao)
        protected.DELETE("/series/:sexNrId", serieExecutadaHandler.DeletarSerieExecutada)
    }

	router.Run(":8080")
}