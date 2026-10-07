# watchlist

Список фильмов/сериалов пользователя, которые он хочет посмотреть, смотрит или уже посмотрел.

## Структура

```
watchlist/
  domain/      модель WatchlistItem, статусы, ошибки, интерфейсы repo и service
  dto/         запросы и ответы для хендлеров
  handler/     gin-хендлеры + mapServiceError (ошибка -> http статус)
  repository/  работа с postgres
  service/     проверка статуса, дальше вызов repo
```

## Таблица

Миграция `migrations/000003_create_watchlist_item.up.sql`

| поле | тип | что это |
|---|---|---|
| id | SERIAL | id записи |
| user_id | UUID | чей список, FK на users, при удалении юзера записи тоже удаляются |
| title_id | INTEGER | id фильма/сериала (FK пока нет, каталога еще нет) |
| status | TEXT | статус, см. ниже |
| added_at | TIMESTAMPTZ | когда добавили |

На `(user_id, title_id)` стоит UNIQUE, один и тот же тайтл два раза добавить нельзя.

## Статусы

- `planned` - хочу посмотреть
- `watching` - смотрю
- `watched` - посмотрел

Константы и проверка `IsValidStatus` лежат в `domain/favorites.go`. Проверка в сервисе на Create и Update, в базе ограничения пока нет.

## Эндпоинты

Роуты еще не подключены в main.go, пути ниже примерные. Все требуют авторизацию (`middlewares.Auth()`).

| метод | путь | что делает | ответ |
|---|---|---|---|
| POST | /watchlist | добавить тайтл | 201 + созданная запись |
| GET | /watchlist | список текущего юзера | 200 + массив |
| PATCH | /watchlist/:id | поменять статус | 200 |
| DELETE | /watchlist/:id | удалить запись | 200 |

Пример тела для POST:

```json
{
  "title_id": 42,
  "status": "planned"
}
```

Для PATCH только `{"status": "watched"}`.

**Важно:** user_id никогда не берется из запроса (body/query), только из контекста после Auth (`c.GetString("user_id")`). Update и Delete в базе идут с `WHERE id = $1 AND user_id = $2`, так что чужую запись поменять или удалить нельзя, будет 404.

## Ошибки

| ошибка | статус | code |
|---|---|---|
| кривой json | 400 | FAILED_TO_BIND |
| id не число | 400 | INVALID_ID |
| ErrInvalidStatus | 400 | INVALID_STATUS |
| ErrNotFound | 404 | NOT_FOUND |
| ErrAlreadyExists | 409 | ALREADY_EXISTS |
| все остальное | 500 | INTERNAL_ERROR |

На 500 наружу уходит просто "internal server error", текст ошибки базы клиенту не отдаем.

Дубль ловится через `ON CONFLICT (user_id, title_id) DO NOTHING RETURNING ...`: если ничего не вставилось, приходит `sql.ErrNoRows` и repo превращает его в `ErrAlreadyExists`.

## Тесты

```
go test ./internal/watchlist/...
```

Есть тесты на `IsValidStatus` (domain) и на сервис (с фейковым repo, без базы).

## TODO

- подключить роуты в main.go
- тесты на хендлеры (httptest)
- тесты на repo (нужен postgres)
- CHECK на status в миграции
- FK `title_id REFERENCES titles(id)` когда появится каталог
- `tx` в сервисе пока не используется, оставлен под будущие операции из нескольких запросов
