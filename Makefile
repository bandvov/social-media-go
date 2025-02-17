export NOTIFICATIONS_VERSION=$(cat ./notifications/VERSION)

up:
	docker-compose up --build
