include .env

compose-up:
	docker-compose up -d --build

compose-down:
	docker-compose down

docs:
	swag init -g internal/app/app.go --pd
.PHONY: docs

keys:
	openssl genpkey -algorithm RSA -out private.key && \
	openssl rsa -pubout -in private.key -out public.key
