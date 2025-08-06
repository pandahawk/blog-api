#!/bin/bash

migrate -path ./internal/database/migrations -database "postgres://admin:admin@localhost:5432/blog?sslmode=disable" up
