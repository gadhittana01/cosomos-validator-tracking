# Cosmos Validator Delegation Tracking System

This system continuously polls validator delegation data every hour from the Cosmos blockchain, stores the complete delegation state, and provides APIs to retrieve historical changes in delegation for both validators and delegators. It also aggregates daily delegation data and supports efficient retrieval and tracking.

## 🎯 Motivation

Validator delegations are essential in understanding validator performance and user behavior in Cosmos' Proof of Stake network. By monitoring delegation changes, we gain insights into staking patterns, validator popularity, and potential risks in validator uptime or credibility.


## System Architecture

### Components

1. **Data Collection Services**
   - Hourly collector for real-time delegation data
   - Daily aggregator for summarized delegation metrics

2. **Storage Layer**
   - PostgreSQL database for persistent storage of delegation data
   - Redis for caching frequently accessed data

3. **API Services**
   - RESTful endpoints for querying validator delegation data
   - Swagger documentation for API reference

### Technology Stack

- **Backend**: Go with Chi router
- **Database**: PostgreSQL with pgx driver
- **Caching**: Redis
- **API Documentation**: Swagger/OpenAPI
- **Containerization**: Docker and Docker Compose

## Features

- Hourly collection of validator delegation data
- Daily aggregation of delegation metrics
- Caching layer for improved performance
- Comprehensive error handling with retries
- RESTful API for data access
- Pagination support for all endpoints

## API Endpoints

API Swagger documentation is accessible at: http://127.0.0.1:8000/v1/validator/docs/#/

### Validator Delegations

- **GET /api/v1/validators/{validatorAddress}/delegations/hourly**
  - Retrieves hourly snapshots of delegations for a specific validator
  - Supports pagination

- **GET /api/v1/validators/{validatorAddress}/delegations/daily**
  - Retrieves daily aggregated delegation data for a specific validator
  - Supports pagination

- **GET /api/v1/validators/{validatorAddress}/delegator/{delegatorAddress}/history**
  - Retrieves the delegation history for a specific delegator to a validator
  - Supports pagination and sorting

### Scheduler Endpoints

- **POST /api/v1/scheduler/validator/hourly**
  - Triggers the hourly collection of validator delegation data

- **POST /api/v1/scheduler/validator/daily**
  - Triggers the daily aggregation of validator delegation data

## Error Handling and Resilience

The system implements comprehensive error handling mechanisms:

- **Transaction Retries**: Automatic retries for failed transactions
- **Recovery Middleware**: Panic recovery middleware to prevent service crashes
- **Contextual Timeout**: Context-based timeouts for external API calls

## Caching Strategy

The system uses Redis for caching with the following features:

- **Cache Invalidation**: Automatic cache invalidation after data updates
- **Prefix-based Clearing**: Ability to clear cache by key prefixes
- **Configurable TTL**: Time-to-live configuration for cached data

## Getting Started

### Prerequisites

- Go 1.24.2 or later
- PostgreSQL
- Redis
- Docker and Docker Compose (optional)

### Running with Docker

```bash
# Build and start services
docker-compose up -d

# Check service status
docker-compose ps
```

### Running Tests

```bash
# Run all tests
make test
```

## Deployment
To deploy this project in a production-ready environment using Google Cloud Platform (GCP), follow the steps below:

### Set Up CI/CD Pipeline
Create a CI/CD pipeline (using GitHub Actions, GitLab CI, or Cloud Build) to automate the deployment process.

Example Workflow Steps:
- Lint and test your code.
- Build Docker image and push to Artifact Registry.
- Deploy to Cloud Run, Instance Group, or GKE.

## Configure Secrets
Use Google Secret Manager to store sensitive information such as:
- Database credentials
- API keys
- Configurable Variables

## Create a Dedicated Service Account
Create and attach the service account to your GCP deployment.

## Deploy the Application
Choose one of the following options based on your architecture preference:
- Cloud Run (Recommended for Simplicity)
- Instance Group (for custom VM workloads)
- GKE (for container orchestration)
