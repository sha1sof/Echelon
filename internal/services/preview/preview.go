package preview

import (
	"context"
	"errors"
	"fmt"
	"github.com/sha1sof/Echelon-/internal/storage"
	"github.com/sha1sof/Echelon-/internal/storage/sqlite"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type Preview struct {
	log      *slog.Logger
	urlSaver URLSaver
	url      Url
	lifetime time.Duration
}

// URLSaver for saving links and hash
type URLSaver interface {
	URLSaver(
		ctx context.Context,
		videoID string,
		hash []byte,
		lifetime time.Duration,
	) (uid int64, err error)
}

// Url for getting the hash
type Url interface {
	Url(ctx context.Context, videoID string) (cache []byte, err error)
}

// New конструктор.
func New(log *slog.Logger, urlSaver URLSaver, url Url, s *sqlite.Storage, lifeTime time.Duration) *Preview {
	return &Preview{
		log:      log,
		urlSaver: urlSaver,
		url:      url,
		lifetime: lifeTime,
	}
}

// GetThumbnail.
func (p *Preview) GetThumbnail(ctx context.Context, videoURL string) ([]byte, error) {
	const op = "Preview.GetThumbnail"

	log := p.log.With(
		slog.String("op", op),
		slog.String("url", videoURL))

	log.Info("Get video")

	videoID, err := extractVideoID(videoURL)
	if err != nil {
		log.Error("failed to extract video ID", slog.Any("error", err))
		return nil, err
	}

	cache, err := p.url.Url(ctx, videoID)
	if err == nil {
		return cache, nil
	}

	if errors.Is(err, storage.ErrVideoIDNotFound) {
		log.Warn("video not found", slog.String("videoID", videoID))

		cache, err = fetchThumbnail(videoID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		uid, err := p.urlSaver.URLSaver(ctx, videoID, cache, p.lifetime)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		log.Info("Saved", slog.Int64("uid", uid))
		return nil, nil
	}

	log.Warn("failed to get video URL", slog.Any("error", err))
	return nil, fmt.Errorf("%s: %w", op, err)
}

// extractVideoID поиск id видео.
func extractVideoID(videoURL string) (string, error) {
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", err
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")
	if videoID == "" {
		videoID = path.Base(parsedURL.Path)
	}

	videoID = strings.Split(videoID, "&")[0]

	if videoID == "" {
		return "", fmt.Errorf("video ID not found")
	}

	return videoID, nil
}

// fetchThumbnail скачивает видео.
func fetchThumbnail(videoID string) ([]byte, error) {
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID)
	req, err := http.NewRequest("GET", thumbnailURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching thumbnail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch thumbnail: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return data, nil
}
