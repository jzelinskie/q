package q

import "testing"

func TestSimpleQuery(t *testing.T) {
	query := Select([]string{"Users"}, []string{"*"}).
		Where(Eq(Raw("name"), Raw("Jimmy"))).
		OrderBy(Random(nil)).
		SQL(MySQL)
	expected := "SELECT * FROM Users WHERE name = Jimmy ORDER BY RAND()"
	if query != expected {
		t.Errorf("failed to render correct SQL from query => wanted: %q, got: %q", expected, query)
	}
}
