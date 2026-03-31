package main

import (
	"context"
	"fmt"
	"github.com/gokulrajk-dev/funk-cli-demo/commands"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "funk",
		Usage: "suite of useful tools for pesky problems",
		Commands:commands.AvailableCommands,
		OnUsageError: func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
			if strings.Contains(err.Error(), "invalid value") {
				fmt.Println("❌ Invalid input: seconds must be an integer (e.g., --s 10)")
				return nil
			}
			return err
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
