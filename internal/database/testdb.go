package database

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()

	dbUser := "admin"
	dbPass := "admin"
	dbName := "blog"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPass,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres",
			func(host string, port nat.Port) string {
				return fmt.Sprintf(
					"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
					host, port.Port(), dbUser, dbPass, dbName)
			}).WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.
		GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true})

	require.NoError(t, err)

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate test container: %v", err)
		}
	})

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port.Port(), dbUser, dbPass, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	applyMigrations(db)
	SeedDevData(db)
	require.NoError(t, err)
	return db
}

func applyMigrations(db *gorm.DB) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Println(wd)
	sqlDB, _ := db.DB()
	driver, _ := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../database/migrations",
		"postgres", driver,
	)
	if err != nil {
		panic(err)
	}
	_ = m.Up()
}
