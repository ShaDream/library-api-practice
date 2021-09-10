package cmd

import (
	"github.com/ShaDream/library-api-practice/database"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var feedCommand = &cobra.Command{
	Use:   "feed",
	Short: "",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	RunE:  executeFeed,
}

func executeFeed(cmd *cobra.Command, args []string) error {
	tx := database.GetTransaction()
	for _, arg := range args {
		file, err := ioutil.ReadFile(arg)
		if err != nil {
			tx.Rollback()
			return err
		}
		result := tx.Exec(string(file))
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}
	err := tx.Commit().Error
	return err
}
