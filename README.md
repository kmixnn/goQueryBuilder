# Go Query Builder

Простой и гибкий конструктор SQL-запросов для Go.

## Возможности

- Построение SELECT-запросов с поддержкой всех основных SQL-операторов
- Поддержка JOIN-ов
- WHERE условия с параметрами
- GROUP BY и HAVING
- ORDER BY
- LIMIT и OFFSET
- Подзапросы
- UNION операции
- Рекурсивные CTE (WITH RECURSIVE)
- Сырые SQL-запросы
- Отладка запросов

## Примеры использования

1. Простой SELECT:
```go
qb := New()
```

2. JOIN:
```go
qb := New()
qb.Select("u.id", "u.name", "o.order_date")
  .From("users u")
  .Join("LEFT JOIN orders o ON o.user_id = u.id")
  .Where("u.active = ?", true)
```

3. Подзапрос:
```
subQuery := New().Select("id").From("users").Where("age > ?", 21)
mainQuery := New()
  .Select("*")
  .From(subQuery.SubQuery("eligible_users"))
```

4. WITH RECURSIVE:
```
recursive := New()
  .Select("id", "parent_id", "name")
  .From("categories")
withQuery := WithRecursive("category_tree", recursive)
```

## Методы

- Select(columns ...string) - выбор колонок
- From(table string) - указание таблицы
- Join(join string, params ...interface{}) - добавление JOIN
- Where(condition string, params ...interface{}) - добавление условий WHERE
- GroupBy(columns ...string) - группировка
- OrderBy(order string) - сортировка
- Limit(limit int) - ограничение количества записей
- Offset(offset int) - пропуск записей
- Having(condition string, params ...interface{}) - условия HAVING
- Raw(sql string, params ...interface{}) - сырой SQL
- Build() - построение запроса
- BuildCount() - построение COUNT запроса
- Debug() - отладка запроса
