package commands

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func ConvertCommand() *cli.Command {
	return &cli.Command{
		Name:    "convert",
		Suggest: true,
		Usage:   "converts various units",
		Action:  Converto,
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  "miles",
				Usage: "Convert kilometers to Miles",
			},
			&cli.Float64Flag{
				Name:  "km",
				Usage: "Convert miles to kilometers",
			},
			&cli.Float64Flag{
				Name:  "lbs",
				Usage: "Convert kilograms to pounds",
			},
			&cli.Float64Flag{
				Name:  "kg",
				Usage: "Convert pounds to kilograms",
			},
		},
		//	OnUsageError: ErrorHandle,
	}
}

func Converto(ctx context.Context, cmd *cli.Command) error {

	// converts kilometer values to miles and vice versa
	if cmd.IsSet("miles") {
		km := cmd.Float64("miles")
		mil := km * 0.6213712
		fmt.Printf("\n%.2f km = %.2f miles \n", km, mil)
		return nil
	} else if cmd.IsSet("km") {
		miles := cmd.Float64("km")
		kim := miles / 0.6213712
		fmt.Printf("%.2f miles = %.2f km \n", miles, kim)
		return nil
	}

	// Converts kilograms to pounds and vice versa
	if cmd.IsSet("lbs") {
		kilogram := cmd.Float64("lbs")
		pounds := kilogram * 2.20462
		fmt.Printf("\n%.2f kg = %.2f lbs \n", kilogram, pounds)
		return nil
	} else if cmd.IsSet("kg") {
		pounds := cmd.Float64("kg")
		kilogram := pounds / 2.20462
		fmt.Printf("\n%.2f lbs <->  %.2f kg \n", pounds, kilogram)
		return nil
	}
	return nil

	// Converts Celsius to Fahrenheit

}

func ErrorHandle(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
	return nil
}
