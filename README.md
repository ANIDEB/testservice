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

## Notes
- This is intentionally minimal. Consider adding request logging, structured logs, and metrics when you extend it.
