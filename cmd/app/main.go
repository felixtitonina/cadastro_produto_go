package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		println(err)
		panic(err)
	}

	defer db.Close()

	repository := repository.NewProductRepositoryMysql(db)
	createProductUsecase := usecase.NewCreateProductsUseCase(repository)
	listProductsUsecase := usecase.NewListProductsUseCase(repository)

	productHandlers := web.NewProductHandlers(createProductUsecase, listProductsUsecase)

	r := chi.NewRouter()
	r.Post("/products", productHandlers.CreateProductHandler)
	r.Get("/products", productHandlers.ListProductHandler)

	http.ListenAndServe(":3000", r)

	msgChan := make(chan *kafka.Message)
	fmt.Println("msgChan")
	fmt.Println(msgChan)
	akafka.Consume([]string{"product"}, "localhost:9094", msgChan)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			println(err)
		}
		_, err = createProductUsecase.Execute(dto)
	}
}

// inferencia
