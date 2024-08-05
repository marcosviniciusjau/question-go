package main
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
  "errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/marcosviniciusjau/question-go/internal/store/pgstore/pgstore"
	"github.com/marcosviniciusjau/question-go/internal/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}	
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
	     "user=%s password=%s host=%s port=%s dbname=%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("POSTGRES_DB"),
		))

	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func(){
		if err:= http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed){
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}