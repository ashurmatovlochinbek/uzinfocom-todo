migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_car_db?sslmode=disable -path migrations up

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_car_db?sslmode=disable -path migrations down

migrate_create:
	migrate create -ext sql -dir migrations/ -seq create_tables_mg

run:
	go run ./cmd/main.go