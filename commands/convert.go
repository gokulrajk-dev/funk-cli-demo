package commands

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func ConvertCommand() *cli.Command {
	return &cli.Command{
		Name:   "convert",
		Usage:  "converts various units",
		Action: Converto,
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  "miles",
				Usage: "Convert kilometers to Miles",
			},
			&cli.Float64Flag{
				Name:  "km",
				Usage: "Convert miles to kilometers",
			},
		},
	}
}

func Converto(ctx context.Context, cmd *cli.Command) error {
	if cmd.IsSet("miles") {
		km := cmd.Float64("miles")
		mil := km * 0.6213712
		fmt.Printf("%.2f km = %.2f miles \n", km, mil)
		return nil
	} else if cmd.IsSet("km") {
		miles := cmd.Float64("km")
		kim := miles / 0.6213712
		fmt.Printf("%.2f miles = %.2f km \n", miles, kim)
		return nil
	} else {
		return fmt.Errorf("use --km or --miles")
	}

}
