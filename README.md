# Veksel (тестовое задание)
#### Перед запуском требуется файл .env содержащий:
```cfg
DB_USER=someuser
POSTGRES_PASSWORD=somepassword
```
Я использовал стандартного пользователя postgres из PostgreSQL

#### Запуск проекта:   
Проект запускается через: `make up`  
Если проект запускается впервые, то после запуска postgres контейнера нужно выполнить:
`make migrate`


## API Reference

#### Получить список заметок

```http
  GET /api/notes
```
**Header:** Authorization: Basic `base64 string` (required)

---

#### Создать заметку

```http
  POST /api/notes
```
**Header:** Authorization: Basic `base64 string` (required)  
**Body:** `json`
```json
{
    "header": "YorBestHeader",
    "content": "Your gorgeous note"
}
```

---

[@сountenum404](https://www.github.com/сountenum404) Шабашов Денис
