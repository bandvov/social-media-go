export NOTIFICATIONS_VERSION=$(cat ./notifications/VERSION)
export USERS_VERSION=$(cat ./users/VERSION)

up:
	docker-compose up --build
