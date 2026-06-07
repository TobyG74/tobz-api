package search

import (
	"context"
	"fmt"
	"net/url"
)

// pixivClient searches Pixiv artworks/manga/novels
type pixivClient struct{ *base }

// PixivResult is one search hit from Pixiv.
type PixivResult struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Type     string `json:"type"` // artworks | slide | manga | novel
}

var pixivHeaders = map[string]string{"Referer": "https://www.pixiv.net/"}

type pixivItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	PageCount int    `json:"pageCount"`
}

func pixivSearchURL(kind, keyword string, page int) string {
	kw := url.PathEscape(keyword)
	word := url.QueryEscape(keyword)
	switch kind {
	case "manga":
		return fmt.Sprintf("https://www.pixiv.net/ajax/search/manga/%s?word=%s&order=date_d&mode=safe&p=%d&s_mode=s_tag&type=manga&work_lang=en&lang=en", kw, word, page)
	case "novel":
		return fmt.Sprintf("https://www.pixiv.net/ajax/search/novels/%s?word=%s&order=date_d&mode=all&p=%d&s_mode=s_tag&gs=0&lang=en", kw, word, page)
	default: // artworks
		return fmt.Sprintf("https://www.pixiv.net/ajax/search/artworks/%s?word=%s&order=date_d&mode=all&p=%d&s_mode=s_tag&type=all&lang=en", kw, word, page)
	}
}

// Search dispatches to the right Pixiv search by kind (artworks|manga|novel).
func (p *pixivClient) Search(ctx context.Context, kind, keyword string, page int) ([]PixivResult, error) {
	if page <= 0 {
		page = 1
	}
	var resp struct {
		Error bool `json:"error"`
		Body  struct {
			IllustManga struct{ Data []pixivItem } `json:"illustManga"`
			Manga       struct{ Data []pixivItem } `json:"manga"`
			Novel       struct{ Data []pixivItem } `json:"novel"`
		} `json:"body"`
	}
	if err := p.httpGetJSON(ctx, pixivSearchURL(kind, keyword, page), pixivHeaders, &resp); err != nil {
		return nil, err
	}

	var items []pixivItem
	switch kind {
	case "manga":
		items = resp.Body.Manga.Data
	case "novel":
		items = resp.Body.Novel.Data
	default:
		kind = "artworks"
		items = resp.Body.IllustManga.Data
	}
	if len(items) == 0 {
		return nil, ErrNoResults
	}

	out := make([]PixivResult, 0, len(items))
	for _, it := range items {
		t := kind
		if kind == "artworks" && it.PageCount > 1 {
			t = "slide"
		}
		out = append(out, PixivResult{
			ID: it.ID, Title: it.Title, UserID: it.UserID, UserName: it.UserName, Type: t,
		})
	}
	return out, nil
}
