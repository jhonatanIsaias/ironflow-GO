package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"


	"ironflow/internal/database"
	"ironflow/internal/handler"
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

	exercicioRepository := repository.NovoExercicioRepository(database.DB)
	treinoRepository := repository.NovoTreinoRepository(database.DB)
	fichaTrinoRepository := repository.NovoFichaTreinoRepository(database.DB)

	treinoHandler := handler.NovoTreinoHandler(treinoRepository)
	exercicioHandler := handler.NovoExercicioHandler(exercicioRepository)
	fichaTreinoHandler := handler.NovoFichaTreinoHandler(fichaTrinoRepository)

	router := gin.Default()

	router.POST("/exercicios", exercicioHandler.CriarExercicio)
	router.GET("/exercicios/:exeNrId", exercicioHandler.BuscarPorID)
	router.PUT("/exercicios", exercicioHandler.EditarExercicio)
	router.DELETE("/exercicios/:exeNrId", exercicioHandler.DeletarPorID)

	router.POST("/treinos", treinoHandler.CriarTreino)
	router.PUT("/treinos/:treNrId", treinoHandler.EditarTreino)
	router.GET("/treinos/:treNrId", treinoHandler.BuscarPorID)
	router.GET("/treinos", treinoHandler.BuscarTodos)
	router.DELETE("/treinos/:treNrId", treinoHandler.DeletarPorID)

	router.POST("/fichas", fichaTreinoHandler.SalvarFichaTreino)
	router.PUT("/fichas/:fitNrId", fichaTreinoHandler.EditarFichaTreino)
	router.GET("/fichas/:fitNrId", fichaTreinoHandler.BuscarPorID)
	router.GET("/fichas/:treNrId", fichaTreinoHandler.BuscarTodos)
	router.DELETE("/fichas/:fitNrId", fichaTreinoHandler.DeletarPorID)



	router.Run(":8080")
}