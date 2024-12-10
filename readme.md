Make file : 
1. migrate create -ext sql -dir db/migrations create_table_users
2. migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/simple_crud_clean_architecture?sslmode=disable" up

