package main

import (
	"context"
	"database/sql"
	"github.com/eiiches/mysql-protobuf-functions/internal/mysqlprofiler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"time"
)

func main() {
	cmd := &cli.Command{
		Name:  "mysql-profiler",
		Usage: "Generates collapsed MySQL stack from performance_schema that can be fed into flamegraph.pl",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "database",
				Usage:    "Database connection string. Example: user:password@tcp(127.0.0.1:3306)/dbname",
				Required: true,
			},
		},
		Action: dumpProfileDataAction,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func dumpProfileDataAction(ctx context.Context, command *cli.Command) error {
	db, err := sql.Open("mysql", command.String("database"))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	mysqlprofiler.RunProfile(db)
	return nil
}
