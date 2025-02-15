package cmd

import (
	"log"
	"merchshop/internal/server"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var (
	port          string
	dbSource      string
	migrationsDir string

	rootCmd = &cobra.Command{
		Use:   "merch",
		Short: "Avito merch app",
		Run:   Run,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&dbSource, "db_source",
		"postgres://root:secret@localhost:5432/example?sslmode=disable",
		"Postgres database connection string")
	rootCmd.PersistentFlags().StringVar(&migrationsDir, "migrations",
		"migrations",
		"Path to the db migrations folder")
	rootCmd.PersistentFlags().StringVar(&port, "port", "8080", "HTTP Server port")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func Run(cmd *cobra.Command, args []string) {
	s := server.NewServer("0.0.0.0:" + port)
	log.Fatal(s.ListenAndServe())
}
