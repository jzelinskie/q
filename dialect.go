package q

type Dialect map[string]string

var MySQL = Dialect{
	"quote":      "`",
	"equals":     "=",
	"ascending":  "ASC",
	"descending": "DESC",
	"random":     "RAND",
}
