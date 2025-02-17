export VERSION_1 := $(shell cat notifications/VERSION)

up:
	docker-compose up --build
