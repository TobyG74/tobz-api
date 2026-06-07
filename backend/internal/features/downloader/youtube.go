package downloader

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// youtubeDownloader uses the ytdown.to proxy. ytdown returns a per-format poll
// URL which we pass through as the item URL for the client to resolve.
type youtubeDownloader struct{ *base }

// Proxy URL stored as base64.
const ytdownProxyB64 = "aHR0cHM6Ly9hcHAueXRkb3duLnRvL3Byb3h5LnBocA=="

func (d *youtubeDownloader) ID() string           { return "ytdown" }
func (d *youtubeDownloader) Platform() string     { return "youtube" }
func (d *youtubeDownloader) PlatformName() string { return "YouTube" }

func (d *youtubeDownloader) CanHandle(u string) bool {
	u = strings.ToLower(u)
	return strings.Contains(u, "youtube.com") || strings.Contains(u, "youtu.be")
}

var reYTVideoID = regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/|youtube\.com/shorts/|[?&]v=)([a-zA-Z0-9_-]{11})`)

func (d *youtubeDownloader) Scrape(ctx context.Context, rawURL string) (*MediaResult, error) {
	videoID := ""
	if m := reYTVideoID.FindStringSubmatch(rawURL); m != nil {
		videoID = m[1]
	}

	proxyBytes, err := base64.StdEncoding.DecodeString(ytdownProxyB64)
	if err != nil {
		return nil, err
	}
	proxyURL := string(proxyBytes)

	form := url.Values{"url": {rawURL}}
	headers := map[string]string{
		"Accept":           "application/json, text/plain, */*",
		"Origin":           "https://app.ytdown.to",
		"Referer":          "https://app.ytdown.to/en23/",
		"x-requested-with": "XMLHttpRequest",
	}
	body, _, err := d.httpPostForm(ctx, proxyURL, form, headers)
	if err != nil {
		return nil, err
	}
	body = strings.TrimPrefix(strings.TrimSpace(body), "\ufeff") // strip BOM

	var parsed struct {
		API struct {
			Status          string `json:"status"`
			Message         string `json:"message"`
			Title           string `json:"title"`
			ImagePreviewURL string `json:"imagePreviewUrl"`
			UserInfo        struct {
				Name     string `json:"name"`
				Username string `json:"username"`
			} `json:"userInfo"`
			MediaItems []struct {
				Type           string `json:"type"`
				MediaURL       string `json:"mediaUrl"`
				MediaQuality   string `json:"mediaQuality"`
				MediaExtension string `json:"mediaExtension"`
				MediaFileSize  string `json:"mediaFileSize"`
				MediaRes       string `json:"mediaRes"`
			} `json:"mediaItems"`
		} `json:"api"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return nil, fmt.Errorf("respons ytdown tidak terbaca")
	}
	api := parsed.API
	if api.Status != "ok" {
		if api.Message != "" {
			return nil, fmt.Errorf("ytdown: %s", api.Message)
		}
		return nil, ErrNoMedia
	}

	res := &MediaResult{
		Platform: d.Platform(), PlatformName: d.PlatformName(), Downloader: "YTDown",
		Title:      orDefault(api.Title, "YouTube Video"),
		AuthorName: orDefault(api.UserInfo.Name, api.UserInfo.Username),
		Thumbnail:  api.ImagePreviewURL,
	}
	if res.Thumbnail == "" && videoID != "" {
		res.Thumbnail = "https://i.ytimg.com/vi/" + videoID + "/maxresdefault.jpg"
	}

	for i, it := range api.MediaItems {
		if it.MediaURL == "" {
			continue
		}
		isVideo := strings.EqualFold(it.Type, "Video")
		isAudio := strings.EqualFold(it.Type, "Audio")
		if !isVideo && !isAudio {
			continue
		}
		label := buildYTLabel(isVideo, it.MediaRes, it.MediaQuality, it.MediaExtension, it.MediaFileSize)
		typ := "audio"
		mime := audioMime(strings.ToLower(it.MediaExtension))
		if isVideo {
			typ = "video"
			mime = "video/mp4"
		}
		res.DownloadItems = append(res.DownloadItems, DownloadItem{
			Key:      fmt.Sprintf("%s_%s_%d", strings.ToLower(it.Type), strings.ToLower(it.MediaExtension), i+1),
			Label:    label,
			Type:     typ,
			URL:      it.MediaURL,
			MimeType: mime,
			Quality:  it.MediaQuality,
		})
	}
	if len(res.DownloadItems) == 0 {
		return nil, ErrNoMedia
	}
	return res, nil
}

func buildYTLabel(isVideo bool, res, quality, ext, size string) string {
	var b strings.Builder
	if isVideo {
		switch {
		case res != "":
			b.WriteString(res + " ")
		case quality != "":
			b.WriteString(quality + " ")
		}
		if ext != "" {
			b.WriteString(ext)
		}
	} else {
		b.WriteString("Audio")
		if quality != "" {
			b.WriteString(" " + quality)
		}
		if ext != "" {
			b.WriteString(" " + ext)
		}
	}
	if size != "" {
		b.WriteString(" (" + size + ")")
	}
	return strings.TrimSpace(b.String())
}

func audioMime(ext string) string {
	switch ext {
	case "mp3":
		return "audio/mpeg"
	case "m4a":
		return "audio/mp4"
	case "opus":
		return "audio/opus"
	default:
		if ext == "" {
			return "audio/mpeg"
		}
		return "audio/" + ext
	}
}

func orDefault(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}
