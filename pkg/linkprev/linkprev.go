package linkprev

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"github.com/nekomeowww/insights-bot/pkg/opengraph"
)

var (
	ErrNetworkError  = errors.New("network error")
	ErrRequestFailed = errors.New("request failed")
)

type Client struct {
	reqClient *req.Client
}

func NewClient() *Client {
	return &Client{
		reqClient: req.C().
			SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 Edg/111.0.1661.54").
			SetTimeout(time.Minute),
	}
}

func (c *Client) Debug() *Client {
	c.reqClient = c.reqClient.EnableDumpAll()
	return c
}

func (c *Client) Preview(ctx context.Context, url string) (Meta, error) {
	body, err := c.request(ctx, url)
	if err != nil {
		return Meta{}, err
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return Meta{}, fmt.Errorf("failed to parse response body with goquery: %v", err)
	}

	preview := newMetaFrom(doc)

	return preview, nil
}

func (c *Client) request(ctx context.Context, url string) (io.Reader, error) {
	resp, err := c.reqClient.
		R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get a preview of url %s, %w: %v", url, ErrNetworkError, err)
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get url %s, %w, status code: %d, dump: %s", url, ErrRequestFailed, resp.StatusCode, resp.Dump())
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return buf, nil
}

type Meta struct {
	Title       string
	Description string
	Favicon     string
	Author      string
	Keywords    []string
	OpenGraph   opengraph.OpenGraph
}

func newMetaFrom(doc *goquery.Document) Meta {
	meta := Meta{
		Title:       strings.TrimSpace(doc.Find("head > title").Text()),
		Description: strings.TrimSpace(doc.Find("head > meta[name='description']").AttrOr("content", "")),
		Favicon:     strings.TrimSpace(doc.Find("head > link[rel='icon']").AttrOr("href", "")),
		Author:      strings.TrimSpace(doc.Find("head > meta[name='author']").AttrOr("content", "")),
		Keywords: doc.Find("head > meta[name='keywords']").Map(func(i int, s *goquery.Selection) string {
			return strings.TrimSpace(s.AttrOr("content", ""))
		}),
		OpenGraph: opengraph.NewOpenGraphMetadataFromDocument(doc),
	}
	if meta.Title == "" && meta.OpenGraph.Title == "" {
		return Meta{}
	}

	return meta
}
