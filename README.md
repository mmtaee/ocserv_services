# Env
```bash
HOST=0.0.0.0
PORT=8080
ALLOW_ORIGINS=
SECRET_KEY=SECRET_KEY

POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5434
POSTGRES_NAME=ocserv
POSTGRES_USER=ocserv
POSTGRES_PASSWORD=ocserv

RABBIT_MQ_HOST=127.0.0.1
RABBIT_MQ_PORT=5672
RABBIT_MQ_USER=ocserv
RABBIT_MQ_PASSWORD=ocserv
RABBIT_MQ_SECURE=false 
RABBIT_MQ_VHOST=ocserv
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

sudo docker run -d \
  --name ocserv-rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -v rabbitmq_data:/var/lib/rabbitmq \
  -e RABBITMQ_DEFAULT_USER=ocserv \
  -e RABBITMQ_DEFAULT_PASS=ocserv \
  -e RABBITMQ_DEFAULT_VHOST=ocserv \
  rabbitmq
```