#!/bin/bash

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