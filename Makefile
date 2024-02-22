all: d-build-page-parser d-compose

d-build-page-parser: .
	docker build -t versta-page-parser:local ./landing-page-parser

d-compose:
	docker compose up