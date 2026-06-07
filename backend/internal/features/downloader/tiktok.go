package downloader

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// tiktokDownloader uses musicaldown.com: GET page for form token + cookie, POST form, parse HTML.
type tiktokDownloader struct{ *base }

const (
	musicalDownURL = "https://musicaldown.com"
	musicalDownAPI = "https://musicaldown.com/download"
)

func (d *tiktokDownloader) ID() string           { return "musicaldown" }
func (d *tiktokDownloader) Platform() string     { return "tiktok" }
func (d *tiktokDownloader) PlatformName() string { return "TikTok" }

func (d *tiktokDownloader) CanHandle(u string) bool {
	u = strings.ToLower(u)
	return strings.Contains(u, "tiktok.com") || strings.Contains(u, "vm.tiktok.com") || strings.Contains(u, "vt.tiktok.com")
}

func (d *tiktokDownloader) Scrape(ctx context.Context, rawURL string) (*MediaResult, error) {
	html, hdr, err := d.httpGet(ctx, musicalDownURL, map[string]string{"User-Agent": chromeUA})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	sessionCookie := firstCookie(hdr)
	inputs := doc.Find("div > input")
	if inputs.Length() == 0 {
		return nil, ErrNoMedia
	}

	firstName, _ := inputs.First().Attr("name")
	form := url.Values{}
	inputs.Each(func(_ int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "" {
			return
		}
		if name == firstName {
			form.Set(name, rawURL)
		} else {
			v, _ := s.Attr("value")
			form.Set(name, v)
		}
	})

	headers := map[string]string{
		"User-Agent": chromeUA,
		"Origin":     "https://musicaldown.com",
		"Referer":    "https://musicaldown.com/en",
	}
	if sessionCookie != "" {
		headers["Cookie"] = sessionCookie
	}

	resultHTML, _, err := d.httpPostForm(ctx, musicalDownAPI, form, headers)
	if err != nil {
		return nil, err
	}
	resultDoc, err := goquery.NewDocumentFromReader(strings.NewReader(resultHTML))
	if err != nil {
		return nil, err
	}

	videos := map[string]string{}
	containers := resultDoc.Find("div.row > div")
	if containers.Length() > 1 {
		containers.Eq(1).Find("a").Each(func(_ int, a *goquery.Selection) {
			href, _ := a.Attr("href")
			if href == "" || href == "#modal2" {
				return
			}
			ev, _ := a.Attr("data-event")
			ev = strings.ToLower(ev)
			switch {
			case strings.Contains(ev, "hd"):
				videos["videoHD"] = href
			case strings.Contains(ev, "mp4"):
				videos["videoSD"] = href
			case strings.Contains(ev, "watermark"):
				videos["videoWatermark"] = href
			case strings.Contains(strings.ToLower(href), "type=mp3"):
				videos["music"] = href
			}
		})
	}

	var images []string
	resultDoc.Find("div.row > div.col.s12.m3").Each(func(_ int, s *goquery.Selection) {
		if src, ok := s.Find("img").First().Attr("src"); ok {
			if src = strings.TrimSpace(src); src != "" {
				images = append(images, src)
			}
		}
	})

	res := &MediaResult{
		Platform: d.Platform(), PlatformName: d.PlatformName(), Downloader: "MusicalDown",
	}

	if len(images) > 0 {
		res.Images = images
		if a := videos["music"]; a != "" {
			res.DownloadItems = append(res.DownloadItems, DownloadItem{
				Key: "audio", Label: "Audio", Type: "audio", URL: a, MimeType: "audio/mpeg",
			})
		}
		return res, nil
	}

	res.Thumbnail, _ = resultDoc.Find("div.img-area > img").First().Attr("src")
	res.AuthorName = strings.TrimSpace(resultDoc.Find("h2.video-author > b").First().Text())
	res.Title = strings.TrimSpace(resultDoc.Find("p.video-desc").First().Text())

	add := func(key, label, url, quality string) {
		if url == "" {
			return
		}
		res.DownloadItems = append(res.DownloadItems, DownloadItem{
			Key: key, Label: label, Type: "video", URL: url, MimeType: "video/mp4", Quality: quality,
		})
	}
	add("video", "Video", videos["videoSD"], "")
	add("video_hd", "Video HD", videos["videoHD"], "HD")
	add("video_watermark", "Video Watermark", videos["videoWatermark"], "Watermark")
	if a := videos["music"]; a != "" {
		res.DownloadItems = append(res.DownloadItems, DownloadItem{
			Key: "audio", Label: "Audio", Type: "audio", URL: a, MimeType: "audio/mpeg",
		})
	}

	if len(res.DownloadItems) == 0 {
		return nil, ErrNoMedia
	}
	return res, nil
}

func firstCookie(h http.Header) string {
	sc := h.Get("Set-Cookie")
	if sc == "" {
		return ""
	}
	return strings.TrimSpace(strings.SplitN(sc, ";", 2)[0])
}
