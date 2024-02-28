all: d-build-auth-service d-build-orders-service d-build-page-parser d-compose

d-build-page-parser: .
	docker build -t versta-page-parser:local ./landing-page-parser

d-build-orders-service: .
	docker build -t versta-orders-service:local ./orders-service

d-build-auth-service: .
	docker build -t versta-auth-service:local ./auth-service

d-compose:
	docker compose up