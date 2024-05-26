package main

import "github.com/B-Dmitriy/music-store/internal/app"

/*
Категории товаров
Товары
Пользоавтели
Аутентификация, авторизация

MW example
func Logging(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    start := time.Now()
    next.ServeHTTP(w, req)
    log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
  })
}
*/

func main() {
	musicShop, err := app.New()
	if err != nil {
		panic("application initialization failed")
	}

	// TODO: graceful shutdown
	musicShop.Run()
}
