//go:build file || all
package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

// COLORS
var (
	colorCyan   = color.New(color.FgCyan, color.Bold)
	colorYellow = color.New(color.FgYellow)
)

func init() {
	AvailableCommands = append(AvailableCommands, FileDetectCommand())
}



// TERMINAL WIDTH
func termWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < 40 {
		return 60
	}
	if width > 80 {
		return 80
	}
	return width
}

// SECTION HEADER
func printSection(icon, title string) {
	w := termWidth()
	dashes := w - len(icon) - len(title) - 6
	if dashes < 2 {
		dashes = 2
	}
	fmt.Println()
	colorCyan.Printf("┌─ %s %s %s\n\n", icon, title, strings.Repeat("─", dashes))
}

// DIVIDER
func printDivider() {
	colorCyan.Println("└" + strings.Repeat("─", termWidth()-1))
}

// NO RESULT
func noResult(msg string) {
	colorYellow.Println("  ⚠  " + msg)
}

// TABLE HELPER
func newTable(headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(true)
	table.SetCenterSeparator("┼")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	return table
}

func FileDetectCommand() *cli.Command {
	return &cli.Command{
		Name:  "sift",
		Usage: "Detect and scan files with various filters",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "emt", Aliases: []string{"e"}, Usage: "Find empty files"},
			&cli.IntFlag{Name: "rec", Aliases: []string{"r"}, Usage: "Find files modified in last N days"},
			&cli.Float64Flag{Name: "lrg", Aliases: []string{"l"}, Usage: "Find files larger than N GB"},
			&cli.BoolFlag{Name: "ext", Aliases: []string{"x"}, Usage: "Show file extension summary"},
			&cli.IntFlag{Name: "top", Aliases: []string{"t"}, Usage: "Show top N largest files"},
			&cli.BoolFlag{Name: "dup", Aliases: []string{"d"}, Usage: "Find duplicate files by name"},
			&cli.BoolFlag{Name: "log", Aliases: []string{"L"}, Usage: "Find all log files (.log)"},
			&cli.StringFlag{Name: "p", Aliases: []string{"P"}, Value: ".", Usage: "Directory path to scan"},
		},
		Action: runDetect,
	}
}

func runDetect(ctx context.Context, c *cli.Command) error {

	path := c.String("p")

	// PATH VALIDATION
	if _, err := os.Stat(path); os.IsNotExist(err) {
		colorYellow.Println("  ⚠  Invalid path:", path)
		return nil
	}

	// no flag check
	if !c.Bool("emt") && !c.Bool("ext") && !c.Bool("dup") && !c.Bool("log") &&
		c.Int("rec") == 0 && c.Float64("lrg") == 0 && c.Int("top") == 0 {
		fmt.Println()
		colorYellow.Println("  ⚠  No flag provided. Use --help to see all available flags.")
		return nil
	}

	// EMPTY FILES
	if c.Bool("emt") {
		printSection("🔍", "Empty Files")
		table := newTable([]string{"#", "FILE"})
		count := 0

		filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if info.Size() == 0 {
				count++
				table.Append([]string{fmt.Sprintf("%d", count), p})
			}
			return nil
		})

		if count == 0 {
			noResult("No empty files found.")
		} else {
			table.Render()
		}
		printDivider()
	}

	// RECENT FILES
	if c.Int("rec") > 0 {
		days := c.Int("rec")
		cutoff := time.Now().AddDate(0, 0, -days)

		printSection("🕐", fmt.Sprintf("Recently Modified — Last %d Day(s)", days))
		table := newTable([]string{"#", "FILE", "MODIFIED"})
		count := 0

		filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if info.ModTime().After(cutoff) {
				count++
				table.Append([]string{
					fmt.Sprintf("%d", count),
					p,
					info.ModTime().Format("2006-01-02 15:04:05"),
				})
			}
			return nil
		})

		if count == 0 {
			noResult(fmt.Sprintf("No files modified in last %d day(s).", days))
		} else {
			table.Render()
		}
		printDivider()
	}

	// LARGE FILES
	if c.Float64("lrg") > 0 {
		sizeGB := c.Float64("lrg")
		minBytes := int64(sizeGB * 1024 * 1024 * 1024)

		printSection("💾", "Large Files")
		table := newTable([]string{"#", "FILE", "SIZE"})
		count := 0

		filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if info.Size() > minBytes {
				count++
				table.Append([]string{
					fmt.Sprintf("%d", count),
					p,
					fmt.Sprintf("%.2f GB", float64(info.Size())/(1024*1024*1024)),
				})
			}
			return nil
		})

		if count == 0 {
			noResult("No large files found.")
		} else {
			table.Render()
		}
		printDivider()
	}

	// TOP FILES
	if c.Int("top") > 0 {
		n := c.Int("top")

		type file struct {
			path string
			size int64
		}
		var files []file

		filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			files = append(files, file{p, info.Size()})
			return nil
		})

		sort.Slice(files, func(i, j int) bool {
			return files[i].size > files[j].size
		})

		if n > len(files) {
			n = len(files)
		}

		printSection("🏆", fmt.Sprintf("Top %d Files", n))
		table := newTable([]string{"RANK", "FILE", "SIZE"})
		for i := 0; i < n; i++ {
			table.Append([]string{
				fmt.Sprintf("#%d", i+1),
				files[i].path,
				fmt.Sprintf("%.2f MB", float64(files[i].size)/(1024*1024)),
			})
		}

		if n == 0 {
			noResult("No files found.")
		} else {
			table.Render()
		}
		printDivider()
	}

	// EXTENSIONS
	if c.Bool("ext") {
		printSection("📋", "Extensions")
		extMap := make(map[string]int)

		filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(p))
			if ext == "" {
				ext = "(none)"
			}
			extMap[ext]++
			return nil
		})

		table := newTable([]string{"EXT", "COUNT"})
		for k, v := range extMap {
			table.Append([]string{k, fmt.Sprintf("%d", v)})
		}

		if len(extMap) == 0 {
			noResult("No files found.")
		} else {
			table.Render()
		}
		printDivider()
	}

	// DUPLICATES
	if c.Bool("dup") {
		printSection("♊", "Duplicates")
		nameMap := make(map[string][]string)

		filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			nameMap[d.Name()] = append(nameMap[d.Name()], p)
			return nil
		})

		table := newTable([]string{"FILE", "PATHS"})
		count := 0
		for name, paths := range nameMap {
			if len(paths) > 1 {
				count++
				table.Append([]string{name, strings.Join(paths, " | ")})
			}
		}

		if count == 0 {
			noResult("No duplicates found.")
		} else {
			table.Render()
		}
		printDivider()
	}

	// LOG FILES
	if c.Bool("log") {
		printSection("📄", "Log Files")
		table := newTable([]string{"FILE"})
		count := 0

		filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if strings.HasSuffix(strings.ToLower(d.Name()), ".log") {
				count++
				table.Append([]string{p})
			}
			return nil
		})

		if count == 0 {
			noResult("No log files found.")
		} else {
			table.Render()
		}
		printDivider()
	}

	return nil
}
