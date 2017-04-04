package q

import "testing"

func TestSimpleQuery(t *testing.T) {
	table := []struct {
		query    Query
		expected string
	}{
		{
			Select(Raw("Users"), Star).Where(Eq(Raw("name"), Raw("'Jimmy'"))).OrderBy(Random(nil)),
			"SELECT * FROM Users WHERE name = 'Jimmy' ORDER BY RAND()",
		},
		{
			Select(Raw("Users"), Star).GroupBy(Raw("User.name")).Having(GreaterThan(Count(Raw("Users.id")), Raw("5"))),
			"SELECT * FROM Users GROUP BY User.name HAVING COUNT(Users.id) > 5",
		},
		{
			Select(Raw("Users"), Star).Having(GreaterThan(Count(Raw("Users.id")), Raw("5"))),
			"SELECT * FROM Users",
		},
	}
	for _, tt := range table {
		sql := MySQL.RenderSQL(tt.query)
		if sql != tt.expected {
			t.Errorf("failed to render SQL => wanted: %q, got: %q", tt.expected, sql)
		}
	}
}

func TestQuerySharing(t *testing.T) {
	baseQuery := Select(Raw("Users"), Star).Where(Eq(Raw("name"), Raw("'Jimmy'")))
	baseQueryExpected := "SELECT * FROM Users WHERE name = 'Jimmy'"
	extendedQuery := baseQuery.Where(Eq(Raw("id"), Raw("1")))
	extendedQueryExpected := "SELECT * FROM Users WHERE name = 'Jimmy' AND id = 1"

	baseQueryRendered := MySQL.RenderSQL(baseQuery)
	extendedQueryRendered := MySQL.RenderSQL(extendedQuery)

	if baseQueryRendered != baseQueryExpected {
		t.Errorf("base query failed to render SQL => wanted: %q, got: %q", baseQueryExpected, baseQueryRendered)
	}
	if extendedQueryRendered != extendedQueryExpected {
		t.Errorf("base query failed to render SQL => wanted: %q, got: %q", extendedQueryExpected, extendedQueryRendered)
	}
}
