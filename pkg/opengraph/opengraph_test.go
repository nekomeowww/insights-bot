package opengraph

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	html := `<html>
  <head>
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
	og := NewOpenGraphMetadataFromDocument(lo.Must(goquery.NewDocumentFromReader(strings.NewReader(html))))
	assert.Equal(t, OpenGraph{
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
	}, og)
}
