#!/bin/bash

# Script to run the full application locally (database, backend, frontend) and open Cypress for testing.
# Prerequisites: Docker, Go, Node.js, npm.
# This setup does not interfere with docker-compose or CI pipelines.

set -e  # Exit on error

echo "Starting local development environment..."

# Start the database if not running
if ! docker ps | grep -q local-postgres-ingw3; then
    echo "Starting database..."
    bash run-local-db.sh
fi

# Set environment variables
export DATABASE_URL="postgresql://postgres:password@localhost:5432/ingsw3_integrated?sslmode=disable"
export REACT_APP_BACKEND_URL="http://localhost:8080"

# Start backend in the background
echo "Starting backend..."
(PORT=8080 cd backend && go run cmd/api/main.go > ../backend.log 2>&1 ) &
BACKEND_PID=$!

# Wait a bit for backend to start
sleep 5

# Start frontend in the background
echo "Starting frontend..."
cd frontend
npm start > frontend.log 2>&1 &
FRONTEND_PID=$!

# Wait a bit for frontend to start
sleep 10

# Open Cypress
echo "Opening Cypress interface..."
npm run cypress:open &
CYPRESS_PID=$!

# Function to cleanup on exit
cleanup() {
    echo "Stopping services..."
    kill $BACKEND_PID $FRONTEND_PID $CYPRESS_PID 2>/dev/null || true
    docker stop local-postgres-ingw3 >/dev/null 2>&1 || true
    exit 0
}

# Trap SIGINT (Ctrl+C) to cleanup
trap cleanup SIGINT

echo "All services started successfully!"
echo "Backend: http://localhost:8080"
echo "Frontend: http://localhost:3000"
echo "Database: localhost:5432"
echo ""
echo "Cypress interface should open shortly. Close with Ctrl+C."

# Wait indefinitely
wait
