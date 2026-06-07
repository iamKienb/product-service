package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"product-shared-module/alias"
	"product-worker-module/internal/bootstrap/config"
	"product-worker-module/internal/bootstrap/module"

	"github.com/elastic/go-elasticsearch/v8"
	configx "github.com/iamKienb/go-core/config"
)

type MigrationLog struct {
	Filename  string    `json:"filename"`
	AppliedAt time.Time `json:"applied_at"`
}

type aliasMigration struct {
	Alias string
	Glob  string
}

const MigrationIndex = ".migration_history"

func main() {
	ctx := context.Background()

	cfg, err := configx.Loader[config.ProductWorkerConfig]()
	if err != nil {
		log.Fatal("failed to load config")
	}

	infra, err := module.NewInfraModule(ctx, cfg)
	if err != nil {
		log.Fatal("failed connect to infra")
	}
	defer infra.ESService.Close(ctx)

	client := infra.ESService.GetClient()
	ensureMigrationIndex(ctx, client)

	migrations := []aliasMigration{
		{Alias: alias.ProductAlias, Glob: "migrations/products/*.json"},
	}

	for _, migration := range migrations {
		if err := migrateAlias(ctx, infra, client, migration); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("All migrations done")
}

func migrateAlias(ctx context.Context, infra *module.InfraModule, client *elasticsearch.TypedClient, migration aliasMigration) error {
	files, err := filepath.Glob(migration.Glob)
	if err != nil {
		return fmt.Errorf("failed to get list migration files: %w", err)
	}

	sort.Strings(files)
	for _, file := range files {
		filename := filepath.Base(file)
		historyKey := fmt.Sprintf("%s_%s", migration.Alias, filename)

		applied, err := isApplied(ctx, client, historyKey)
		if err != nil {
			return fmt.Errorf("failed to check history file %s: %w", filename, err)
		}
		if applied {
			log.Printf("Skip: %s (applied)", historyKey)
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("cannot read file %s: %w", filename, err)
		}

		log.Printf("Migrating [%s] for alias [%s]", filename, migration.Alias)
		if err = infra.ESService.BootstrapIndex(ctx, migration.Alias, string(content)); err != nil {
			return fmt.Errorf("failed to init bootstrap for %s: %w", migration.Alias, err)
		}

		if err := saveLogHistoryFile(ctx, client, historyKey); err != nil {
			return fmt.Errorf("failed to write log for file %s: %w", filename, err)
		}

		log.Printf("Success: %s", historyKey)
	}

	return nil
}

func isApplied(ctx context.Context, client *elasticsearch.TypedClient, fileName string) (bool, error) {
	res, err := client.Get(MigrationIndex, fileName).Do(ctx)
	if err != nil {
		return false, nil
	}

	return res.Found, nil
}

func saveLogHistoryFile(ctx context.Context, client *elasticsearch.TypedClient, fileName string) error {
	logEntry := MigrationLog{
		Filename:  fileName,
		AppliedAt: time.Now(),
	}
	_, err := client.Index(MigrationIndex).Id(fileName).Document(logEntry).Do(ctx)

	return err
}

func ensureMigrationIndex(ctx context.Context, client *elasticsearch.TypedClient) {
	exists, _ := client.Indices.Exists(MigrationIndex).Do(ctx)
	if !exists {
		_, _ = client.Indices.Create(MigrationIndex).Do(ctx)
		log.Printf("Create index [%s] managed migration history", MigrationIndex)
	}
}
