package main

import (
	db "avito_test/pkg/postgres_db"
	"avito_test/pkg/postgres_db/user"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func load_env() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}
func createUsers() error {
	var u user.User
	ctx := context.Background()
	conn, tx, err := db.ConnectPoolTrx(ctx)
	defer conn.Release()
	defer func() {
		if err != nil {
			tx.Rollback(ctx)

		} else {
			tx.Commit(ctx)
		}
	}()
	if err != nil {
		return fmt.Errorf("createUsers error - %v", err.Error())
	}
	for i := 1; i <= 3; i++ {
		u = user.User{UserName: fmt.Sprintf("admin_%v", i),
			Password: fmt.Sprintf("admin_%v", i),
			IsAdmin:  true}
		if err = u.Insert(ctx, conn); err != nil {
			return fmt.Errorf("createUsers error - %v", err.Error())
		}
	}
	for i := 1; i <= 5; i++ {
		u = user.User{UserName: fmt.Sprintf("user_%v", i),
			Password: fmt.Sprintf("user_%v", i),
			IsAdmin:  false}
		if err = u.Insert(ctx, conn); err != nil {
			return fmt.Errorf("createUsers error - %v", err.Error())
		}
	}
	return nil
}

func main() {
	load_env()
	err := createUsers()
	if err != nil {
		fmt.Println("err ", err.Error())
	}
}
