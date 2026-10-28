# history

Прогресс просмотра: на какой секунде пользователь остановился в каждом тайтле. Из этого собирается блок «продолжить смотреть».

## Структура

```
history/
  domain/      модель WatchProgress, Validate, ошибки, интерфейсы repo и service
  dto/         запрос и ответ для хендлеров
  handler/     gin-хендлеры + mapServiceError (ошибка -> http статус)
  mocks/       testify-моки repo и service для тестов
  repository/  работа с postgres
  service/     Validate, дальше вызов repo
```

## Таблица

Миграция `migrations/000004_create_watch_progress.up.sql`

| поле | тип | что это |
|---|---|---|
| user_id | UUID | чей прогресс, FK на users, при удалении юзера записи тоже удаляются |
| title_id | INTEGER | id фильма/сериала, FK нет: домены связаны только по id (ТЗ) |
| position_sec | INTEGER | на какой секунде остановился, CHECK `>= 0` |
| updated_at | TIMESTAMPTZ | когда последний раз сохраняли прогресс |

PK составной: `(user_id, title_id)`, отдельного `id` нет (так в ТЗ). Одна запись на пару юзер + тайтл, повторное сохранение обновляет позицию (upsert).

Индекс `(user_id, updated_at DESC)` под запрос «продолжить смотреть».

## Эндпоинты

Роуты подключены в `cmd/api/main.go`, вся группа под `middlewares.Auth()`.

| метод | путь | что делает | ответ |
|---|---|---|---|
| PUT | /history/:titleId/progress | сохранить позицию (нет записи - создать, есть - обновить) | 200 + запись |
| GET | /history/continue-watching | прогресс текущего юзера, последние просмотренные сверху (`ORDER BY updated_at DESC`) | 200 + массив (пустой список = `[]`) |

Запись в ответе:

```json
{
  "title_id": 42,
  "position_sec": 1200,
  "updated_at": "2026-10-29T12:00:00Z"
}
```

Тело для PUT, title_id берется из пути:

```json
{
  "position_sec": 1200
}
```

`position_sec` в DTO это `*int` с `binding:"required,gte=0"`: обычный `required` на `int` не пропустил бы `0`, а позиция 0 валидна (только начал смотреть).

**Важно:** user_id никогда не берется из запроса (body/query), только из контекста после Auth (`c.GetString("user_id")`).

## Ошибки

| ошибка | статус | code |
|---|---|---|
| кривой json / нет position_sec / position_sec < 0 | 400 | FAILED_TO_BIND |
| titleId не число или <= 0 | 400 | INVALID_TITLE_ID |
| ErrInvalidPosition | 400 | INVALID_POSITION |
| нет user_id в контексте | 401 | UNAUTHORIZED |
| все остальное | 500 | INTERNAL_ERROR |

На 500 наружу уходит просто "internal server error", текст ошибки базы клиенту не отдаем.

Сервис оборачивает ошибку Validate через `%w`, поэтому `errors.Is` в `mapServiceError` находит `ErrInvalidPosition`.

## Тесты

```
go test ./internal/history/...
```

Тесты только на сервис и хендлеры, моки на testify лежат в `mocks/`:

- `service/watch_progress_service_test.go` - сервис с `WatchProgressRepoMock`: валидная позиция, ноль, отрицательная (repo не вызывается), проброс ошибок repo
- `handler/watch_progress_handler_test.go` - хендлеры через httptest с `WatchProgressServiceMock`: коды ответов, кривой titleId/json, позиция 0 / нет позиции / отрицательная, 401 без user_id, что user_id берется из контекста, а не из query, и что текст ошибки базы не уходит клиенту

Тестов на repo нет, SQL проверяется только руками на реальной базе.

## TODO

- событие `ProgressUpdated` в `watch.progress.events` (Kafka + outbox), по ТЗ это поздний этап
- `domain.ErrDatabase` пока используется только в тестах, repo возвращает обернутые ошибки базы
