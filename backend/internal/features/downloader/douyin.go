package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// douyinDownloader fetches the Douyin page and extracts the embedded JSON state.
type douyinDownloader struct{ *base }

const douyinMobileUA = "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36"
const douyinVideoTmpl = "https://www.douyin.com/aweme/v1/play/?video_id=%s&ratio=1080p&line=0"

func (d *douyinDownloader) ID() string           { return "douyin" }
func (d *douyinDownloader) Platform() string     { return "douyin" }
func (d *douyinDownloader) PlatformName() string { return "Douyin" }

func (d *douyinDownloader) CanHandle(u string) bool {
	u = strings.ToLower(u)
	return strings.Contains(u, "douyin.com") || strings.Contains(u, "iesdouyin.com")
}

var reFirstURL = regexp.MustCompile(`https?://[^\s"'<>]+`)

func (d *douyinDownloader) Scrape(ctx context.Context, rawURL string) (*MediaResult, error) {
	normalized := normalizeDouyinURL(rawURL)
	html, _, err := d.httpGet(ctx, normalized, map[string]string{"User-Agent": douyinMobileUA})
	if err != nil {
		return nil, err
	}

	item, err := douyinFindItem(html)
	if err != nil {
		return nil, err
	}

	res := &MediaResult{
		Platform: d.Platform(), PlatformName: d.PlatformName(), Downloader: "Douyin",
		Title: jsonString(item, "desc"),
	}
	if author := jsonObject(item, "author"); author != nil {
		res.AuthorName = orDefault(jsonString(author, "nickname"), jsonString(author, "unique_id"))
	}
	res.Duration = int(jsonNumber(item, "duration") / 1000)

	video := jsonObject(item, "video")
	if video != nil {
		if cover := jsonObject(video, "cover"); cover != nil {
			res.Thumbnail = firstStringInArray(cover, "url_list")
		}
		if playAddr := jsonObject(video, "play_addr"); playAddr != nil {
			uri := jsonString(playAddr, "uri")
			if uri != "" && !strings.HasSuffix(strings.ToLower(uri), "mp3") {
				res.DownloadItems = append(res.DownloadItems, DownloadItem{
					Key: "video", Label: "Video 1080p", Type: "video",
					URL: fmt.Sprintf(douyinVideoTmpl, uri), MimeType: "video/mp4", Quality: "1080p",
				})
			}
		}
	}

	if imgs, ok := item["images"].([]interface{}); ok {
		for _, e := range imgs {
			if m, ok := e.(map[string]interface{}); ok {
				if s := firstStringInArray(m, "url_list"); s != "" {
					res.Images = append(res.Images, s)
				}
			}
		}
		if res.Thumbnail == "" && len(res.Images) > 0 {
			res.Thumbnail = res.Images[0]
		}
	}

	if len(res.DownloadItems) == 0 && len(res.Images) == 0 {
		return nil, ErrNoMedia
	}
	return res, nil
}

func normalizeDouyinURL(raw string) string {
	t := strings.TrimSpace(raw)
	if m := reFirstURL.FindString(t); m != "" {
		t = m
	}
	low := strings.ToLower(t)
	if !strings.HasPrefix(low, "http://") && !strings.HasPrefix(low, "https://") {
		t = "https://" + t
	}
	return strings.TrimRight(t, "),.;!?")
}

// douyinFindItem finds the JSON node at loaderData -> page -> videoInfoRes.item_list[0].
func douyinFindItem(html string) (map[string]interface{}, error) {
	for _, obj := range extractJSONObjects(html) {
		loader := jsonObject(obj, "loaderData")
		if loader == nil {
			continue
		}
		page := jsonObject(loader, "video_(id)/page")
		if page == nil {
			page = jsonObject(loader, "note_(id)/page")
		}
		if page == nil {
			continue
		}
		info := jsonObject(page, "videoInfoRes")
		if info == nil {
			continue
		}
		if list, ok := info["item_list"].([]interface{}); ok && len(list) > 0 {
			if first, ok := list[0].(map[string]interface{}); ok {
				return first, nil
			}
		}
	}
	return nil, fmt.Errorf("info video Douyin tidak ditemukan")
}

// extractJSONObjects returns each valid top-level JSON object, using string-aware
// brace-matching so braces inside string literals don't break parsing.
func extractJSONObjects(text string) []map[string]interface{} {
	var results []map[string]interface{}
	start, depth := -1, 0
	inString, escape := false, false

	for i := 0; i < len(text); i++ {
		ch := text[i]
		if start == -1 {
			if ch == '{' {
				start, depth = i, 1
			}
			continue
		}
		if inString {
			switch {
			case escape:
				escape = false
			case ch == '\\':
				escape = true
			case ch == '"':
				inString = false
			}
			continue
		}
		switch ch {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				candidate := text[start : i+1]
				var m map[string]interface{}
				if json.Unmarshal([]byte(candidate), &m) == nil {
					results = append(results, m)
				}
				start = -1
			}
		}
	}
	return results
}

// JSON map navigation helpers.

func jsonObject(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return nil
}

func jsonString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func jsonNumber(m map[string]interface{}, key string) int64 {
	if v, ok := m[key].(float64); ok {
		return int64(v)
	}
	return 0
}

func firstStringInArray(m map[string]interface{}, key string) string {
	if arr, ok := m[key].([]interface{}); ok && len(arr) > 0 {
		if s, ok := arr[0].(string); ok {
			return s
		}
	}
	return ""
}
