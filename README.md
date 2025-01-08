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
  -p 5432:5432 \
  -e POSTGRES_USER=ocserv \
  -e POSTGRES_PASSWORD=ocserv \
  -e POSTGRES_DB=ocserv \
  -v ./.volumes/db:/var/lib/postgresql/data \
  --restart always \
  postgres:latest 

sudo docker run -d \
  --name rabbitmq-ocserv \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=ocserv \
  -e RABBITMQ_DEFAULT_PASS=ocserv \
  -e RABBITMQ_DEFAULT_VHOST=/ocserv\
  --restart always \
  rabbitmq:management

swag init -g cmd/main.go

env $(cat .env | xargs) go run cmd/main.go

go run cmd/main.go -debug -drop 

go run cmd/main.go -debug -migrate
```

# develop & Deploy
```bash
POSTGRES_HOST=127.0.0.1 go run cmd/main.go -debug -drop

go build  -o build/ocserv_api cmd/main.go  

sudo docker build -t ocserv:api .

sudo docker run -it --rm -v "./build:/app" \
    -v "./.volumes/ocserv:/etc/ocserv" \
    -v "./.volumes/logs:/var/log/ocserv" \
    --env-file=.env -p "8080:8080" -p "20443:443" \
    --link ocserv-postgres:ocserv-postgres \
    --link rabbitmq-ocserv:rabbitmq-ocserv \
    --name ocserv_api --privileged ocserv:api

sudo docker run -it --rm -v "./build:/app" \
    -v "./.volumes/ocserv:/etc/ocserv" \
    -v "./.volumes/logs:/var/log/ocserv" \
    --env-file=.env -p "8080:8080" -p "20443:443" \
    --link ocserv-postgres:ocserv-postgres \
    --name ocserv_api --privileged ocserv:api
 
ocpasswd -c /etc/ocserv/ocpasswd USERNAME

sudo docker exec -it ocserv_api bash 

echo -e "1234\n1234\n" | ocpasswd -c /etc/ocserv/ocpasswd test
     
```

```text
main[test]:192.168.100.21:41448 user disconnected (reason: user disconnected, rx: 483392, tx: 2330)
```