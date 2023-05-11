# RabbitMQ consumer

Run multiple containers on 1 virtual machine

## Build and Push up to Docker Hub
```
docker build -t mushufeels/a2-rmqconsumer .
docker push mushufeels/a2-rmqconsumer:latest
```

## Pull down from Docker Hub into VM
```
docker pull mushufeels/a2-rmqconsumer:latest
```

## Create .env file in VM
```
touch ~/rabbitmq.env
echo "RABBITMQ_USERNAME={username}" >> ~/rabbitmq.env
echo "RABBITMQ_PASSWORD={password}" >> ~/rabbitmq.env
echo "RABBITMQ_HOST={ip_address}" >> ~/rabbitmq.env
```

## Create internal Docker network
```
docker network create my_network
```

## Run Container #1
```
docker run -d --name consumer-1 --network=my_network \
--env-file rabbitmq.env \
-e LOG_LEVEL=warn \
mushufeels/a2-rmqconsumer
```

## Run Container #2
```
docker run -d --name consumer-2 --network=my_network \
--env-file rabbitmq.env \
-e LOG_LEVEL=warn \
mushufeels/a2-rmqconsumer
```

## Run Nginx
```
docker run --name nginx-proxy -v ~/nginx.conf:/etc/nginx/nginx.conf:ro -d -p 80:80 --network=my_network nginx
```

## Stop containers
```
docker ps
docker stop {container_name}
```

## Restart stopped containers
```
docker ps -a
docker start {container_name}
```

## View logs of containers
```
docker logs {container_name}
```