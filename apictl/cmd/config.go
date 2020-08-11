package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Save settings in config file.",
	Long: `Save settings in config file.

Example:

apictl -s 127.0.0.1:8080 config
`,
	Run: func(cmd *cobra.Command, args []string) {
		for k, v := range viper.AllSettings() {
			log.Printf("viper: %s = %v", k, v)
		}
		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
