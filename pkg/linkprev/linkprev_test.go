package linkprev

import (
	"context"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/nekomeowww/insights-bot/pkg/opengraph"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestPreview(t *testing.T) {
	t.Run("GeneralWebsite", func(t *testing.T) {
		meta, err := NewClient().Preview(context.Background(), "https://nolebase.ayaka.io")
		assert.NoError(t, err)
		assert.Equal(t, Meta{
			Title:       "Nólëbase | 记录回忆，知识和畅想的地方",
			Description: "记录回忆，知识和畅想的地方",
			Favicon:     "/logo.svg",
			Author:      "絢香猫, 絢香音",
			Keywords: []string{
				"markdown, knowledge-base, 知识库, vitepress, obsidian, notebook, notes, nekomeowww, LittleSound",
			},
			OpenGraph: opengraph.OpenGraph{
				Title:           "Nólëbase",
				Image:           "https://nolebase.ayaka.io/og.png",
				Description:     "记录回忆，知识和畅想的地方",
				SiteName:        "Nólëbase",
				LocaleAlternate: make([]string, 0),
			},
		}, meta)
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
		Keywords:    make([]string, 0),
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
