# social-media-go

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
