# Cloud Run-ready Go service

This repository contains a minimal Go HTTP service ready to run on Cloud Run. It includes a Dockerfile, unit tests, a GitHub Actions workflow for CI (and optional deploy), and local run instructions.

## What you get
- `main.go` + `handlers.go` – simple HTTP server with `/` and `/healthz` endpoints
- `Dockerfile` – multi-stage build for a small runtime image
- `.github/workflows/ci-cd.yml` – runs tests and can deploy to Cloud Run when secrets are set
- `Makefile` – helpers for build, run, docker-build, test

## Local development

Prerequisites:
- Go 1.20+
- Docker (for container testing)

Run locally with Go:

```bash
go run ./
```

Run tests:

```bash
go test ./...
```

Build and run with Docker:

```bash
docker build -t testservice:local .
docker run -p 8080:8080 testservice:local
```

Then open http://localhost:8080/ and http://localhost:8080/healthz

## Deploy to Cloud Run (manual)

Using gcloud:

1. Build and push container to Artifact Registry or Container Registry.

```bash
# Example using Container Registry
docker build -t gcr.io/PROJECT-ID/testservice:latest .
docker push gcr.io/PROJECT-ID/testservice:latest

# Deploy to Cloud Run
gcloud run deploy testservice \
  --image gcr.io/PROJECT-ID/testservice:latest \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated
```

Replace `PROJECT-ID` and region as needed.

## GitHub Actions / CI

The workflow at `.github/workflows/ci-cd.yml` will:
- Run `go test` and `go vet` on PRs and pushes to `main`.
- On successful push to `main`, it builds the Docker image, pushes to Artifact Registry, and deploys to Cloud Run (when secrets are configured).

**Required repository secrets:**

- `GCP_SA_KEY` – JSON service account key with permissions to deploy to Cloud Run and push images
- `GCP_PROJECT` – your Google Cloud project id
- `CLOUD_RUN_SERVICE` – desired Cloud Run service name (e.g., `testservice`)
- `CLOUD_RUN_REGION` – region (e.g., `us-central1`)

### Step-by-step: create service account, Artifact Registry repo, and set GitHub secrets

1. Enable APIs in your GCP project:

```bash
PROJECT=your-gcp-project-id
REGION=us-central1

gcloud services enable run.googleapis.com artifactregistry.googleapis.com --project $PROJECT
```

2. Create an Artifact Registry repository for Docker images:

```bash
gcloud artifacts repositories create testservice \
  --repository-format=docker \
  --location=$REGION \
  --description="Docker images for testservice" \
  --project=$PROJECT
```

3. Create a service account and grant it the required roles:

```bash
SA_NAME=github-actions-deployer
gcloud iam service-accounts create $SA_NAME --project $PROJECT --display-name "GH Actions deployer"

# Grant roles required to deploy to Cloud Run and push images to Artifact Registry
gcloud projects add-iam-policy-binding $PROJECT \
  --member="serviceAccount:$SA_NAME@$PROJECT.iam.gserviceaccount.com" \
  --role="roles/run.admin"
gcloud projects add-iam-policy-binding $PROJECT \
  --member="serviceAccount:$SA_NAME@$PROJECT.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"
gcloud projects add-iam-policy-binding $PROJECT \
  --member="serviceAccount:$SA_NAME@$PROJECT.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"
```

4. Create and download a JSON key for the service account:

```bash
gcloud iam service-accounts keys create key.json \
  --iam-account=$SA_NAME@$PROJECT.iam.gserviceaccount.com --project $PROJECT
```

**Important: Do NOT commit key.json to git. Keep it private.**

5. In your GitHub repository settings → Secrets → Actions, add the following secrets:

- `GCP_SA_KEY` — the contents of `key.json` (paste the whole JSON)
- `GCP_PROJECT` — your GCP project id
- `CLOUD_RUN_SERVICE` — desired Cloud Run service name (e.g., `testservice`)
- `CLOUD_RUN_REGION` — region (e.g., `us-central1`)

After that, pushes to `main` will trigger the workflow which builds, pushes the image to `$REGION-docker.pkg.dev/$PROJECT/testservice/testservice:<sha>`, and deploys it to Cloud Run.

## Notes
- This is intentionally minimal. Consider adding request logging, structured logs, and metrics when you extend it.
