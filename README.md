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
- On successful push to `main`, it builds the Docker image. If you add the following repository secrets, it will also deploy to Cloud Run:

- `GCP_SA_KEY` – JSON service account key with permissions to deploy to Cloud Run
- `GCP_PROJECT` – your Google Cloud project id
- `CLOUD_RUN_SERVICE` – desired Cloud Run service name
- `CLOUD_RUN_REGION` – region, e.g. `us-central1`

To create `GCP_SA_KEY`, follow the Google Cloud docs to create a service account with the `roles/run.admin` and `roles/iam.serviceAccountUser` roles and download the JSON key. Add the JSON content to the `GCP_SA_KEY` secret.

### Step-by-step: create service account and set GitHub secrets

1. Enable APIs in your GCP project:

```bash
gcloud services enable run.googleapis.com containerregistry.googleapis.com
```

2. Create a service account and grant it the required roles:

```bash
PROJECT=your-gcp-project-id
SA_NAME=github-actions-deployer
gcloud iam service-accounts create $SA_NAME --project $PROJECT --display-name "GH Actions deployer"

# Grant roles required to deploy to Cloud Run and push images
gcloud projects add-iam-policy-binding $PROJECT \
  --member="serviceAccount:$SA_NAME@$PROJECT.iam.gserviceaccount.com" \
  --role="roles/run.admin"
gcloud projects add-iam-policy-binding $PROJECT \
  --member="serviceAccount:$SA_NAME@$PROJECT.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"
gcloud projects add-iam-policy-binding $PROJECT \
  --member="serviceAccount:$SA_NAME@$PROJECT.iam.gserviceaccount.com" \
  --role="roles/storage.admin"

# (Optional) If you're using Artifact Registry instead of GCR, add the appropriate roles.
```

3. Create and download a JSON key for the service account:

```bash
gcloud iam service-accounts keys create key.json \
  --iam-account=$SA_NAME@$PROJECT.iam.gserviceaccount.com --project $PROJECT
```

4. In your GitHub repository settings → Secrets → Actions, add the following secrets:

- `GCP_SA_KEY` — the contents of `key.json` (paste the whole JSON)
- `GCP_PROJECT` — your GCP project id
- `CLOUD_RUN_SERVICE` — desired Cloud Run service name (e.g., `testservice`)
- `CLOUD_RUN_REGION` — region (e.g., `us-central1`)

After that, pushes to `main` will trigger the workflow which builds, pushes the image to `gcr.io/<PROJECT>/testservice:<sha>`, and deploys it to Cloud Run.

### Artifact Registry (optional, recommended)

If you prefer Artifact Registry (recommended for new projects), configure the workflow to push there by adding the `ARTIFACT_REGISTRY_REPO` secret (and ensure `GCP_REGION` is set as well). When `ARTIFACT_REGISTRY_REPO` is present the workflow will push to:

```
${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT}/${ARTIFACT_REGISTRY_REPO}/testservice:<sha>
```

Steps to create an Artifact Registry repository and grant permissions:

```bash
PROJECT=your-gcp-project-id
REGION=us-central1
REPO=testservice-repo

# Create a Docker repository in Artifact Registry
gcloud artifacts repositories create $REPO --repository-format=docker --location=$REGION --description="Docker repo for testservice" --project=$PROJECT

# Grant the service account permission to write images
gcloud projects add-iam-policy-binding $PROJECT \
  --member="serviceAccount:github-actions-deployer@$PROJECT.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

# (Also ensure roles/storage.admin or necessary permissions for pushing images are granted as per your setup.)
```

Add the following repository secrets in GitHub:

- `ARTIFACT_REGISTRY_REPO` — the Artifact Registry repo name (e.g., `testservice-repo`)
- `GCP_REGION` — region used for the Artifact Registry repo (e.g., `us-central1`)

With those secrets present the Actions workflow will push to Artifact Registry and deploy the pushed image to Cloud Run.

## Notes
- This is intentionally minimal. Consider adding request logging, structured logs, and metrics when you extend it.
