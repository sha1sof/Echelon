package main

import (
	"context"
	"flag"
	"fmt"
	previewgrpc "github.com/sha1sof/Echelon-/internal/clients/preview/grpc"
	"github.com/sha1sof/Echelon-/internal/config"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	asyncFlag := flag.Bool("async", false, "async flag")
	flag.Parse()

	videoLinks := flag.Args()
	if len(videoLinks) == 0 {
		fmt.Println("No video links provided")
		return
	}

	cfg := config.MustLoadPath("./config/prod.yaml")
	outputDir := cfg.Clients.Preview.OutputDir

	previewClient, err := previewgrpc.New(
		context.Background(),
		cfg.Clients.Preview.Address,
		cfg.Clients.Preview.Timeout,
		cfg.Clients.Preview.RetriesCount,
	)
	if err != nil {
		panic(err)
	}

	if *asyncFlag {
		var wg sync.WaitGroup
		for _, videoLink := range videoLinks {
			wg.Add(1)
			go func(videoLink string) {
				defer wg.Done()
				err := downloadThumbnail(previewClient, videoLink, outputDir)
				if err != nil {
					fmt.Printf("download thumbnail error: %v\n", err)
				}
			}(videoLink)
		}
		wg.Wait()
	} else {
		for _, videoLink := range videoLinks {
			err := downloadThumbnail(previewClient, videoLink, outputDir)
			if err != nil {
				fmt.Printf("download thumbnail error: %v\n", err)
			}
		}
	}
}

// downloadThumbnail сохраняет превью видео.
func downloadThumbnail(previewClient *previewgrpc.Client, videoLink, outputDir string) error {
	if videoLink == "" {
		return fmt.Errorf("video link is empty")
	}

	thumbnailData, err := previewClient.GetThumbnail(context.Background(), videoLink)
	if err != nil {
		return fmt.Errorf("failed to get thumbnail: %w", err)
	}

	fileName := fmt.Sprintf("%s.jpg", generateName(videoLink))
	filePath := filepath.Join(outputDir, fileName)

	err = ioutil.WriteFile(filePath, thumbnailData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	if len(thumbnailData) != 0 {
		fmt.Printf("Save: %s | Size: %d\n", filePath, len(thumbnailData))
	} else {
		fmt.Printf("Thumbnail data is empty. Please check the video link again: %s\n", videoLink)
	}

	return nil
}

// generateName извлекает id из ссылки на видео youtube.
func generateName(videoLink string) string {
	if strings.Contains(videoLink, "youtu.be/") {
		parts := strings.Split(videoLink, "youtu.be/")
		if len(parts) == 2 {
			return parts[1]
		}
	}

	parsedURL, err := url.Parse(videoLink)
	if err != nil {
		return "unknown"
	}

	videoID := parsedURL.Query().Get("v")
	if videoID == "" {
		return "unknown"
	}

	return videoID
}
