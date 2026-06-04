# IronFlow API 🏋️‍♂️💪

O **IronFlow** é uma API REST robusta e de alta performance desenvolvida em Go para o monitoramento de treinos, evolução de cargas e acompanhamento de hipertrofia muscular. 

O projeto foi arquitetado do zero focando em controle explícito de concorrência, baixo acoplamento e máxima eficiência nas consultas ao banco de dados, utilizando o framework **Gin** e o driver nativo **pgx**.

## 🚀 Tecnologias Utilizadas

* **Go (Golang)** (v1.20+)
* **Gin Gonic** — Framework web HTTP de alta performance
* **PostgreSQL** — Banco de dados relacional
* **pgx/v5** — Driver nativo e pool de conexões para PostgreSQL
* **Godotenv** — Gerenciamento de variáveis de ambiente

## 🏗️ Arquitetura do Projeto

A estrutura segue o padrão idiomático da comunidade Go, separando rigidamente os dados (structs), o comportamento (banco de dados) e a camada de transporte (HTTP):

```text
ironflow/
├── cmd/
│   └── api/
│       └── main.go       # Ponto de entrada (Injeção de dependências e Boot)
├── internal/
│   ├── database/
│   │   └── database.go   # Gerenciamento do ciclo de vida do pgxpool
│   ├── handler/
│   │   └── ...           # Camada Web (Gin, JSON binding e Interfaces de consumo)
│   ├── model/
│   │   └── ...           # Structs de domínio (BaseEntity, Exercicio, Treino)
│   └── repository/
│       └── ...           # SQL explícito, transações e persistência pura
├── .env                  # Variáveis de ambiente (ignorado no git)
├── go.mod
└── go.sum