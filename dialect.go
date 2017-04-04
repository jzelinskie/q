package q

type QueryRenderer interface {
	RenderSQL(Query) string
}

type Dialect map[string]string

var MySQL = Dialect{
	"quote":       "`",
	"equals":      "=",
	"greaterthan": ">",
	"lessthan":    "<",
	"and":         "AND",
	"ascending":   "ASC",
	"descending":  "DESC",
	"random":      "RAND",
	"count":       "COUNT",
}

func (d Dialect) RenderSQL(q Query) string {
	clauses := make(map[ClauseKind]Clause)
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
		if havingClause, ok := clauses[HavingClause]; ok && clause != nil {
			sql += " " + havingClause.Prelude()(d) + havingClause.Predicate()(d)
		}
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
