#!/bin/bash

migrate -path ./internal/database/migrations -database "postgres://blogadmin:blogadmin@localhost:5432/blog?sslmode=disable" up
