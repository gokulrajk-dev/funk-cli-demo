package commands

import (
	"codeberg.org/go-pdf/fpdf"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
)

func morphi(ctx context.Context, cmd *cli.Command) error {

	file := cmd.String("tp")
	if cmd.IsSet("tp") {
		text, err := os.ReadFile(file)
		pdf := fpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.MoveTo(0, 10)
		width, _ := pdf.GetPageSize()
		pdf.MultiCell(width, 10, string(text), "", "", false)
		err = pdf.OutputFileAndClose("hello.pdf")

		if err == nil {
			fmt.Println("PDF generated successfully")
		}
	}
	return nil

}

func Morph() *cli.Command {
	return &cli.Command{
		Name:   "morph",
		Usage:  "Converts files and images to various formats",
		Action: morphi,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "txt2pdf",
				Usage:   "Converts txt files to pdf",
				Aliases: []string{"tp"},
			},
		},
	}
}
