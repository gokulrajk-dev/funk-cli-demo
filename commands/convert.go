//go:build convert || all
package commands

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"strconv"
	"strings"
)

func init() {
	AvailableCommands = append(AvailableCommands, ConvertCommand())
}


func ConvertCommand() *cli.Command {
	return &cli.Command{
		Name:    "conv",
		Suggest: true,
		Usage:   "converts various units",
		Action:  Converto,
		CustomHelpTemplate: `NAME: 
	{{.FullName}} - {{.Usage}}
	 
USAGE:
   {{.FullName}} [INPUT FLAGS] TARGET [CONVERSION FLAGS]	
{{if .VisibleFlagCategories}}  
FLAGS:  
	{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}  
FLAGS:  
	{{range .VisibleFlags}}   {{.}}  
	{{end}}{{end}}  	
EXAMPLES:
   1. To convert kilometer to meter 
      $ {{.FullName}} -k[Input Flag] 7 -tm[Conversion Flag] 
`,
		Flags: []cli.Flag{
			// distance converters.
			&cli.Float64Flag{
				Name:     "miles",
				Category: "Distance",
				Aliases:  []string{"M"},
				Usage:    "Enter values in miles",
			},
			&cli.Float64Flag{
				Name:     "km",
				Category: "Distance",
				Aliases:  []string{"k"},
				Usage:    "Enter values in kilometers",
			},
			&cli.Float64Flag{
				Name:     "m",
				Category: "Distance",
				Aliases:  []string{"m"},
				Usage:    "Enter values in meters",
			},
			&cli.BoolFlag{
				Name:     "to-km",
				Category: "Distance",
				Aliases:  []string{"tk"},
				Usage:    "convert to kilometer",
			},
			&cli.BoolFlag{
				Name:     "to-miles",
				Category: "Distance",
				Aliases:  []string{"tM"},
				Usage:    "Convert to miles",
			},
			&cli.BoolFlag{
				Name:     "to-m",
				Category: "Distance",
				Aliases:  []string{"tm"},
				Usage:    "convert to meters",
			},
			&cli.BoolFlag{
				Name:     "to-cm",
				Category: "Distance",
				Aliases:  []string{"tc"},
				Usage:    "convert to centimeters",
			},

			// Weight converters
			&cli.Float64Flag{
				Name:     "lbs",
				Category: "Weight",
				Aliases:  []string{"p"},
				Usage:    "Enter values in pounds",
			},
			&cli.BoolFlag{
				Name:     "to-kg",
				Category: "Weight",
				Aliases:  []string{"tw"},
				Usage:    "Convert to kg",
			},
			&cli.Float64Flag{
				Name:     "kg",
				Category: "Weight",
				Aliases:  []string{"w"},
				Usage:    "Enter values in kilograms",
			},
			&cli.BoolFlag{
				Name:     "to-lbs",
				Category: "Weight",
				Aliases:  []string{"tp"},
				Usage:    "Convert to pounds",
			},
			&cli.BoolFlag{
				Name:     "to-gm",
				Category: "Weight",
				Aliases:  []string{"tg"},
				Usage:    "Convert to grams",
			},

			// temperature converters
			&cli.Float64Flag{
				Name:     "celsius",
				Category: "Temperature",
				Aliases:  []string{"c"},
				Usage:    "Enter values in degree celsius",
			},
			&cli.Float64Flag{
				Name:     "fahrenheit",
				Category: "Temperature",
				Aliases:  []string{"f"},
				Usage:    "Enter values in Fahrenheit",
			},
			&cli.BoolFlag{
				Name:     "to-c",
				Category: "Temperature",
				Aliases:  []string{"tC"},
				Usage:    "Convert fahrenheit to celsius",
			},
			&cli.BoolFlag{
				Name:     "to-f",
				Category: "Temperature",
				Aliases:  []string{"tf"},
				Usage:    "Convert celsius to fahrenheit",
			},

			// Binary to hexadecimal conversion
			&cli.StringSliceFlag{
				Name:     "binary",
				Category: "Number system",
				Aliases:  []string{"b"},
				Usage:    "Enter values in binary in quotes",
			},
			&cli.BoolFlag{
				Name:     "to-hex",
				Category: "Number system",
				Aliases:  []string{"tH"},
				Usage:    "Convert binary numbers to hexadecimal",
			},
			&cli.StringSliceFlag{
				Name:     "hex",
				Category: "Number system",
				Aliases:  []string{"H"},
				Usage:    "Enter values in hexadecimal",
			},
			&cli.BoolFlag{
				Name:     "to-binary",
				Category: "Number system",
				Aliases:  []string{"tb"},
				Usage:    "Convert hexadecimal numbers to binary number",
			},
		},
		//	OnUsageError: ErrorHandle,
	}
}

func Converto(ctx context.Context, cmd *cli.Command) error {

	// distance conversion
	meters := cmd.Float64("m")
	if cmd.IsSet("m") {
		if cmd.Bool("to-km") {
			km := meters / 1000
			fmt.Printf("\n %.2f m <-> %.2f km \n", meters, km)
		}
		if cmd.Bool("to-cm") {
			cm := meters * 100
			fmt.Printf("\n %.2f m <-> %.2f cm \n", meters, cm)
		}
		if cmd.Bool("to-miles") {
			miles := meters * 0.00062137
			fmt.Printf("\n %.2f m <-> %.2f miles \n", meters, miles)
		}
	}
	km := cmd.Float64("km")
	if cmd.IsSet("km") {
		if cmd.Bool("to-miles") {
			miles := km * 0.621371
			fmt.Printf("\n %.2f km <-> %.2f miles\n", km, miles)
		}
		if cmd.Bool("to-m") {
			meters := km * 1000
			fmt.Printf("\n %.2f km <-> %.2f meters\n", km, meters)
		}
		if cmd.Bool("to-cm") {
			cm := km * 100000
			fmt.Printf("\n %.2f km <-> %.2f centimeters\n", km, cm)
		}
	}
	miles := cmd.Float64("miles")
	if cmd.IsSet("miles") {
		if cmd.Bool("to-km") {
			kim := miles / 0.621371
			fmt.Printf("\n %.2f miles <-> %.2f km\n", miles, kim)
		}
		if cmd.Bool("to-m") {
			meters := miles * 1609.344
			fmt.Printf("\n %.2f miles <-> %.2f meters\n", miles, meters)
		}
		if cmd.Bool("to-cm") {
			cm := miles * 160934.4
			fmt.Printf("\n %.2f miles <-> %.2f centimeters\n", miles, cm)
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

	// Number system conversion - binary to hexadecimal

	binary := cmd.StringSlice("binary")
	if cmd.IsSet("binary") {
		if cmd.Bool("to-hex") {
			hexSlice := []string{}
			for _, binstr := range binary {
				if val, err := strconv.ParseInt(binstr, 2, 64); err == nil {
					hexstr := fmt.Sprintf("0x%X", val)
					fmt.Printf("%s (binary) = %s (hex) \n", binstr, hexstr)
					hexSlice = append(hexSlice, hexstr)
				}
			}
			fmt.Printf("\n%s \n", hexSlice)
		}

	}
	hex := cmd.StringSlice("hex")
	if cmd.IsSet("hex") {
		if cmd.Bool("to-binary") {
			binSlice := []string{}
			for _, hexstr := range hex {
				cleanhex := strings.TrimPrefix(hexstr, "0x")
				if val, err := strconv.ParseInt(cleanhex, 16, 64); err == nil {
					binStr := strconv.FormatInt(val, 2)
					fmt.Printf("%s (hex) = %s (binary) \n", hexstr, binStr)
					binSlice = append(binSlice, binStr)
				}
			}
			fmt.Printf("\n%s\n", binSlice)
		} else {
			fmt.Printf("Try -tb?\n")
		}

	}

	return nil
}
