# Lab 6: RabbitMQ

Dennis Ping

## Getting Started

RabbitMQ hosting options:

1. AWS EC2 + Elastic Load Balancer + Elastic Block Storage
2. CloudAMQP

### Create EC2 instances

1. Create EC2 instances running Ubuntu
2. Add env variables to `.zshrc` or `.bashrc` for convenience
   ```
    export RABBITMQ_USERNAME=your_username
    export RABBITMQ_PASSWORD=your_password
    export RABBITMQ_HOST=public_elastic_ipv4
    export EC2_DNS=public_ipv4_dns
   ```
3. SSH into EC2
   ```
   ssh -i "~/ec2-t2micro.pem" ubuntu@${EC2_DNS}
   ```
4. Install Docker
   ```
   https://docs.docker.com/engine/install/ubuntu/
   ```

### Add 'ubuntu' to the Docker group

This is so that we can execute Docker commands without needing "sudo"
```
sudo usermod -a -G docker ubuntu
```

### Push RabbitMQ container from local to Docker Hub

```
docker build -t mushufeels/rabbitmq:latest .

docker push mushufeels/rabbitmq:latest
```

### Run Docker container

```
docker pull mushufeels/rabbitmq:latest

docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=[your_username] -e RABBITMQ_DEFAULT_PASS=[your_password] mushufeels/rabbitmq:latest
```

### Log into the RabbitMQ management console

```
http://[elastic_ip]:15672
```

## Consumer

Todo

## Producer

Todo
