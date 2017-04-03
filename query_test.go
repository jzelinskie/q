package q

import "testing"

func TestSimpleQuery(t *testing.T) {
	table := []struct{ query, expected string }{
		{
			Select(Raw("Users"), Star).Where(Eq(Raw("name"), Raw("'Jimmy'"))).OrderBy(Random(nil)).SQL(MySQL),
			"SELECT * FROM Users WHERE name = 'Jimmy' ORDER BY RAND()",
		},
	}
	for _, tt := range table {
		if tt.query != tt.expected {
			t.Errorf("failed to render correct SQL from query => wanted: %q, got: %q", tt.expected, tt.query)
		}
	}
}
