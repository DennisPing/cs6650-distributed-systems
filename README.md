# cs6650-distributed-systems

https://gortonator.github.io/bsds-6650/

Done in Go and deployed to Google Cloud Platform for self-learning purposes.

## Create new Cloud Run Project

1. Create a new project in Google Cloud Platform
2. Make sure you have billing enabled
3. Enable the Cloud Run API

**Project name = cs6650-dping**

## Install and Initialize Google Cloud CLI if not already installed
```
brew install google-cloud-sdk
gcloud init
```

## Login to Google Cloud
```
gcloud auth login
```

## Set project as default
```
gcloud config set project cs6650-dping
```

## Give yourself permissions to deploy to Cloud Run

https://cloud.google.com/run/docs/reference/iam/roles#additional-configuration 

1. Go to IAM & Admin > Service Accounts
2. Click on PROJECT_NUMBER-compute@developer.gserviceaccount.com
3. Click the Permissions tab
4. Click the Grant Access button
5. Enter your email (or other's email)
6. In the Select a role dropdown, select the Service Accounts > Service Account User role
7. Click Save

## Configure Docker to use Google Container Registry
```
gcloud auth configure-docker
```