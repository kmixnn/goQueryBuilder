package goquerybuilder

import (
	"fmt"
	"strings"
	"time"
)

type QueryBuilder struct {
	select_ []string
	from    string
	joins   []string
	where   []string
	groupBy []string
	orderBy []string
	limit   int
	offset  int
	params  []interface{}
	having  []string
	with    []string
	rawSQL  []string
}

func New() *QueryBuilder {
	return &QueryBuilder{
		select_: make([]string, 0),
		joins:   make([]string, 0),
		where:   make([]string, 0),
		groupBy: make([]string, 0),
		orderBy: make([]string, 0),
		params:  make([]interface{}, 0),
	}
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.select_ = append(qb.select_, columns...)
	return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.from = table
	return qb
}

func (qb *QueryBuilder) Join(join string, params ...interface{}) *QueryBuilder {
	qb.joins = append(qb.joins, join)
	qb.params = append(qb.params, params...)
	return qb
}

func (qb *QueryBuilder) Where(condition string, params ...interface{}) *QueryBuilder {
	qb.where = append(qb.where, condition)
	qb.params = append(qb.params, params...)
	return qb
}

func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, columns...)
	return qb
}

func (qb *QueryBuilder) OrderBy(order string) *QueryBuilder {
	qb.orderBy = append(qb.orderBy, order)
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	var parts []string

	// SELECT
	if len(qb.select_) > 0 {
		parts = append(parts, "SELECT "+strings.Join(qb.select_, ", "))
	} else {
		parts = append(parts, "SELECT *")
	}

	// FROM
	parts = append(parts, "FROM "+qb.from)

	// JOIN
	if len(qb.joins) > 0 {
		parts = append(parts, strings.Join(qb.joins, " "))
	}

	// WHERE
	if len(qb.where) > 0 {
		parts = append(parts, "WHERE "+strings.Join(qb.where, " AND "))
	}

	// GROUP BY
	if len(qb.groupBy) > 0 {
		parts = append(parts, "GROUP BY "+strings.Join(qb.groupBy, ", "))
	}

	// ORDER BY
	if len(qb.orderBy) > 0 {
		parts = append(parts, "ORDER BY "+strings.Join(qb.orderBy, ", "))
	}

	// LIMIT & OFFSET
	if qb.limit > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", qb.limit))
		if qb.offset > 0 {
			parts = append(parts, fmt.Sprintf("OFFSET %d", qb.offset))
		}
	}

	return strings.Join(parts, " "), qb.params
}

func (qb *QueryBuilder) BuildCount() (string, []interface{}) {
	var parts []string

	parts = append(parts, "SELECT COUNT(*) as count")
	parts = append(parts, "FROM "+qb.from)

	if len(qb.joins) > 0 {
		parts = append(parts, strings.Join(qb.joins, " "))
	}

	if len(qb.where) > 0 {
		parts = append(parts, "WHERE "+strings.Join(qb.where, " AND "))
	}

	return strings.Join(parts, " "), qb.params
}

// SubQuery создает подзапрос
func (qb *QueryBuilder) SubQuery(alias string) string {
	query, _ := qb.Build()
	return "(" + query + ") AS " + alias
}

// Union объединяет два запроса
func Union(queries ...*QueryBuilder) *QueryBuilder {
	var unionParts []string
	var allParams []interface{}

	for i, q := range queries {
		query, params := q.Build()
		unionParts = append(unionParts, "("+query+")")
		allParams = append(allParams, params...)

		if i < len(queries)-1 {
			unionParts = append(unionParts, "UNION ALL")
		}
	}

	result := New()
	result.from = strings.Join(unionParts, " ")
	result.params = allParams
	return result
}

// Having добавляет условие HAVING
func (qb *QueryBuilder) Having(condition string, params ...interface{}) *QueryBuilder {
	if qb.having == nil {
		qb.having = make([]string, 0)
	}
	qb.having = append(qb.having, condition)
	qb.params = append(qb.params, params...)
	return qb
}

// WithRecursive добавляет рекурсивный CTE
func WithRecursive(name string, query *QueryBuilder) *QueryBuilder {
	sql, params := query.Build()
	result := New()
	result.with = append(result.with, "WITH RECURSIVE "+name+" AS ("+sql+")")
	result.params = append(result.params, params...)
	return result
}

// Raw добавляет сырой SQL
func (qb *QueryBuilder) Raw(sql string, params ...interface{}) *QueryBuilder {
	qb.rawSQL = append(qb.rawSQL, sql)
	qb.params = append(qb.params, params...)
	return qb
}

// String возвращает SQL запрос в виде строки
func (qb *QueryBuilder) String() string {
	query, _ := qb.Build()
	return query
}

// Params возвращает все параметры запроса
func (qb *QueryBuilder) Params() []interface{} {
	_, params := qb.Build()
	return params
}

// GetQuery возвращает SQL запрос и параметры отдельно
func (qb *QueryBuilder) GetQuery() (string, []interface{}) {
	return qb.Build()
}

func (qb *QueryBuilder) Debug() *QueryBuilder {
	q, p := qb.GetQuery()
	ReplaceQueryParams(q, p)
	return qb
}

func ReplaceQueryParams(query string, args []interface{}) string {
	// Если нет параметров, возвращаем исходный запрос
	if len(args) == 0 {
		return query
	}

	// Для каждого параметра
	for _, arg := range args {
		var value string
		switch v := arg.(type) {
		case string:
			value = fmt.Sprintf("'%s'", v)
		case time.Time:
			value = fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
		case nil:
			value = "NULL"
		default:
			value = fmt.Sprintf("%v", v)
		}

		// Заменяем первый ? на значение параметра
		query = strings.Replace(query, "?", value, 1)
	}

	fmt.Printf(query)

	return query
}
