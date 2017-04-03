package q

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
	Prelude() SQLMarshalFunc
	Predicate() SQLMarshalFunc
}

type Query struct {
	Prelude SQLMarshalFunc
	Clauses []Clause
}

func (q Query) Where(f SQLMarshalFunc) Query {
	return Query{
		q.Prelude,
		append(q.Clauses, Where(f)),
	}
}

func (q Query) OrderBy(f SQLMarshalFunc) Query {
	return Query{
		q.Prelude,
		append(q.Clauses, OrderBy(f)),
	}
}

func (q Query) SQL(d Dialect) string {
	clauses := map[ClauseKind]Clause{
		WhereClause:   nil,
		OrderByClause: nil,
		GroupByClause: nil,
		LimitClause:   nil,
		OffsetClause:  nil,
	}

	for _, clause := range q.Clauses {
		if clause.ClauseKind() == WhereClause && clauses[WhereClause] != nil {
			clauses[WhereClause] = Where(And(clauses[WhereClause].Predicate(), clause.Predicate()))
		} else {
			clauses[clause.ClauseKind()] = clause
		}
	}

	sql := q.Prelude(d)
	if clause, ok := clauses[WhereClause]; ok && clause != nil {
		sql += " " + clause.Prelude()(d) + clause.Predicate()(d)
	}
	if clause, ok := clauses[GroupByClause]; ok && clause != nil {
		sql += " " + clause.Prelude()(d) + clause.Predicate()(d)
	}
	if clause, ok := clauses[OrderByClause]; ok && clause != nil {
		sql += " " + clause.Prelude()(d) + clause.Predicate()(d)
	}
	if clause, ok := clauses[LimitClause]; ok && clause != nil {
		sql += " " + clause.Prelude()(d) + clause.Predicate()(d)
	}
	if clause, ok := clauses[OffsetClause]; ok && clause != nil {
		sql += " " + clause.Prelude()(d) + clause.Predicate()(d)
	}

	return sql
}

func Select(tables, columns SQLMarshalFunc) Query {
	return Query{
		func(d Dialect) string {
			return "SELECT " + columns(d) + " FROM " + tables(d)
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
	Star       = Raw("*")
	And        = NewConjunction("and")
	Eq         = NewConjunction("equals")
	Ascending  = NewSuffix("ascending")
	Descending = NewSuffix("descending")
	Random     = NewFunction("random")
)

type Where SQLMarshalFunc

func (w Where) ClauseKind() ClauseKind    { return WhereClause }
func (w Where) Prelude() SQLMarshalFunc   { return func(d Dialect) string { return "WHERE " } }
func (w Where) Predicate() SQLMarshalFunc { return SQLMarshalFunc(w) }

type GroupBy SQLMarshalFunc

func (g GroupBy) ClauseKind() ClauseKind    { return GroupByClause }
func (g GroupBy) Prelude() SQLMarshalFunc   { return func(d Dialect) string { return "GROUP BY " } }
func (g GroupBy) Predicate() SQLMarshalFunc { return SQLMarshalFunc(g) }

type OrderBy SQLMarshalFunc

func (o OrderBy) ClauseKind() ClauseKind    { return OrderByClause }
func (o OrderBy) Prelude() SQLMarshalFunc   { return func(d Dialect) string { return "ORDER BY " } }
func (o OrderBy) Predicate() SQLMarshalFunc { return SQLMarshalFunc(o) }

type Limit SQLMarshalFunc

func (l Limit) ClauseKind() ClauseKind    { return LimitClause }
func (l Limit) Prelude() SQLMarshalFunc   { return func(d Dialect) string { return "LIMIT " } }
func (l Limit) Predicate() SQLMarshalFunc { return SQLMarshalFunc(l) }

type Offset SQLMarshalFunc

func (o Offset) ClauseKind() ClauseKind    { return OffsetClause }
func (o Offset) Prelude() SQLMarshalFunc   { return func(d Dialect) string { return "OFFSET " } }
func (o Offset) Predicate() SQLMarshalFunc { return SQLMarshalFunc(o) }
