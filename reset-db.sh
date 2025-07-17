#!/bin/bash
echo " Looking for Blog API process..."
PID=$(lsof -ti:8080) # assumes your API runs on port 8080

if [ -n "$PID" ]; then
  echo "Stopping Go API (PID $PID)..."
  kill "$PID"
else
  echo "No running Go API found on port 8080"
fi

echo "Stopping containers and removing volumes..."
docker-compose down -v

echo "Starting containers..."
docker-compose up -d

echo "Done."