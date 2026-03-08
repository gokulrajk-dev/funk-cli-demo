package commands

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func ConvertCommand() *cli.Command {
	return &cli.Command{
		Name:    "conv",
		Suggest: true,
		Usage:   "converts various units",
		Action:  Converto, Flags: []cli.Flag{
			// distance converters.
			&cli.Float64Flag{
				Name:    "miles",
				Aliases: []string{"m"},
				Usage:   "Enter values in miles",
			},
			&cli.Float64Flag{
				Name:    "km",
				Aliases: []string{"k"},
				Usage:   "Enter values in kilometers",
			},
			&cli.BoolFlag{
				Name:    "to-km",
				Aliases: []string{"tk"},
				Usage:   "convert to kilometer",
			},
			&cli.BoolFlag{
				Name:    "to-miles",
				Aliases: []string{"tM"},
				Usage:   "Convert to miles",
			},
			&cli.BoolFlag{
				Name:    "to-meters",
				Aliases: []string{"tm"},
				Usage:   "convert to meters",
			},
			// Weight converters
			&cli.Float64Flag{
				Name:    "lbs",
				Aliases: []string{"p"},
				Usage:   "Enter values in pounds",
			},
			&cli.BoolFlag{
				Name:    "to-kg",
				Aliases: []string{"tw"},
				Usage:   "Convert to kg",
			},
			&cli.Float64Flag{
				Name:    "kg",
				Aliases: []string{"w"},
				Usage:   "Enter values in kilograms",
			},
			&cli.BoolFlag{
				Name:    "to-lbs",
				Aliases: []string{"tp"},
				Usage:   "Convert to pounds",
			},
			&cli.BoolFlag{
				Name:    "to-gm",
				Aliases: []string{"tg"},
				Usage:   "Convert to grams",
			},
			// temperature converters
			&cli.Float64Flag{
				Name:    "celsius",
				Aliases: []string{"c"},
				Usage:   "Enter values in degree celsius",
			},
			&cli.Float64Flag{
				Name:    "fahrenheit",
				Aliases: []string{"f"},
				Usage:   "Enter values in Fahrenheit",
			},
			&cli.BoolFlag{
				Name:    "to-c",
				Aliases: []string{"tc"},
				Usage:   "Convert fahrenheit to celsius",
			},
			&cli.BoolFlag{
				Name:    "to-f",
				Aliases: []string{"tf"},
				Usage:   "Convert celsius to fahrenheit",
			},
		},
		//	OnUsageError: ErrorHandle,
	}
}

func Converto(ctx context.Context, cmd *cli.Command) error {

	// converts kilometer values to miles and vice versa
	km := cmd.Float64("km")
	if cmd.IsSet("km") {
		if cmd.Bool("to-miles") {
			miles := km * 0.621371
			fmt.Printf("\n %.2f km <-> %.2f miles\n", km, miles)
		}
		if cmd.Bool("to-meters") {
			meters := km * 1000
			fmt.Printf("\n %.2f km <-> %.2f meters\n", km, meters)
		}
	}

	miles := cmd.Float64("miles")
	if cmd.IsSet("miles") {
		if cmd.Bool("to-km") {
			kim := miles / 0.621371
			fmt.Printf("\n %.2f miles <-> %.2f km\n", miles, kim)
		}
		if cmd.Bool("to-meters") {
			meters := miles * 1609.344
			fmt.Printf("\n %.2f miles <-> %.2f meters\n", miles, meters)
		}
	}

	// Weight conversion

	lbs := cmd.Float64("lbs")
	if cmd.IsSet("lbs") {
		if cmd.Bool("to-kg") {
			kg := lbs * 0.4535924
			fmt.Printf("\n %.2f lbs <->  %.2f kg\n", lbs, kg)
		}
		if cmd.Bool("to-gm") {
			grams := lbs * 453.5924
			fmt.Printf("\n %.2f lbs <-> %.2f gm\n", lbs, grams)
		}
	}

	kg := cmd.Float64("kg")
	if cmd.IsSet("kg") {
		if cmd.Bool("to-lbs") {
			lbs := kg * 2.204623
			fmt.Printf("\n %.2f kg <-> %.2f lbs\n", kg, lbs)
		}
		if cmd.Bool("to-gm") {
			gm := kg * 100
			fmt.Printf("\n %.2f kg <-> %.2f gm\n", kg, gm)
		}
	}

	// temperature conversion

	celsius := cmd.Float64("celsius")
	if cmd.IsSet("celsius") {
		if cmd.Bool("to-f") {
			fahrenheit := (celsius * 9 / 5) + 32
			fmt.Printf("\n %.2f celsius <-> %.2f fahrenheit\n", celsius, fahrenheit)
		}
	}
	fahrenheit := cmd.Float64("fahrenheit")
	if cmd.IsSet("fahrenheit") {
		if cmd.Bool("to-c") {
			celsius := (fahrenheit - 32) * 5 / 9
			fmt.Printf("\n %.2f fahrenheit <-> %.2f celsius\n", fahrenheit, celsius)
		}
	}

	return nil

}

func ErrorHandle(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
	return nil

}
