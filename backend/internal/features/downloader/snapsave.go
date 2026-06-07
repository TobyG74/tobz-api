package downloader

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// snapsaveDownloader handles snapsave.app platforms (Instagram, Facebook, Twitter/X).
// snapsave returns obfuscated JS that we decrypt into an HTML snippet, then parse.
type snapsaveDownloader struct {
	*base
	platform     string
	platformName string
	hosts        []string
	endpoint     snapsaveEndpoint
}

// snapsaveEndpoint describes the API call for one platform.
type snapsaveEndpoint struct {
	apiURL  string
	origin  string
	referer string
}

var (
	snapsaveIGEndpoint = snapsaveEndpoint{"https://snapsave.app/id/action.php?lang=id", "https://snapsave.app", "https://snapsave.app/id"}
	snapsaveFBEndpoint = snapsaveEndpoint{"https://snapsave.app/action.php?lang=en", "https://snapsave.app", "https://snapsave.app/id/facebook-reels-download"}
	snapsaveTWEndpoint = snapsaveEndpoint{"https://snapsave.app/action.php", "https://snapsave.app", "https://snapsave.app/"}
)

const snapsaveAPIBase = "https://snapsave.app"

func newSnapsave(b *base, platform, name string, hosts []string, ep snapsaveEndpoint) *snapsaveDownloader {
	return &snapsaveDownloader{base: b, platform: platform, platformName: name, hosts: hosts, endpoint: ep}
}

func (d *snapsaveDownloader) ID() string           { return "snapsave_" + d.platform }
func (d *snapsaveDownloader) Platform() string     { return d.platform }
func (d *snapsaveDownloader) PlatformName() string { return d.platformName }

func (d *snapsaveDownloader) CanHandle(u string) bool {
	u = strings.ToLower(u)
	for _, h := range d.hosts {
		if strings.Contains(u, h) {
			return true
		}
	}
	return false
}

func (d *snapsaveDownloader) Scrape(ctx context.Context, rawURL string) (*MediaResult, error) {
	form := url.Values{"url": {rawURL}}
	headers := map[string]string{
		"Accept":           "*/*",
		"Origin":           d.endpoint.origin,
		"Referer":          d.endpoint.referer,
		"X-Requested-With": "XMLHttpRequest",
	}
	encrypted, _, err := d.httpPostForm(ctx, d.endpoint.apiURL, form, headers)
	if err != nil {
		return nil, err
	}
	html, err := decryptSnapSave(encrypted)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(html) == "" {
		return nil, ErrNoMedia
	}
	return d.parse(html)
}

func (d *snapsaveDownloader) parse(htmlContent string) (*MediaResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var videoURL, thumbnail string
	var images []string

	switch {
	case doc.Find("table.table").Length() > 0:
		videoURL, thumbnail, err = parseSnapTable(doc)
	case doc.Find("div.download-items").Length() > 0:
		videoURL, images, err = parseSnapDownloadItems(doc)
	case doc.Find("div.card").Length() > 0:
		videoURL, images = parseSnapCard(doc)
	default:
		videoURL, images = parseSnapSingle(doc)
	}
	if err != nil {
		return nil, err
	}

	res := &MediaResult{
		Platform: d.platform, PlatformName: d.platformName, Downloader: "SnapSave",
		Thumbnail: thumbnail, Images: images,
	}
	if videoURL != "" {
		res.DownloadItems = append(res.DownloadItems, DownloadItem{
			Key: "video", Label: "Video", Type: "video", URL: videoURL, MimeType: "video/mp4", Quality: "Original",
		})
	}
	if len(res.DownloadItems) == 0 && len(res.Images) == 0 {
		return nil, ErrNoMedia
	}
	return res, nil
}

var reProgressAPI = regexp.MustCompile(`get_progressApi\('(.+?)'\)`)

func resolveProgress(onclick string) string {
	if strings.Contains(onclick, "get_progressApi") {
		if m := reProgressAPI.FindStringSubmatch(onclick); m != nil {
			return snapsaveAPIBase + m[1]
		}
	}
	return ""
}

func parseSnapTable(doc *goquery.Document) (string, string, error) {
	rows := doc.Find("table.table tbody tr")
	if rows.Length() == 0 {
		return "", "", ErrNoMedia
	}
	thumb, _ := doc.Find("article.media > figure img").First().Attr("src")
	cells := rows.First().Find("td")
	if cells.Length() < 3 {
		return "", "", ErrNoMedia
	}
	btn := cells.Eq(2).Find("button").First()
	onclick, _ := btn.Attr("onclick")
	videoURL := resolveProgress(onclick)
	if videoURL == "" {
		videoURL, _ = btn.Attr("href")
	}
	if videoURL == "" {
		videoURL, _ = cells.Eq(2).Find("a").First().Attr("href")
	}
	if videoURL == "" {
		return "", "", ErrNoMedia
	}
	return videoURL, thumb, nil
}

func parseSnapDownloadItems(doc *goquery.Document) (string, []string, error) {
	items := doc.Find("div.download-items")
	first := items.First()
	hasVideo := first.Find("video").Length() > 0
	span := strings.ToLower(strings.TrimSpace(first.Find("div.download-items__btn span").First().Text()))
	if hasVideo || strings.Contains(span, "video") {
		videoURL, _ := first.Find("div.download-items__btn a").First().Attr("href")
		if videoURL == "" {
			return "", nil, ErrNoMedia
		}
		return videoURL, nil, nil
	}
	var images []string
	items.Each(func(_ int, it *goquery.Selection) {
		src, ok := it.Find("div.download-items__thumb > img").First().Attr("src")
		if !ok || src == "" {
			src, _ = it.Find("div.download-items__btn a").First().Attr("href")
		}
		if src != "" && !strings.Contains(src, ".mp4") {
			images = append(images, src)
		}
	})
	if len(images) == 0 {
		return "", nil, ErrNoMedia
	}
	return "", images, nil
}

func parseSnapCard(doc *goquery.Document) (string, []string) {
	card := doc.Find("div.card").First()
	link := card.Find("div.card-body a").First()
	text := strings.ToLower(strings.TrimSpace(link.Text()))
	href, _ := link.Attr("href")
	if href == "" {
		return "", nil
	}
	if strings.Contains(text, "photo") {
		return "", []string{href}
	}
	return href, nil
}

func parseSnapSingle(doc *goquery.Document) (string, []string) {
	link := doc.Find("a").First()
	text := strings.ToLower(strings.TrimSpace(link.Text()))
	href, _ := link.Attr("href")
	if href == "" {
		onclick, _ := doc.Find("button").First().Attr("onclick")
		href = resolveProgress(onclick)
	}
	if href == "" {
		return "", nil
	}
	if strings.Contains(text, "photo") {
		return "", []string{href}
	}
	return href, nil
}

// ---- snapapp packer decryption ----

const snapAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/"

func decryptSnapSave(encrypted string) (string, error) {
	params := snapEncodedParams(encrypted)
	if len(params) == 0 {
		return "", ErrNoMedia
	}
	decoded := snapDecodeApp(params)
	if decoded == "" {
		return "", ErrNoMedia
	}
	return snapExtractHTML(decoded)
}

// snapEncodedParams extracts the packer call arguments from decodeURIComponent(escape(r))}(ARGS)).
func snapEncodedParams(data string) []string {
	parts := strings.SplitN(data, "decodeURIComponent(escape(r))}(", 2)
	if len(parts) < 2 {
		return nil
	}
	paramStr := strings.SplitN(parts[1], "))", 2)[0]
	raw := strings.Split(paramStr, ",")
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		out = append(out, strings.TrimSpace(strings.ReplaceAll(p, "\"", "")))
	}
	return out
}

func snapDecodeApp(args []string) string {
	if len(args) < 6 {
		return ""
	}
	t := args[0]
	o := args[2]
	b, _ := strconv.Atoi(args[3])
	z, _ := strconv.Atoi(args[4])
	if z <= 0 || z > len(o) {
		return ""
	}
	baseChars := []rune(snapAlphabet[:z])
	delim := rune(o[z])
	oRunes := []rune(o)

	var sb strings.Builder
	tRunes := []rune(t)
	i := 0
	for i < len(tRunes) {
		var s strings.Builder
		for i < len(tRunes) && tRunes[i] != delim {
			s.WriteRune(tRunes[i])
			i++
		}
		i++

		piece := s.String()
		for j := 0; j < len(oRunes); j++ {
			piece = strings.ReplaceAll(piece, string(oRunes[j]), strconv.Itoa(j))
		}
		charCode := snapDecodeBase(piece, z, baseChars) - b
		if charCode > 0 {
			sb.WriteRune(rune(charCode))
		}
	}
	return sb.String()
}

func snapDecodeBase(value string, base int, baseChars []rune) int {
	r := []rune(value)
	output := 0
	pow := 1
	for idx := len(r) - 1; idx >= 0; idx-- {
		digit := indexRune(baseChars, r[idx])
		if digit != -1 {
			output += digit * pow
		}
		pow *= base
	}
	return output
}

func indexRune(s []rune, target rune) int {
	for i, c := range s {
		if c == target {
			return i
		}
	}
	return -1
}

func snapExtractHTML(data string) (string, error) {
	// Surface any error message in #alert.
	if ep := strings.SplitN(data, `document.querySelector("#alert").innerHTML = "`, 2); len(ep) > 1 {
		msg := strings.TrimSpace(strings.SplitN(ep[1], `";`, 2)[0])
		if msg != "" {
			return "", errors.New("snapsave: " + msg)
		}
	}
	parts := strings.SplitN(data, `getElementById("download-section").innerHTML = "`, 2)
	if len(parts) < 2 {
		return "", ErrNoMedia
	}
	html := strings.SplitN(parts[1], `"; document.getElementById("inputData").remove(); `, 2)[0]
	return strings.ReplaceAll(html, "\\", ""), nil
}
