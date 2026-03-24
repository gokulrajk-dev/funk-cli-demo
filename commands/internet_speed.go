package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/urfave/cli/v3"
)

func Isps() *cli.Command {
	return &cli.Command{
		Name:    "internet_speed",
		Usage:   "use to measure the internet speed",
		Aliases: []string{"is"},
		Action:  Net_Speed,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "start",
				Usage: "use to start the testing fuctions",
			},
		},
	}
}

func Net_Speed(cxt context.Context, cmd *cli.Command) error {

	switch {
	case cmd.Bool("start"):
		download := testDownload()
		upload := testUpload()
		fmt.Printf("Download: %.2f Mbps\n", download)
		fmt.Printf("Upload: %.2f Mbps\n", upload)
	}
	return nil
}



const (
	testDuration = 10 * time.Second
	workers      = 8
)



func downloadWorker(url string, wg *sync.WaitGroup, bytesChan chan int64, stop <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			return
		default:
			resp, err := http.Get(url)
			if err != nil {
				continue
			}

			n, _ := io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			bytesChan <- n
		}
	}
}

func testDownload() float64 {
	url := "https://speed.cloudflare.com/__down?bytes=10000000"

	var wg sync.WaitGroup
	bytesChan := make(chan int64, workers)
	stop := make(chan struct{})

	start := time.Now()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go downloadWorker(url, &wg, bytesChan, stop)
	}

	var totalBytes int64

loop:
	for {
		select {
		case n := <-bytesChan:
			totalBytes += n
		default:
			if time.Since(start) > testDuration {
				break loop
			}
		}
	}

	close(stop)
	wg.Wait()

	duration := time.Since(start).Seconds()
	return float64(totalBytes*8) / duration / 1_000_000
}

func uploadWorker(wg *sync.WaitGroup, bytesChan chan int64, stop <-chan struct{}) {
	defer wg.Done()

	data := bytes.Repeat([]byte("a"), 1*1024*1024) // 1MB chunk

	for {
		select {
		case <-stop:
			return
		default:
			resp, err := http.Post("https://speed.cloudflare.com/__up",
				"application/octet-stream",
				bytes.NewReader(data))

			if err != nil {
				continue
			}

			resp.Body.Close()
			bytesChan <- int64(len(data))
		}
	}
}

func testUpload() float64 {
	var wg sync.WaitGroup
	bytesChan := make(chan int64, workers)
	stop := make(chan struct{})

	start := time.Now()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go uploadWorker(&wg, bytesChan, stop)
	}

	var totalBytes int64

loop:
	for {
		select {
		case n := <-bytesChan:
			totalBytes += n
		default:
			if time.Since(start) > testDuration {
				break loop
			}
		}
	}

	close(stop)
	wg.Wait()

	duration := time.Since(start).Seconds()
	return float64(totalBytes*8) / duration / 1_000_000
}