# Lab 2

Dennis Ping

## Build and Test Docker image

### Build Docker image
```
docker build -t skier-server:lab2 .
```

### Run Docker container locally to test for issues
```
docker run -p 8080:8080 -e PORT=8080 skier-server:lab2
```

### Stop local Docker container
```
docker ps
docker stop <container id>
```

## Deploy to Cloud Run

### Tag and push Docker image to Google Container Registry
```
docker tag skier-server:lab2 gcr.io/cs6650-dping/skier-server:lab2
docker push gcr.io/cs6650-dping/skier-server:lab2
```

## Create new Cloud Run service
```
gcloud run deploy skier-server --image gcr.io/cs6650-dping/skier-server:lab2 --platform managed --region us-central1 --allow-unauthenticated
```
