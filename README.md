# Veksel (тестовое задание)

Особенности реализации:   
- Заметка с ошибкой сохраняется в сервисе, но в ответ пользователю посылается сохраненная заметка и список найденых ошибок. Пользователь сам должен решить какой вариант ему подходит для редактирования.
- Валидируется только Content заметки, Header всегда принимается как есть.
- Пользователи захардкожены в таблице users в PostgreSQL.

#### Перед запуском требуется файл .env содержащий:
```cfg
DB_USER=someuser
POSTGRES_PASSWORD=somepassword
REDIS_HOST=redis:6379
```
Я использовал стандартного пользователя postgres из PostgreSQL

#### Запуск проекта:   
Проект запускается через: `make up` или командой `docker-compose up --build veksel`


## API Reference

#### Получить список заметок

```http
  GET /api/notes
```
**Header:** Authorization: Basic `base64 string` (required)   
**Example:**
```bash
curl --location 'http://localhost:4567/api/notes' \
--header 'Authorization: Basic cmdvc2xpbmc6ZHJpdmU='
```
**Response:**
```json
[
    {
        "id": 1,
        "header": "HEADER",
        "content": "Спеллер, найди пожалста ашибке"
    },
    {
        "id": 2,
        "header": "HEADER1",
        "content": "съешь ещё эти мягках французскил Дулок, да выпей чаю"
    }
]
```
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

**Example:**
```bash
curl --location 'http://localhost:4567/api/notes' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic cmdvc2xpbmc6ZHJpdmU=' \
--data '{
    "header": "HEADER",
    "content": "Спеллер, найди пожаста ашибки"
}'
```

**Response:**
```json
{
    "created_note": {
        "header": "HEADER",
        "content": "Спеллер, найди пожаста ашибки"
    },
    "spells": [
        {
            "code": 1,
            "pos": 15,
            "row": 0,
            "col": 15,
            "len": 7,
            "word": "пожаста",
            "s": [
                "пожалуйста",
                "пожалста"
            ]
        },
        {
            "code": 1,
            "pos": 23,
            "row": 0,
            "col": 23,
            "len": 6,
            "word": "ашибки",
            "s": [
                "ошибки"
            ]
        }
    ]
}
```

---

[@сountenum404](https://www.github.com/сountenum404) Шабашов Денис
