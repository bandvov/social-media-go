# social-media-go
## Software requirement specifications (SRS)
- [SRS](docs/SRS.md)
- [registration sequence diagram](docs/registration_sequnce.md)
- [login sequence diagram](docs/login_sequence.md)
- [fetch user posts sequence diagram](docs/fetch_user_posts_sequence.md)
- [registration dataflow diagram](docs/registration_dataflow.md)
- [login dataflow diagram](docs/login_dataflow.md)
## Start server
1. Copy `.env.example`
2. Paste it and rename to `.env`
3. Fill with correct data
4. Open terminal and run `docker compose up` to start database
5. Open another terminal and run `source ./.env && go run ./...` to start server 
## Generate selfsined certificates
```bash
openssl req -x509 -newkey rsa:4096 -keyout certs/key.pem -out certs/cert.pem -days 365 -nodes
```
