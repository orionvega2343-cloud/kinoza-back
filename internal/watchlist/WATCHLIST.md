# watchlist

Список фильмов/сериалов пользователя, которые он хочет посмотреть, смотрит или уже посмотрел.

## Структура

```
watchlist/
  domain/      модель WatchlistItem, статусы, ошибки, интерфейсы repo и service
  dto/         запросы и ответы для хендлеров
  handler/     gin-хендлеры + mapServiceError (ошибка -> http статус)
  mocks/       testify-моки repo и service для тестов
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
| updated_at | TIMESTAMPTZ | когда последний раз меняли статус (в API пока не отдается) |

На `(user_id, title_id)` стоит UNIQUE, повторный POST того же тайтла обновляет статус (upsert).

## Статусы

- `planned` - хочу посмотреть
- `watching` - смотрю
- `watched` - посмотрел

Константы и проверка `IsValidStatus` лежат в `domain/favorites.go`. Проверка в сервисе на Create и Update, плюс CHECK в базе.

## Эндпоинты

Роуты подключены в `cmd/api/main.go`, вся группа под `middlewares.Auth()`.

| метод | путь | что делает | ответ |
|---|---|---|---|
| POST | /watchlist/:titleId | добавить тайтл (если уже есть, обновить статус) | 201 + запись (и при обновлении тоже) |
| GET | /watchlist | список текущего юзера, новые сверху (`ORDER BY added_at DESC`) | 200 + массив (пустой список = `[]`) |
| PATCH | /watchlist/:titleId | поменять статус | 200 + `{}` |
| DELETE | /watchlist/:titleId | удалить тайтл из списка | 200 + `{}` |

Запись в ответе:

```json
{
  "id": 1,
  "title_id": 42,
  "status": "planned",
  "added_at": "2026-10-21T12:00:00Z"
}
```

Тело для POST и PATCH одинаковое, title_id берется из пути:

```json
{
  "status": "planned"
}
```

**Важно:** user_id никогда не берется из запроса (body/query), только из контекста после Auth (`c.GetString("user_id")`). Update и Delete в базе идут с `WHERE user_id = $1 AND title_id = $2`, так что чужую запись поменять или удалить нельзя, будет 404.

## Ошибки

| ошибка | статус | code |
|---|---|---|
| кривой json / нет status | 400 | FAILED_TO_BIND |
| titleId не число или <= 0 | 400 | INVALID_TITLE_ID |
| ErrInvalidStatus | 400 | INVALID_STATUS |
| ErrNotFound | 404 | NOT_FOUND |
| нет user_id в контексте | 401 | UNAUTHORIZED |
| все остальное | 500 | INTERNAL_ERROR |

На 500 наружу уходит просто "internal server error", текст ошибки базы клиенту не отдаем.

Повторное добавление того же тайтла не ошибка, а upsert: `ON CONFLICT (user_id, title_id) DO UPDATE SET status = EXCLUDED.status, updated_at = NOW()`. Статус обновляется, `added_at` остается старым.

## Тесты

```
go test ./internal/watchlist/...
```

Тесты только на сервис и хендлеры, моки на testify лежат в `mocks/`:

- `service/favorites_test.go` - сервис с `WatchlistRepoMock`, тут же проверка статусов (валидные, пустой, в другом регистре) и проброс ошибок repo
- `handler/favorites_handler_test.go` - хендлеры через httptest с `WatchlistServiceMock`: коды ответов, кривой titleId/json, 401 без user_id, что user_id берется из контекста, а не из query, и что текст ошибки базы не уходит клиенту

Тестов на repo нет, SQL проверяется только руками на реальной базе.

## TODO

- FK `title_id REFERENCES titles(id)` когда появится каталог
- `tx` в сервисе пока не используется, оставлен под будущие операции из нескольких запросов
- `updated_at` отдавать в ответе
- `domain.ErrDB` пока используется только в тестах, repo возвращает обернутые ошибки базы
