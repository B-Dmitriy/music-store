package main

import "github.com/B-Dmitriy/music-store/internal/app"

func main() {
	musicShop, err := app.New()
	if err != nil {
		panic("application initialization failed")
	}

	// TODO: graceful shutdown
	musicShop.Run()
}
