package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Substitua pelo caminho real do seu projeto!
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

	treinoHandler := handler.NovoTreinoHandler(treinoRepository)
	exercicioHandler := handler.NovoExercicioHandler(exercicioRepository)


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

	router.Run(":8080")
}