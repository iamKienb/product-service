DB_URL=postgres://postgres:postgres@127.0.0.1:5435/shopify_product_service?sslmode=disable
sqlc:
	sqlc generate

migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-drop:
	migrate -path db/migrations -database "$(DB_URL)" drop

seed:
	migrate create -ext sql -dir db/seeds -seq $(name)

seed-up:
	migrate -path db/seeds -database "$(DB_URL)&x-migrations-table=seed_migrations" up

seed-drop:
	migrate -path db/seeds -database "$(DB_URL)&x-migrations-table=seed_migrations" drop

proto:
	@echo "--- Generating Proto  ---"
	buf generate
	@echo "--- Done! Check your /contact folder ---"

run:
	air

dev: migrate-up run