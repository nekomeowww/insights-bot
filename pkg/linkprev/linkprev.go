package linkprev

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"github.com/nekomeowww/insights-bot/pkg/opengraph"
	"github.com/nekomeowww/xo"
	"github.com/samber/lo"
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
		reqClient: req.C().SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 Edg/111.0.1661.54"),
	}
}

func (c *Client) Preview(ctx context.Context, urlStr string) (Meta, error) {
	r := c.newRequest(ctx, urlStr)

	body, err := c.request(r, urlStr)
	if err != nil {
		return Meta{}, err
	}

	defer body.Reset()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return Meta{}, fmt.Errorf("failed to parse response body with goquery: %v", err)
	}

	preview := newMetaFrom(doc)

	return preview, nil
}

func (c *Client) newRequest(ctx context.Context, urlStr string) *req.Request {
	request := c.reqClient.
		R().
		SetContext(ctx)

	c.alterRequestForTwitter(request, urlStr)

	return request
}

// requestForTwitter is a special request for Twitter.
//
// We need to ask Twitter server to generate a SSR rendered HTML for us to get the metadata
// Learn more at:
//  1. https://stackoverflow.com/a/64332370/19954520
//  2. https://stackoverflow.com/a/64164115/19954520
//
// Other alternative User-Agent for Twitter:
//  1. Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)
//  2. Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)
//  3. facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)
//  4. Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)
func (c *Client) alterRequestForTwitter(request *req.Request, urlStr string) *req.Request {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return request
	}
	if !lo.Contains([]string{"twitter.com", "vxtwitter.com", "fxtwitter.com"}, parsedURL.Host) {
		return request
	}

	return request.SetHeader("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
}

func (c *Client) request(r *req.Request, urlStr string) (*bytes.Buffer, error) {
	dumpBuffer := new(bytes.Buffer)
	defer func() {
		dumpBuffer.Reset()
		dumpBuffer = nil
	}()

	request := r.
		EnableDumpTo(dumpBuffer).
		DisableAutoReadResponse()
	defer func() {
		request.EnableDumpTo(xo.NewNopIoWriter())
	}()

	resp, err := request.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get a preview of url %s, %w: %v", urlStr, ErrNetworkError, err)
	}
	if !resp.IsSuccessState() {
		errorBuf := new(bytes.Buffer)
		defer errorBuf.Reset()

		_, err = io.Copy(errorBuf, resp.Body)
		if err != nil {
			fmt.Fprintf(errorBuf, "failed to read response body: %v", err)
		}

		dumpBuffer.WriteString("\n")
		dumpBuffer.Write(errorBuf.Bytes())

		return nil, fmt.Errorf("failed to get url %s, %w, status code: %d, dump:\n%s", urlStr, ErrRequestFailed, resp.StatusCode, dumpBuffer.String())
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
