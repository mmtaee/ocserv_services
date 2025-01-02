# Env
```bash
HOST=0.0.0.0
PORT=8080
ALLOW_ORIGINS=
SECRET_KEY=SECRET_KEY

POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
POSTGRES_NAME=ocserv
POSTGRES_USER=ocserv
POSTGRES_PASSWORD=ocserv
```

# Services
```bash
sudo docker run -d \
  --name ocserv-postgres \
  -e POSTGRES_USER=ocserv \
  -e POSTGRES_PASSWORD=ocserv \
  -e POSTGRES_DB=ocserv \
  -v /home/masoud/Documents/docker-volumes/ocserv/db:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:latest 


swag init -g cmd/main.go

env $(cat .env | xargs) go run cmd/main.go

go run cmd/main.go -debug -drop 

go run cmd/main.go -debug -migrate
```