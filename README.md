# q

the search for a flexible SQL query composition library for Go

the goal is to enable writing something like this:

```go
query := q.
	Select(Users, q.Star).
	Where(q.Eq(Users.name, q.Str("Jimmy"))).
	OrderBy(q.Random(nil)).
	Limit(1)

println(query.SQL(q.MySQL))
println(query.SQL(q.PostgreSQL))

// Output:
// SELECT * FROM `Users` WHERE `Users`.`name` = 'Jimmy' ORDER BY RAND() LIMIT 1
// SELECT * FROM "Users" WHERE "Users"."name" = 'Jimmy' ORDER BY RANDOM() LIMIT 1
```
