# Lab 3

Dennis Ping

## Build and Test Docker image

### Build Docker image
```
docker build -t hello-server:lab3 -f Dockerfile.server .
docker build -t hello-client:lab3 -f Dockerfile.client .
```

### Run Docker container locally to test for issues
```
docker run -p 8080:8080 -e PORT=8080 hello-server:lab3
docker run -e HELLO_SERVER_URL=http://host.docker.internal:8080 hello-client:lab3 
```

### Stop local Docker container
```
docker ps
docker stop <container id>
```

## Deploy to Cloud Run

### Tag and push Docker image to Google Container Registry

#### Server
```
docker tag hello-server:lab3 gcr.io/cs6650-dping/hello-server:lab3
docker push gcr.io/cs6650-dping/hello-server:lab3
```

#### Client
```
docker tag hello-client:lab3 gcr.io/cs6650-dping/hello-client:lab3
docker push gcr.io/cs6650-dping/hello-client:lab3
```

## Create new Cloud Run service

#### Server
```
gcloud run deploy hello-server \
    --image gcr.io/cs6650-dping/hello-server:lab3 \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated
```
Record the URL of the Cloud Run service for the server!  
Eg. **https://hello-server-[random]-uc.a.run.app**

#### Client
```
gcloud run deploy hello-client \
    --image gcr.io/cs6650-dping/hello-client:lab3 \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated \
    --set-env-vars="HELLO_SERVER_URL=http://hello-server-[random]-uc.a.run.app"
```

## Cloud Run service configuration

### Server
- Increase server concurrent connections to 100
- Lower request timeout to 120 seconds

### Client
- If you want to cap server's max concurrent connections, then you'll need
  to implement retry logic in the client