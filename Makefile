cogigo-prod:
	rm .env; touch .env; echo ENVIRONMENT = "PRODUCTION" >> .env
	# Uncomment the following line to enable migrations
	# go run cmd/migrations/main.go
	go run cmd/main.go

cogigo-dev:
	docker compose down
	rm .env; touch .env; echo ENVIRONMENT = "DEVELOPMENT" >> .env
	# Uncomment the following line to enable migrations
	# go run cmd/migrations/main.go
	docker compose up -d
	go run cmd/main.go
