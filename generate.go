package main

//go:generate migrate create -ext sql -dir sql_config/migration -seq xxx

//go:generate migrate -path sql_config/migration -database 'mysql://root:root@tcp(localhost:3306)/cl_db' up

//go:generate sqlc -f sql_config/sqlc.yaml generate

///go:generate sqlboiler -c sql_config/sqlboiler.yaml
