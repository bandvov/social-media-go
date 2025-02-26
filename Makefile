export NOTIFICATIONS_VERSION=$(cat ./notifications/VERSION)
export USERS_VERSION=$(cat ./users/VERSION)
export POSTS_VERSION=$(cat ./posts/VERSION)
export FOLLOWERS_VERSION=$(cat ./followers/VERSION)
export COMMENTS_VERSION=$(cat ./comments/VERSION)
export REACTIONS_VERSION=$(cat ./reactions/VERSION)


up:
	docker-compose up --build
