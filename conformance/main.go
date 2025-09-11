package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "mysql-conformance",
		Usage: "Protocol Buffers conformance testing for MySQL protobuf implementation",
		Description: `This tool implements the official Protocol Buffers conformance testing protocol
to validate that the MySQL protobuf implementation correctly handles the same test cases
used to validate other protobuf implementations.

For more information about conformance testing, see:
https://github.com/protocolbuffers/protobuf/tree/main/conformance`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "database",
				Usage:    "Database connection string. Example: user:password@tcp(127.0.0.1:3306)/dbname",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
			&cli.StringFlag{
				Name:  "test",
				Usage: "Test name (for logging/identification)",
			},
			&cli.BoolFlag{
				Name:  "use-legacy-conversion",
				Usage: "Use legacy direct Protobuf↔JSON conversion instead of ProtoNumberJSON intermediate format",
				Value: false,
			},
		},
		Action: runConformanceTest,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runConformanceTest(ctx context.Context, command *cli.Command) error {
	database := command.String("database")
	debug := command.Bool("debug")
	testName := command.String("test")
	useLegacyConversion := command.Bool("use-legacy-conversion")

	if debug && testName != "" {
		log.Printf("Running test: %s", testName)
	}

	if debug {
		if useLegacyConversion {
			log.Printf("Using legacy direct Protobuf↔JSON conversion")
		} else {
			log.Printf("Using ProtoNumberJSON intermediate format (default)")
		}
	}

	// Open database connection
	db, err := sql.Open("mysql", database)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create conformance test handler
	handler := &ConformanceHandler{
		db:                  db,
		debug:               debug,
		useLegacyConversion: useLegacyConversion,
	}

	// Run the conformance test protocol
	return handler.RunProtocol()
}
