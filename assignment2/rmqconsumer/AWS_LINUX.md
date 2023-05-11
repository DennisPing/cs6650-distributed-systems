# AWS Linux 2 on ARM architecture

```
sudo yum update
sudo yum -y install docker
sudo usermod -a -G docker ec2-user
newgrp docker
sudo systemctl enable docker.service
sudo systemctl start docker.service
```