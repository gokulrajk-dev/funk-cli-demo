//go:build timer || all
package commands

import (
	"context"
	"fmt"
	"funk/sqldb"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

var colorSuccess = color.New(color.FgGreen)



func TimerCommand() *cli.Command {
	return &cli.Command{
		Name:      "timer",
		Usage:     "Set a countdown timer and show Windows toast when done",
		Action:    TimerSet,
		ArgsUsage: "[task name (optional)]",
		Aliases:   []string{"tm"},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "sec",
				Usage:   "timer duration in seconds",
				Aliases: []string{"s"},
			},
			&cli.IntFlag{
				Name:    "min",
				Usage:   "timer duration in minutes",
				Aliases: []string{"m"},
			},
			&cli.IntFlag{
				Name:  "hr",
				Usage: "timer duration in hours",
			},
			&cli.BoolFlag{
				Name:  "his",
				Usage: "show the timer history",
			},
			&cli.IntFlag{
				Name:    "del",
				Usage:   "use to delete the timer record",
				Aliases: []string{"d", "rm"},
			},
			&cli.BoolFlag{
				Name:  "delete_all",
				Usage: "use to delete all timer record",
			},
		},
	}
}

var digits = map[rune][]string{
	'0': {
		" ███ ",
		"█   █",
		"█   █",
		"█   █",
		" ███ ",
	},
	'1': {
		"  █  ",
		" ██  ",
		"  █  ",
		"  █  ",
		"█████",
	},
	'2': {
		"████ ",
		"    █",
		" ███ ",
		"█    ",
		"█████",
	},
	'3': {
		"████ ",
		"    █",
		" ███ ",
		"    █",
		"████ ",
	},
	'4': {
		"█  █ ",
		"█  █ ",
		"█████",
		"   █ ",
		"   █ ",
	},
	'5': {
		"████ ",
		"█    ",
		"████ ",
		"    █",
		"████ ",
	},
	'6': {
		"████ ",
		"█    ",
		"████ ",
		"█   █",
		" ███ ",
	},
	'7': {
		"█████",
		"    █",
		"   █ ",
		"  █  ",
		"  █  ",
	},
	'8': {
		" ███ ",
		"█   █",
		" ███ ",
		"█   █",
		" ███ ",
	},
	'9': {
		" ███ ",
		"█   █",
		" ████",
		"    █",
		" ███ ",
	},
	':': {
		"     ",
		"  █  ",
		"     ",
		"  █  ",
		"     ",
	},
}

func print_timer(t string, s int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
		height = 24
	}

	timerHeight := 5

	topPadding := (height - timerHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	for i := 0; i < topPadding; i++ {
		fmt.Println()
	}

	for row := 0; row < timerHeight; row++ {
		line := ""

		for _, ch := range t {
			line += digits[ch][row] + " "
		}

		lineWidth := len([]rune(line))

		leftPadding := (width - lineWidth) / 2
		if leftPadding < 0 {
			leftPadding = 0
		}

		if s > 4 {
			colorSuccess = color.New(color.FgGreen)
		} else {
			colorSuccess = color.New(color.FgRed)
		}

		colorSuccess.Printf("%*s%s\n", leftPadding, "", line)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func TimerSet(ctx context.Context, cmd *cli.Command) error {

	var totalSeconds int

	if cmd.Bool("his") {
		sqldb.Show_history()
		return nil
	}

	colorSuccess = color.New(color.FgGreen)

	if !cmd.IsSet("sec") && !cmd.IsSet("min") && !cmd.IsSet("hr") && !cmd.IsSet("del") && !cmd.IsSet("delete_all") {
		colorSuccess.Println("  ⚠  No flag provided. Use --help to see all available flags.")
		return nil
	}

	if cmd.IsSet("del") {
		index := cmd.Int("del")
		sqldb.Delete_Record(index)
		return nil
	}

	if cmd.Bool("delete_all") {
		colorSuccess.Println("if you want to delete all record (y/n)")
		var y_n string
		fmt.Scanln(&y_n)
		if len(y_n) > 0 && y_n[0] == 'y' {
			sqldb.Delete_All_Record()
		} else {
			colorSuccess = color.New(color.FgRed)
			colorSuccess.Println("command is rejected")
		}
		return nil
	}

	if cmd.IsSet("sec") {
		totalSeconds += cmd.Int("sec")
	}

	if cmd.IsSet("min") {
		totalSeconds += cmd.Int("min") * 60
	}

	if cmd.IsSet("hr") {
		totalSeconds += cmd.Int("hr") * 3600
	}

	if totalSeconds <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	task := "No Name"

	if cmd.NArg() > 0 {
		task = strings.Join(cmd.Args().Slice(), " ")
	}

	sqldb.Create_db()

	switch runtime.GOOS {
	case "windows":
		fmt.Printf("Timer started — %d seconds\n", totalSeconds)

		psCommand := fmt.Sprintf(
			"Start-Sleep -Seconds %d; Import-Module BurntToast; New-BurntToastNotification -Text '%d Seconds Timer Finished'",
			totalSeconds, totalSeconds,
		)

		cmdExec := exec.Command(
			"powershell.exe",
			"-NoProfile",
			"-Command",
			psCommand,
		)

		err := cmdExec.Start()

		if err != nil {
			return fmt.Errorf("failed to start timer: %v", err)
		}

		fmt.Println("Timer running in background")

	default:
		err_d := termbox.Init()

		if err_d != nil {
			panic(err_d)
		}

		defer termbox.Close()

		control := make(chan string)

		state := "resume"

		go func() {
			for {
				ev := termbox.PollEvent()
				if ev.Type == termbox.EventKey {
					if ev.Key == termbox.KeyCtrlC || ev.Ch == 'q' {
						control <- "stop"
						fmt.Println("timer stop")
						colorSuccess =color.New(color.FgGreen)
						colorSuccess.Println("insert timer data in db successfully \nCommand 'funk timer --his' ")
						return
					}
					if ev.Key == termbox.KeySpace {
						if state == "resume" {
							control <- "pause"
						} else {
							control <- "resume"
						}
					}
				}
			}
		}()

		st_timer := totalSeconds

		var h int

		var m int

		var s int

		h1, m1, s1 := timer_cal(totalSeconds)

		for st_timer >= 0 {

			h, m, s = timer_cal(st_timer)

			select {
			case con := <-control:
				if con == "stop" {
					save_timer(h, m, s, task)
					return nil
				}

				if con == "pause" {
					state = "pause"
				}
				if con == "resume" {
					state = "resume"
				}

			default:
				a := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
				if state == "pause" {
					clearScreen()
					fmt.Println("\nTimer Pause")
				} else if state == "resume" {
					clearScreen()
					st_timer--
				}
				print_timer(a, st_timer)
				time.Sleep(time.Second)
				if st_timer < 0 {
					save_timer(h1, m1, s1, task)
				}
			}
		}
		fmt.Println("\ntimer finish")
		colorSuccess =color.New(color.FgGreen)
		colorSuccess.Println("insert timer data in db successfully \nCommand 'funk timer --his' ")
	}
	return nil
}

func save_timer(h int, m int, s int, task string) {
	termbox.Close()
	sqldb.Insert_data(h, m, s, task)
}

func timer_cal(i int) (int, int, int) {
	h := i / 3600
	m := (i % 3600) / 60
	s := i % 60
	return h, m, s
}