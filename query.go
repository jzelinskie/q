package q

type SQLMarshalFunc func(Dialect) string

type ClauseKind string

var (
	WhereClause   = ClauseKind("WHERE")
	GroupByClause = ClauseKind("GROUP BY")
	HavingClause  = ClauseKind("HAVING")
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
	clauses := make([]Clause, len(q.Clauses))
	copy(clauses, q.Clauses)
	return Query{
		q.Prelude,
		append(clauses, Where(f)),
	}
}
func (q Query) GroupBy(f SQLMarshalFunc) Query {
	clauses := make([]Clause, len(q.Clauses))
	copy(clauses, q.Clauses)
	return Query{
		q.Prelude,
		append(clauses, GroupBy(f)),
	}
}

// Having will only be used if there is a GroupBy present.
func (q Query) Having(f SQLMarshalFunc) Query {
	clauses := make([]Clause, len(q.Clauses))
	copy(clauses, q.Clauses)
	return Query{
		q.Prelude,
		append(clauses, Having(f)),
	}
}

func (q Query) OrderBy(f SQLMarshalFunc) Query {
	clauses := make([]Clause, len(q.Clauses))
	copy(clauses, q.Clauses)
	return Query{
		q.Prelude,
		append(clauses, OrderBy(f)),
	}
}

func (q Query) Limit(f SQLMarshalFunc) Query {
	clauses := make([]Clause, len(q.Clauses))
	copy(clauses, q.Clauses)
	return Query{
		q.Prelude,
		append(clauses, Limit(f)),
	}
}

func (q Query) Offset(f SQLMarshalFunc) Query {
	clauses := make([]Clause, len(q.Clauses))
	copy(clauses, q.Clauses)
	return Query{
		q.Prelude,
		append(clauses, Offset(f)),
	}
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
	Star        = Raw("*")
	And         = NewConjunction("and")
	Eq          = NewConjunction("equals")
	GreaterThan = NewConjunction("greaterthan")
	LessThan    = NewConjunction("lessthan")
	Ascending   = NewSuffix("ascending")
	Descending  = NewSuffix("descending")
	Random      = NewFunction("random")
	Count       = NewFunction("count")
)

type Where SQLMarshalFunc

func (w Where) ClauseKind() ClauseKind    { return WhereClause }
func (w Where) Prelude() SQLMarshalFunc   { return func(d Dialect) string { return "WHERE " } }
func (w Where) Predicate() SQLMarshalFunc { return SQLMarshalFunc(w) }

type GroupBy SQLMarshalFunc

func (g GroupBy) ClauseKind() ClauseKind    { return GroupByClause }
func (g GroupBy) Prelude() SQLMarshalFunc   { return func(d Dialect) string { return "GROUP BY " } }
func (g GroupBy) Predicate() SQLMarshalFunc { return SQLMarshalFunc(g) }

type Having SQLMarshalFunc

func (h Having) ClauseKind() ClauseKind    { return HavingClause }
func (h Having) Prelude() SQLMarshalFunc   { return func(d Dialect) string { return "HAVING " } }
func (h Having) Predicate() SQLMarshalFunc { return SQLMarshalFunc(h) }

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
