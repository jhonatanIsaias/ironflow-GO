package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)


var DB *pgxpool.Pool

func Conectar() error {

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("a variável de ambiente DATABASE_URL não está configurada")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("erro ao criar o pool de conexão: %v", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("erro ao fazer ping no banco de dados: %v", err)
	}

	fmt.Println("Conexão com o PostgreSQL estabelecida com sucesso!")
	return nil
}

func Fechar() {
	if DB != nil {
		DB.Close()
	}
}