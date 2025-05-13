package main

import (
	"fmt"
	"go1f/pkg/db"
	"go1f/pkg/server"

	"go1f/pkg/api"
)

func main() {
	//Инициализация БД
	fmt.Println("Инициализируем БД")

	db.Init("scheduler.db")
	api.Init()
	//Запуск сервера
	fmt.Println("Запуск сервера")

	server.Run()

}
