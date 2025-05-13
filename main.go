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

	errDb := db.Init("scheduler.db")
	if errDb != nil {
		fmt.Println("ошибка при инициализации БД")
		return
	}
	api.Init()

	//Запуск сервера
	fmt.Println("Запуск сервера")

	errSrv := server.Run()
	if errSrv != nil {
		fmt.Println("Ошибка при запуске сервера")
		return
	}
}
