package postgres_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/toadharvard/goxkcd/internal/config"
	"github.com/toadharvard/goxkcd/internal/entity"
	repository "github.com/toadharvard/goxkcd/internal/repository/comix/postgres"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupPostgresContainer(ctx context.Context, t *testing.T) *postgres.PostgresContainer {
	dbName := "goxkcd"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	return postgresContainer
}

func MigratePostgres(t *testing.T, migrationsPath string, DSN string) {
	m, _ := migrate.New(
		"file://"+migrationsPath,
		DSN,
	)
	_ = m.Up()
}

func TestRepo(t *testing.T) {
	// Container setup
	ctx := context.Background()
	container := SetupPostgresContainer(ctx, t)
	DSN, _ := container.ConnectionString(ctx, "sslmode=disable")
	slog.Info("Connected to Postgres", "dsn", DSN)
	// Load config
	cfg, err := config.New(config.DefaultConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	// Migrate
	MigratePostgres(t, cfg.Postgres.Migrations, DSN)
	// Test
	repo, err := repository.New(DSN)

	if err != nil {
		t.Fatal(err)
	}

	t.Run(
		"Test insert", func(t *testing.T) {
			comics := []entity.Comix{
				{
					ID:       1,
					URL:      "https://xkcd.com/1/img/",
					Keywords: []string{"test1", "test2"},
				},
				{
					ID:       2,
					URL:      "https://xkcd.com/2/img/",
					Keywords: []string{"test1", "test4", "test5"},
				},
			}

			if size, _ := repo.Size(); size != 0 {
				t.Fatalf("expected 0, got %d", size)
			}

			err = repo.BulkInsert(comics)
			assert.NoError(t, err)

			if size, _ := repo.Size(); size != 2 {
				t.Fatalf("expected 2, got %d", size)
			}

			comixFromDB, err := repo.GetAll()
			assert.NoError(t, err)
			assert.Equal(t, comics, comixFromDB)
		},
	)
}
