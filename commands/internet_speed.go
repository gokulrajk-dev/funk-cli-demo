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
		Usage:   "measure internet speed",
		Aliases: []string{"is"},
		Action:  Net_Speed,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "start",
				Usage: "start the speed test",
			},
		},
	}
}

func Net_Speed(cxt context.Context, cmd *cli.Command) error {

	if !cmd.Bool("start") {
		fmt.Println("❌ Please use --start to begin the speed test")
		fmt.Println("👉 Example: funk internet_speed --start")
		return nil
	}

	fmt.Println("🚀 Starting Internet Speed Test...\n")

	download := testDownload()
	fmt.Printf("\n✅ Final Download: %.2f Mbps\n\n", download)

	upload := testUpload()
	fmt.Printf("\n✅ Final Upload: %.2f Mbps\n", upload)

	return nil
}

const (
	testDuration = 10 * time.Second
	workers      = 8
)

func downloadWorker(client *http.Client, url string, wg *sync.WaitGroup, bytesChan chan int64, stop <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			return
		default:
			resp, err := client.Get(url)
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

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	bytesChan := make(chan int64, 100)
	stop := make(chan struct{})

	start := time.Now()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go downloadWorker(client, url, &wg, bytesChan, stop)
	}

	var totalBytes int64
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timer := time.After(testDuration)

loop:
	for {
		select {
		case n := <-bytesChan:
			totalBytes += n

		case <-ticker.C:
			elapsed := time.Since(start).Seconds()
			speed := float64(totalBytes*8) / elapsed / 1_000_000
			fmt.Printf("\r📥 Downloading... %.2f Mbps", speed)

		case <-timer:
			break loop
		}
	}

	close(stop)
	wg.Wait()

	duration := time.Since(start).Seconds()
	return float64(totalBytes*8) / duration / 1_000_000
}



func uploadWorker(client *http.Client, wg *sync.WaitGroup, bytesChan chan int64, stop <-chan struct{}) {
	defer wg.Done()

	data := bytes.Repeat([]byte("a"), 1*1024*1024) 

	for {
		select {
		case <-stop:
			return
		default:
			resp, err := client.Post(
				"https://speed.cloudflare.com/__up",
				"application/octet-stream",
				bytes.NewReader(data),
			)

			if err != nil {
				continue
			}

			resp.Body.Close()
			bytesChan <- int64(len(data))
		}
	}
}

func testUpload() float64 {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	bytesChan := make(chan int64, 100)
	stop := make(chan struct{})

	start := time.Now()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go uploadWorker(client, &wg, bytesChan, stop)
	}

	var totalBytes int64
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timer := time.After(testDuration)

loop:
	for {
		select {
		case n := <-bytesChan:
			totalBytes += n

		case <-ticker.C:
			elapsed := time.Since(start).Seconds()
			speed := float64(totalBytes*8) / elapsed / 1_000_000
			fmt.Printf("\r📤 Uploading... %.2f Mbps", speed)

		case <-timer:
			break loop
		}
	}

	close(stop)
	wg.Wait()

	duration := time.Since(start).Seconds()
	return float64(totalBytes*8) / duration / 1_000_000
}