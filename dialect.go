package q

type Dialect map[string]string

var MySQL = Dialect{
	"quote":      "`",
	"equals":     "=",
	"and":        "AND",
	"ascending":  "ASC",
	"descending": "DESC",
	"random":     "RAND",
}
