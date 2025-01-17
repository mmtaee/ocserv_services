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
swag init -g cmd/main.go --parseDependency

env $(cat .env | xargs) go run cmd/main.go

go run cmd/main.go -debug -drop 

go run cmd/main.go -debug -migrate
```

# develop & Deploy
```bash
# API service
POSTGRES_HOST=127.0.0.1 go run cmd/main.go -debug -drop

sudo docker build -t ocserv:api .

go build  -o build/ocserv_api cmd/main.go  

sudo docker run -it --rm -v "./build:/app" \
    -v "./.volumes/ocserv:/etc/ocserv" \
    -v "/tmp/ocserv:/var/log/ocserv" \
    --env-file=.env -p "8080:8080" -p "20443:443" \
    --link ocserv-postgres:ocserv-postgres \
    --name ocserv_api --privileged ocserv:api
 
ocpasswd -c /etc/ocserv/ocpasswd USERNAME

sudo docker exec -it ocserv_api bash 

echo -e "1234\n1234\n" | ocpasswd -c /etc/ocserv/ocpasswd test
     
#Log Service
sudo docker build -t ocserv:log_processor .

sudo docker run -it --rm \
    -v "/tmp/ocserv:/var/log/ocserv" \
    -e "LOG_FILE=/var/log/ocserv/ocserv.log"\
    --env-file=.env\
    --link ocserv-postgres:ocserv-postgres \
    --name ocserv_log_processor ocserv:log_processor   

#Log Broadcaster
sudo docker build -t ocserv:log_broadcaster .

sudo docker run -it --rm \
    -p "8081:8080"\
    -v "/tmp/ocserv:/var/log/ocserv" \
    -e "LOG_FILE=/var/log/ocserv/ocserv.log"\
    --name ocserv_log_broadcaster ocserv:log_broadcaster
```

```text
worker[test]: 172.17.0.1 worker-auth.c:1731: failed authentication for 'test'

main[test]:172.17.0.1:55064 user logged in

main[test]:192.168.100.21:41448 user disconnected (reason: user disconnected, rx: 483392, tx: 2330)

main[test]:172.17.0.1:56906 user disconnected (reason: server disconnected, rx: 378308, tx: 944)
```