#!/bin/bash

# Exit on error
set -e

# Create database if it doesn't exist
psql -U postgres -c "CREATE DATABASE portfolio_db;" 2>/dev/null || true

# Initialize the database schema
psql -U postgres -d portfolio_db -f init_db.sql

# Create application user
psql -U postgres -c "CREATE USER portfolio_user WITH PASSWORD 'your_secure_password';"

# Grant necessary permissions
psql -U postgres -d portfolio_db -c "GRANT CONNECT ON DATABASE portfolio_db TO portfolio_user;"
psql -U postgres -d portfolio_db -c "GRANT USAGE ON SCHEMA public TO portfolio_user;"
psql -U postgres -d portfolio_db -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO portfolio_user;"
psql -U postgres -d portfolio_db -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO portfolio_user;"

# PostgreSQL connection string
PGPASSWORD="your_secure_password" psql -h database-2.crkkai2skkf4.us-east-2.rds.amazonaws.com -U portfolio_user -d portfolio_db -c "
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    date TEXT NOT NULL,
    amount FLOAT NOT NULL,
    type TEXT CHECK (type IN ('Credit', 'Debit')) NOT NULL
);
"

echo "Database initialization complete"

# sudo yum update -y
# sudo yum install -y git gcc make
# sudo amazon-linux-extras enable go1.18
# sudo yum install -y golang


# export DB_USER="portfolio_user"
# export DB_PASSWORD="your_secure_password"
# export DB_HOST="database-2.crkkai2skkf4.us-east-2.rds.amazonaws.com"
# export DB_PORT="5432"
# export DB_NAME="portfolio_db"