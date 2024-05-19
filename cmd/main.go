package main

import "github.com/B-Dmitriy/test-api/internal/app"

/*
Категории товаров
Товары
Пользоавтели
Аутентификация, авторизация
*/

func main() {
	musicShop, err := app.New()
	if err != nil {
		panic("application initialization failed")
	}

	// TODO: graceful shutdown
	musicShop.Run()
}
