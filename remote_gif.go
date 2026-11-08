package main

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	remoteGifDownloadTimeout = 30 * time.Second
	remoteGifMaxSizeInBytes  = 50 << 20
)

const supportedURLFormats = `supported URLs:
  https://giphy.com/gifs/<name>-<id>
  https://i.giphy.com/<id>.gif
  https://media.giphy.com/media/<id>/giphy.gif`

var (
	giphyIDPattern        = regexp.MustCompile(`^[A-Za-z0-9]+$`)
	giphyMediaHostPattern = regexp.MustCompile(`^media[0-9]*\.giphy\.com$`)
)

func isRemoteGifSource(source string) bool {
	lowerSource := strings.ToLower(source)
	return strings.HasPrefix(lowerSource, "http://") || strings.HasPrefix(lowerSource, "https://")
}

// resolveGiphyGifURL turns any supported Giphy URL into a direct link to the
// .gif file, so page URLs never need to be scraped.
func resolveGiphyGifURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := strings.ToLower(parsedURL.Hostname())
	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	switch {
	case host == "giphy.com" || host == "www.giphy.com":
		if len(pathSegments) != 2 || pathSegments[0] != "gifs" {
			return "", fmt.Errorf("unsupported Giphy page %q\n%s", rawURL, supportedURLFormats)
		}
		slug := pathSegments[1]
		giphyID := slug[strings.LastIndex(slug, "-")+1:]
		if !giphyIDPattern.MatchString(giphyID) {
			return "", fmt.Errorf("could not find a Giphy ID in %q\n%s", rawURL, supportedURLFormats)
		}
		return "https://i.giphy.com/" + giphyID + ".gif", nil

	case host == "i.giphy.com" || giphyMediaHostPattern.MatchString(host):
		if extension := path.Ext(parsedURL.Path); extension != ".gif" {
			return "", fmt.Errorf("%q is not a .gif file (got %q), only GIFs are supported\n%s", rawURL, extension, supportedURLFormats)
		}
		parsedURL.Scheme = "https"
		return parsedURL.String(), nil

	default:
		return "", fmt.Errorf("unsupported website %q, only Giphy URLs are supported\n%s", host, supportedURLFormats)
	}
}

func downloadGif(gifURL string) (io.Reader, error) {
	httpClient := http.Client{Timeout: remoteGifDownloadTimeout}

	request, err := http.NewRequest(http.MethodGet, gifURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "gifterm/"+version)

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not download %s: %w", gifURL, err)
	}
	defer response.Body.Close()

	// Giphy answers unknown IDs with 200 and a placeholder GIF; this undocumented
	// header is the only thing that marks it, so at worst the placeholder renders.
	if response.StatusCode == http.StatusNotFound || response.Header.Get("X-Retry-Not-Found-Metric") != "" {
		return nil, fmt.Errorf("GIF not found at %s, check that the URL is correct", gifURL)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloading %s failed with status %s", gifURL, response.Status)
	}

	mediaType, _, _ := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if mediaType != "image/gif" {
		return nil, fmt.Errorf("%s returned %q instead of a GIF", gifURL, mediaType)
	}

	// read one extra byte to tell a file of exactly the max size from a bigger one
	body, err := io.ReadAll(io.LimitReader(response.Body, remoteGifMaxSizeInBytes+1))
	if err != nil {
		return nil, fmt.Errorf("could not download %s: %w", gifURL, err)
	}
	if len(body) > remoteGifMaxSizeInBytes {
		return nil, fmt.Errorf("GIF at %s is larger than %d MB", gifURL, remoteGifMaxSizeInBytes>>20)
	}

	return bytes.NewReader(body), nil
}
