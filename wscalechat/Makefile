defualt: build

build:
	GOOS=linux CGO_ENABLED=0 go build
	docker-compose build

run: build
	docker-compose up --scale server=3


