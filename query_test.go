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
			t.Errorf("failed to render SQL => wanted: %q, got: %q", tt.expected, tt.query)
		}
	}
}

func TestQuerySharing(t *testing.T) {
	baseQuery := Select(Raw("Users"), Star).Where(Eq(Raw("name"), Raw("'Jimmy'")))
	baseQueryExpected := "SELECT * FROM Users WHERE name = 'Jimmy'"
	extendedQuery := baseQuery.Where(Eq(Raw("id"), Raw("1")))
	extendedQueryExpected := "SELECT * FROM Users WHERE name = 'Jimmy' AND id = 1"

	baseQueryRendered := baseQuery.SQL(MySQL)
	extendedQueryRendered := extendedQuery.SQL(MySQL)

	if baseQueryRendered != baseQueryExpected {
		t.Errorf("base query failed to render SQL => wanted: %q, got: %q", baseQueryExpected, baseQueryRendered)
	}
	if extendedQueryRendered != extendedQueryExpected {
		t.Errorf("base query failed to render SQL => wanted: %q, got: %q", extendedQueryExpected, extendedQueryRendered)
	}
}
