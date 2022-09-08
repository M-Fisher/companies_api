package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/M-Fisher/companies_api/app"
	"github.com/M-Fisher/companies_api/app/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "companies-api",
	Short: "A service for managing Companies",
	Long:  `A service for managing Companies`,
}

func init() {
	cfgFile := os.Getenv("SERVICE_CONF")
	if cfgFile != "" {
		err := godotenv.Load(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	rootCmd.AddCommand(runServerCmd)
}

func Execute() {
	if len(os.Args[1:]) == 0 {
		// run server by default
		os.Args = append(os.Args, runServerCmd.Use)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var runServerCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func runServer() {
	cfg := config.NewFromEnv()
	application := app.New(cfg)

	if err := application.Run(); err != nil {
		panic(err)
	}

}
