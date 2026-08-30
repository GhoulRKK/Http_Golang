package main

import (
	"clodeTask/db"
	"clodeTask/httpg"
	"context"
	"fmt"
	"net/http"
)

func main() {

	ctx := context.Background()

	pool, err := db.Simple_connection(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Close()

	if err := db.Create_DB(pool, ctx); err != nil {
		fmt.Println(err)
		return
	}

	h := &httpg.Handler{Pool: pool}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /notes", h.CreateHandler)
	mux.HandleFunc("GET /notes", h.GetAllHandler)
	mux.HandleFunc("GET /notes/{id}", h.GetByIDHandler)
	mux.HandleFunc("PUT /notes/{id}", h.PutByIDHandler)
	mux.HandleFunc("DELETE /notes/{id}", h.DeleteByIDHandler)

	fmt.Println("сервер запущен на :9091")
	if err := http.ListenAndServe(":9091", mux); err != nil {
		fmt.Println(err)
	}
}
