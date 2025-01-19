package seeds

import (
	"database/sql"
	"net/http"
)

func SeedData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Seed(db, "./migrations/create_users_table.sql")
		Seed(db, "./migrations/create_posts_table.sql")
		Seed(db, "./migrations/media_urls_create_table.sql")
		Seed(db, "./migrations/create_reaction_types.table.sql")
		Seed(db, "./migrations/create_reactions_table.sql")
		Seed(db, "./migrations/create_followers_table.sql")
		Seed(db, "./migrations/create_tags_table.sql")
		Seed(db, "./migrations/create_comments_table.sql")

		Seed(db, "./seeds/seed_users.sql")
		Seed(db, "./seeds/seed_posts.sql")
		Seed(db, "./seeds/seed_media_urls.sql")
		Seed(db, "./seeds/seed_reaction_types.sql")
		Seed(db, "./seeds/seed_reactions.sql")
		Seed(db, "./seeds/seed_followers.sql")
		Seed(db, "./seeds/seed_tags.sql")
		Seed(db, "./seeds/seed_comments.sql")

	}
}
