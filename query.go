package q

import "strings"

type SQLMarshalFunc func(Dialect) string

type ClauseKind string

var (
	WhereClause   = ClauseKind("WHERE")
	GroupByClause = ClauseKind("GROUP BY")
	OrderByClause = ClauseKind("ORDER BY")
	LimitClause   = ClauseKind("LIMIT")
	OffsetClause  = ClauseKind("OFFSET")
)

type Clause interface {
	ClauseKind() ClauseKind
	SQLMarshalFunc() SQLMarshalFunc
}

type Query struct {
	prelude SQLMarshalFunc
	clauses []Clause
}

func (q Query) Where(f SQLMarshalFunc) Query {
	return Query{
		q.prelude,
		append(q.clauses, Where(f)),
	}
}

func (q Query) OrderBy(f SQLMarshalFunc) Query {
	return Query{
		q.prelude,
		append(q.clauses, OrderBy(f)),
	}
}

func (q Query) SQL(d Dialect) string {
	renderedClauses := map[ClauseKind]string{
		WhereClause:   "",
		OrderByClause: "",
		GroupByClause: "",
		LimitClause:   "",
		OffsetClause:  "",
	}

	for _, clause := range q.clauses {
		if clause.ClauseKind() == WhereClause {
			renderedClauses[WhereClause] = clause.SQLMarshalFunc()(d)
		} else {
			renderedClauses[clause.ClauseKind()] = clause.SQLMarshalFunc()(d)
		}
	}

	sql := q.prelude(d)
	if renderedClauses[WhereClause] != "" {
		sql += " " + renderedClauses[WhereClause]
	}
	if renderedClauses[GroupByClause] != "" {
		sql += " " + renderedClauses[GroupByClause]
	}
	if renderedClauses[OrderByClause] != "" {
		sql += " " + renderedClauses[OrderByClause]
	}
	if renderedClauses[LimitClause] != "" {
		sql += " " + renderedClauses[LimitClause]
	}
	if renderedClauses[OffsetClause] != "" {
		sql += " " + renderedClauses[OffsetClause]
	}

	return sql
}

func Select(tables, columns []string) Query {
	return Query{
		func(d Dialect) string {
			return "SELECT " + strings.Join(columns, ", ") + " FROM " + strings.Join(tables, ", ")
		},
		nil,
	}
}

type Binomial func(first, second SQLMarshalFunc) SQLMarshalFunc

func NewConjunction(tokenName string) Binomial {
	return func(first, second SQLMarshalFunc) SQLMarshalFunc {
		return func(d Dialect) string { return first(d) + " " + d[tokenName] + " " + second(d) }
	}
}

type Monomial func(SQLMarshalFunc) SQLMarshalFunc

func NewFunction(tokenName string) Monomial {
	return func(f SQLMarshalFunc) SQLMarshalFunc {
		return func(d Dialect) string {
			args := ""
			if f != nil {
				args = f(d)
			}

			return d[tokenName] + "(" + args + ")"
		}
	}
}

func NewSuffix(tokenName string) Monomial {
	return func(f SQLMarshalFunc) SQLMarshalFunc {
		return func(d Dialect) string { return f(d) + " " + d[tokenName] }
	}
}

type Constant func() SQLMarshalFunc

func NewConstant(tokenName string) Constant {
	return func() SQLMarshalFunc {
		return func(d Dialect) string { return d[tokenName] }
	}
}

func Raw(s string) SQLMarshalFunc { return func(d Dialect) string { return s } }

var (
	And        = NewConjunction("and")
	Eq         = NewConjunction("equals")
	Ascending  = NewSuffix("ascending")
	Descending = NewSuffix("descending")
	Random     = NewFunction("random")
)

type orderBy struct{ predicate SQLMarshalFunc }

func OrderBy(predicate SQLMarshalFunc) orderBy { return orderBy{predicate} }
func (o orderBy) ClauseKind() ClauseKind       { return OrderByClause }
func (o orderBy) SQLMarshalFunc() SQLMarshalFunc {
	return func(d Dialect) string { return "ORDER BY " + o.predicate(d) }
}

type where struct{ predicate SQLMarshalFunc }

func Where(predicate SQLMarshalFunc) where { return where{predicate} }
func (w where) ClauseKind() ClauseKind     { return WhereClause }
func (w where) SQLMarshalFunc() SQLMarshalFunc {
	return func(d Dialect) string { return "WHERE " + w.predicate(d) }
}
