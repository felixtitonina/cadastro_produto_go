package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	akafka "github.com/felixtitonina/go-esquenta/internal/infra/kafka"
	"github.com/felixtitonina/go-esquenta/internal/infra/repository"
	"github.com/felixtitonina/go-esquenta/internal/infra/web"
	"github.com/felixtitonina/go-esquenta/internal/usecase"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:mysql@tcp(localhost:3306)/products")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	repository := repository.NewProductRepositoryMysql(db)
	createProductUsecase := usecase.NewCreateProductsUseCase(repository)
	listProductsUsecase := usecase.NewListProductsUseCase(repository)

	productHandlers := web.NewProductHandlers(createProductUsecase, listProductsUsecase)

	r := chi.NewRouter()
	r.Post("/products", productHandlers.ListProductHandler)
	r.Get("/products", productHandlers.ListProductHandler)

	http.ListenAndServe(":4444", r)

	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"products"}, "localhost:9094", msgChan)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			// logar o erro
		}
		_, err = createProductUsecase.Execute(dto)
	}
}

// inferencia
