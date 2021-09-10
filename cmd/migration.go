package cmd

import (
	"fmt"
	"github.com/ShaDream/library-api-practice/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

//todo: обработать исключения о миграции
var migrationCommand = &cobra.Command{
	Use:   "migrate",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return nil
	},
}

var upMigrationCommand = &cobra.Command{
	Use:   "up",
	Short: "",
	Long:  "",
	Args:  migrationArgsChecker,
	RunE:  executeUpMigration,
}

var downMigrationCommand = &cobra.Command{
	Use:   "down",
	Short: "",
	Long:  "",
	Args:  migrationArgsChecker,
	RunE:  executeDownMigration,
}

var sourcePath string

func init() {
	migrationCommand.AddCommand(upMigrationCommand, downMigrationCommand)
	migrationCommand.PersistentFlags().StringVarP(&sourcePath, "path", "p", "", "path to migration folder")
}

func migrationArgsChecker(cmd *cobra.Command, args []string) error {
	if err := cobra.RangeArgs(0, 1)(cmd, args); err != nil {
		return err
	}
	if len(args) == 1 {
		number, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("argument should be a number")
		}
		if number == 0 {
			return fmt.Errorf("number should not be 0")
		}
	}
	return nil
}

func executeUpMigration(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		// skip error checking because it checks in args validation
		number, _ := strconv.Atoi(args[0])
		return executeMigration(number, sourcePath)
	} else {
		return executeFullMigration(false, sourcePath)
	}
}

func executeDownMigration(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		// skip error checking because it checks in args validation
		number, _ := strconv.Atoi(args[0])
		return executeMigration(-number, sourcePath)
	} else {
		return executeFullMigration(true, sourcePath)
	}
}

func getMigration(path string) (*migrate.Migrate, error) {
	tx := database.GetTransaction()
	db, err := tx.DB()
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		viper.GetString("db.name"), driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func executeFullMigration(isDown bool, path string) error {
	m, err := getMigration(path)
	if err != nil {
		return err
	}

	if isDown {
		err = m.Down()
	} else {
		err = m.Up()
	}

	if err != nil {
		return err
	}

	return nil
}

func executeMigration(steps int, path string) error {
	m, err := getMigration(path)
	if err != nil {
		return err
	}

	err = m.Steps(steps)

	switch err {
	case os.ErrNotExist:
		fmt.Println("Nothing to migrate")
		return nil
	case nil:
		break
	default:
		return err
	}

	version, d, err := m.Version()
	if err != nil {
		return err
	}
	fmt.Printf("Successfully migrate. Current version: %d, isDirty: %t \n", version, d)
	return nil
}
