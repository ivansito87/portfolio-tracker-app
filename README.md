# Full-Stack Application Development & Deployment

## Project Overview
A full-stack web application built with React, Golang, and PostgreSQL, deployed on AWS with a focus on scalability, automation, and DevOps best practices.

## Architecture Overview

### Technologies & Services
- **Frontend**: React, TypeScript, Material UI (EC2 + GitHub Actions CI/CD)
- **Backend**: Golang with Gin Framework (AWS EKS)
- **Database**: AWS RDS PostgreSQL
- **Infrastructure**: AWS (EKS, EC2, RDS, ECR, IAM, Load Balancer)
- **Containerization**: Docker, Kubernetes (EKS)
- **CI/CD**: GitHub Actions

## Implementation Steps

### 1. Frontend Development & Deployment
- React application with React Router and Material UI
- EC2 hosting with Nginx
- Automated deployment via GitHub Actions

### 2. Backend Development & Deployment
- REST API using Golang and Gin
- Docker containerization
- AWS ECR for image registry
- EKS deployment with Load Balancer

### 3. Database Setup
- AWS RDS PostgreSQL
- SQL schema migrations
- Secure environment configuration

### 4. AWS Kubernetes (EKS) Setup
- EKS Cluster configuration
- Kubernetes Deployment and Service
- IAM roles and security groups
- Network configuration

### 5. CI/CD Pipeline
- **Backend**:
  - Automated image building
  - ECR push
  - EKS deployment
- **Frontend**:
  - Automated builds
  - EC2 deployment

### 6. Testing & Validation
- API accessibility verification
- Frontend-backend integration
- CORS and networking resolution
- End-to-end testing

## Final Outcome
✅ Full-stack application deployed on AWS
✅ Automated CI/CD for frontend and backend
✅ Scalable Kubernetes architecture
✅ Secure containerized backend
✅ Optimized infrastructure following DevOps best practices
