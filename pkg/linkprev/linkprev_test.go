package linkprev

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nekomeowww/insights-bot/pkg/opengraph"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreview(t *testing.T) {
	t.Run("GeneralWebsite", func(t *testing.T) {
		meta, err := NewClient().Preview(context.Background(), "https://nolebase.ayaka.io")
		assert.NoError(t, err)
		assert.Equal(t, Meta{
			Title:       "Nólëbase | 记录回忆，知识和畅想的地方",
			Description: "记录回忆，知识和畅想的地方",
			Favicon:     "/logo.svg",
			Author:      "Ayaka Neko, Ayaka Rizumu",
			Keywords: []string{
				"markdown, knowledgebase, 知识库, vitepress, obsidian, notebook, notes, nekomeowww, littlesound",
			},
			OpenGraph: opengraph.OpenGraph{
				Title:       "Nólëbase",
				Image:       "https://nolebase.ayaka.io/og.png",
				Description: "记录回忆，知识和畅想的地方",
				SiteName:    "Nólëbase",
			},
		}, meta)
	})

	t.Run("Twitter", func(t *testing.T) {
		t.Run("twitter.com", func(t *testing.T) {
			meta, err := NewClient().Preview(context.Background(), "https://twitter.com/GoogleDevEurope/status/1640667303158198272")
			require.NoError(t, err)
			assert.Equal(t, Meta{
				Title: "Google Developers Europe on Twitter: \"🎉 Happy Birthday @golang!\n\nDid you know that 11 years ago today Go 1 was publicly released? Join us in celebrating this day by:\n\n🎁 Checking out local meetups → https://t.co/TCNAZL0oOj\n🎁 Trying out the Go Playground → https://t.co/nnkaugz32x\n\nRT if you are a fellow Gopher! https://t.co/jiE7UTMHll\" / Twitter",
				OpenGraph: opengraph.OpenGraph{
					Title:       "Google Developers Europe on Twitter",
					Type:        "article",
					Image:       "https://pbs.twimg.com/media/FsTSN8nWwAA278D.png:large",
					URL:         "https://twitter.com/GoogleDevEurope/status/1640667303158198272",
					Description: "“🎉 Happy Birthday @golang!\n\nDid you know that 11 years ago today Go 1 was publicly released? Join us in celebrating this day by:\n\n🎁 Checking out local meetups → https://t.co/TCNAZL0oOj\n🎁 Trying out the Go Playground → https://t.co/nnkaugz32x\n\nRT if you are a fellow Gopher!”",
					SiteName:    "Twitter",
				},
			}, meta)
		})

		time.Sleep(time.Second)

		t.Run("fxtwitter.com", func(t *testing.T) {
			meta, err := NewClient().Preview(context.Background(), "https://fxtwitter.com/GoogleDevEurope/status/1640667303158198272")
			require.NoError(t, err)
			assert.Equal(t, Meta{
				OpenGraph: opengraph.OpenGraph{
					Title:       "Google Developers Europe (@GoogleDevEurope)",
					Image:       "https://pbs.twimg.com/media/FsTSN8nWwAA278D.png",
					Description: "🎉 Happy Birthday @golang!\n\nDid you know that 11 years ago today Go 1 was publicly released? Join us in celebrating this day by:\n\n🎁 Checking out local meetups → https://goo.gle/3zaGgRi\n🎁 Trying out the Go Playground → https://goo.gle/3zaGurC\n\nRT if you are a fellow Gopher!",
					SiteName:    "FixTweet",
				},
			}, meta)
		})

		t.Run("vxtwitter.com", func(t *testing.T) {
			meta, err := NewClient().Preview(context.Background(), "https://vxtwitter.com/GoogleDevEurope/status/1640667303158198272")
			require.NoError(t, err)
			assert.Equal(t, Meta{
				Title: "Google Developers Europe on Twitter: \"🎉 Happy Birthday @golang!\n\nDid you know that 11 years ago today Go 1 was publicly released? Join us in celebrating this day by:\n\n🎁 Checking out local meetups → https://t.co/TCNAZL0oOj\n🎁 Trying out the Go Playground → https://t.co/nnkaugz32x\n\nRT if you are a fellow Gopher! https://t.co/jiE7UTMHll\" / Twitter",
				OpenGraph: opengraph.OpenGraph{
					Title:       "Google Developers Europe on Twitter",
					Type:        "article",
					Image:       "https://pbs.twimg.com/media/FsTSN8nWwAA278D.png:large",
					URL:         "https://twitter.com/GoogleDevEurope/status/1640667303158198272",
					Description: "“🎉 Happy Birthday @golang!\n\nDid you know that 11 years ago today Go 1 was publicly released? Join us in celebrating this day by:\n\n🎁 Checking out local meetups → https://t.co/TCNAZL0oOj\n🎁 Trying out the Go Playground → https://t.co/nnkaugz32x\n\nRT if you are a fellow Gopher!”",
					SiteName:    "Twitter",
				},
			}, meta)
		})
	})
}

func TestNewMetaFrom(t *testing.T) {
	html := `<html>
  <head>
    <title>Example Movie</title>
	<meta name="description" content="Example description">
    <link rel="icon" href="/logo.svg" type="image/svg+xml">
	<meta property="og:title" content="Example Movie" />
    <meta property="og:type" content="video.movie" />
    <meta property="og:url" content="https://example.com/movie" />
    <meta property="og:image" content="https://example.com/movie/poster.png" />
	<meta property="og:audio" content="https://example.com/bond/theme.mp3" />
    <meta property="og:description" content="Example description" />
    <meta property="og:determiner" content="the" />
    <meta property="og:locale" content="en_US" />
    <meta property="og:locale:alternate" content="fr_FR" />
    <meta property="og:locale:alternate" content="es_ES" />
    <meta property="og:site_name" content="Movie" />
    <meta property="og:video" content="https://example.com/bond/trailer.swf" />
  </head>
</html>`

	meta := newMetaFrom(lo.Must(goquery.NewDocumentFromReader(strings.NewReader(html))))
	assert.Equal(t, Meta{
		Title:       "Example Movie",
		Description: "Example description",
		Favicon:     "/logo.svg",
		OpenGraph: opengraph.OpenGraph{
			Title:       "Example Movie",
			Type:        "video.movie",
			Image:       "https://example.com/movie/poster.png",
			URL:         "https://example.com/movie",
			Audio:       "https://example.com/bond/theme.mp3",
			Description: "Example description",
			Determiner:  "the",
			Locale:      "en_US",
			LocaleAlternate: []string{
				"fr_FR",
				"es_ES",
			},
			SiteName: "Movie",
			Video:    "https://example.com/bond/trailer.swf",
		},
	}, meta)
}
