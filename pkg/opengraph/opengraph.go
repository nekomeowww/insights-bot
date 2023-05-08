package opengraph

import (
	"github.com/PuerkitoBio/goquery"
)

type OpenGraph struct {
	// The title of your object as it should appear within the graph,
	// e.g., "The Rock". Property: og:title
	Title string
	// The type of your object, e.g., "video.movie". Depending on the
	// type you specify, other properties may also be required.
	// Property: og:type
	Type string
	// An image URL which should represent your object within the graph.
	// Property: og:image
	Image string
	// The canonical URL of your object that will be used as its permanent ID in the graph, e.g., "https://www.imdb.com/title/tt0117500/".
	// Property: og:url
	URL string

	// A URL to an audio file to accompany this object. Property: og:audio
	Audio string
	// A one to two sentence description of your object. Property: og:description
	Description string
	// The word that appears before this object's title in a sentence.
	// An enum of (a, an, the, "", auto). If auto is chosen, the consumer
	// of your data should chose between "a" or "an". Default is "" (blank).
	// Property: og:determiner
	Determiner string
	// The locale these tags are marked up in. Of the format language_TERRITORY. Default is en_US.
	// Property: og:locale
	Locale string
	// An array of other locales this page is available in.
	// Property: og:locale:alternate
	LocaleAlternate []string
	// If your object is part of a larger web site, the name which should be displayed for the overall site. e.g., "IMDb".
	// Property: og:site_name
	SiteName string
	// A URL to a video file that complements this object.
	// Property: og:video
	Video string
}

func NewOpenGraphMetadataFromDocument(doc *goquery.Document) OpenGraph {
	return OpenGraph{
		Title:       doc.Find("head > meta[property='og:title']").AttrOr("content", ""),
		Type:        doc.Find("head > meta[property='og:type']").AttrOr("content", ""),
		Image:       doc.Find("head > meta[property='og:image']").AttrOr("content", ""),
		URL:         doc.Find("head > meta[property='og:url']").AttrOr("content", ""),
		Audio:       doc.Find("head > meta[property='og:audio']").AttrOr("content", ""),
		Description: doc.Find("head > meta[property='og:description']").AttrOr("content", ""),
		Determiner:  doc.Find("head > meta[property='og:determiner']").AttrOr("content", ""),
		Locale:      doc.Find("head > meta[property='og:locale']").AttrOr("content", ""),
		LocaleAlternate: doc.Find("head > meta[property='og:locale:alternate']").Map(
			func(i int, s *goquery.Selection) string {
				return s.AttrOr("content", "")
			},
		),
		SiteName: doc.Find("head > meta[property='og:site_name']").AttrOr("content", ""),
		Video:    doc.Find("head > meta[property='og:video']").AttrOr("content", ""),
	}
}
