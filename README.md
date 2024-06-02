# Get started
Для запуска приложения необходим флаг CGO_ENABLED=1.
```
 CGO_ENABLED=1 go run cmd/main.go
```
Требуется для пакета github.com/mattn/go-sqlite3

# Техническое задание — Разработка API

## Описание задачи

Необходимо написать простейшее API для каталога товаров. Приложение должно содержать:
- Категории товаров
- Конкретные товары, которые принадлежат к какой-то категории (один товар может принадлежать нескольким категориям)
- Пользователей, которые могут авторизоваться

Возможные действия:
- Получение списка всех категорий
- Получение списка товаров в конкретной категории
- Авторизация пользователей
- Добавление/Редактирование/Удаление категории (для авторизованных пользователей)
- Добавление/Редактирование/Удаление товара (для авторизованных пользователей)

## Технические требования
1. Приложение должно быть написано на Golang
2. Хранилище данных sqlite
3. Приложение не должно быть написано с помощью какого-либо фреймворка, однако можно устанавливать для него различные пакеты
4. Результаты запросов должны быть представлены в формате JSON
5. Должна быть инструкция по запуску проекта.

## Критерии оценки
- Архитектурная организация API
- Корректная обработка внештатных ситуаций
- Покрытие кода тестами


```text
Auth service
------------------------------------------------------------------------------------------------------------------------
curl -i -X POST -d '{"email": "test2@mail.ru", "password": "qwerty123"}' http://localhost:5050/api/login
curl -i -X POST -H "Authorization: Bearer <token>" http://localhost:5050/api/logout
curl -i -X POST -d '{"email": "test2@mail.ru", "password": "qwerty123", "username":"user2"}' http://localhost:5050/api/registration
curl -i -X POST -H "Cookie: refresh_token=eyJhbG.eyJlwMD.Iuu4C4n" http://localhost:5050/api/refresh

Categories service
------------------------------------------------------------------------------------------------------------------------
curl -i -X GET "http://localhost:5050/api/categories"
curl -i -X POST -d '{"name": "Бас гитары"}' http://localhost:5050/api/categories

Products service
------------------------------------------------------------------------------------------------------------------------
curl -i -X GET "http://localhost:5050/api/products"
```