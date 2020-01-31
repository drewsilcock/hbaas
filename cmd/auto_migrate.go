package cmd

import (
	"github.com/drewsilcock/hbaas-server/migrations"
	"github.com/golang-migrate/migrate/v4"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	rootCmd.AddCommand(autoMigrateCmd)
}

var autoMigrateCmd = &cobra.Command{
	Use:   "auto-migrate",
	Short: "Migrate DB schema",
	Long:  "Automatically migrate database schema to the latest version.",
	Run:   autoMigrate,
}

func autoMigrate(cmd *cobra.Command, args []string) {
	assetSource := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	migrationData, err := bindata.WithInstance(assetSource)
	if err != nil {
		log.Fatal("Unable to read migration data:", err)
	}
	migrator, err := migrate.NewWithSourceInstance(
		"go-bindata",
		migrationData,
		viper.GetString("POSTGRES_URL"),
	)
	if err != nil {
		log.Fatal("Unable to read create migrator:", err)
	}
	currentVersion, isDirty, err := migrator.Version()
	if err == migrate.ErrNilVersion {
		log.Println("No migrations currently applied.")
	} else if err != nil {
		log.Fatal("Unable to get migration status:", err)
	} else {
		var status string
		if isDirty {
			status = "dirty"
		} else {
			status = "clean"
		}
		log.Println("Currently", status, "on version", currentVersion)
	}
	log.Println("Migrating database up to latest...")
	if err := migrator.Up(); err != nil {
		log.Fatal("Unable to perform 'up' migrations:", err)
	}
	log.Println("Done")
}
