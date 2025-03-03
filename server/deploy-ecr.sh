#!/bin/bash

# AWS ECR configuration
AWS_REGION="us-east-1"
AWS_ACCOUNT_ID="779846782668"
ECR_REPOSITORY="portfolio-backend"
IMAGE_TAG="latest"

# Login to AWS ECR
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Create ECR repository if it doesn't exist
aws ecr describe-repositories --repository-names $ECR_REPOSITORY || \
    aws ecr create-repository --repository-name $ECR_REPOSITORY

# Build the Docker image
docker build -t $ECR_REPOSITORY:$IMAGE_TAG .

# Tag the image for ECR
docker tag $ECR_REPOSITORY:$IMAGE_TAG $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPOSITORY:$IMAGE_TAG

# Push the image to ECR
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPOSITORY:$IMAGE_TAG

echo "Successfully built and pushed image to ECR" 