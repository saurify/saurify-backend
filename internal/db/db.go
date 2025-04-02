package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(dsn string){
	var err error 
	DB, err = pgxpool.New(context.Background(), dsn)
	if err!=nil{
		log.Fatalf("failed to connect to DB")
	}
	fmt.Println("Db connected!")
}