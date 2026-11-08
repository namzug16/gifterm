package main

import (
	"image"
	"image/color"
	"image/gif"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsRemoteGifSource(t *testing.T) {
	tests := []struct {
		source string
		want   bool
	}{
		{"https://giphy.com/gifs/abc", true},
		{"HTTP://giphy.com/gifs/abc", true},
		{"input.gif", false},
		{"./https.gif", false},
		{"ftp://example.com/a.gif", false},
	}

	for _, tt := range tests {
		if got := isRemoteGifSource(tt.source); got != tt.want {
			t.Errorf("isRemoteGifSource(%q) = %v, want %v", tt.source, got, tt.want)
		}
	}
}

func TestResolveGiphyGifURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		want    string
		wantErr string
	}{
		{name: "page with slug", rawURL: "https://giphy.com/gifs/anime-chainsaw-man-reze-IKFVtPf8jP6KJH16dB", want: "https://i.giphy.com/IKFVtPf8jP6KJH16dB.gif"},
		{name: "page without slug", rawURL: "https://giphy.com/gifs/IKFVtPf8jP6KJH16dB", want: "https://i.giphy.com/IKFVtPf8jP6KJH16dB.gif"},
		{name: "page with www, trailing slash and query", rawURL: "https://www.giphy.com/gifs/reze-IKFVtPf8jP6KJH16dB/?utm=x", want: "https://i.giphy.com/IKFVtPf8jP6KJH16dB.gif"},
		{name: "direct i.giphy.com", rawURL: "https://i.giphy.com/IKFVtPf8jP6KJH16dB.gif", want: "https://i.giphy.com/IKFVtPf8jP6KJH16dB.gif"},
		{name: "media url upgraded to https", rawURL: "http://media0.giphy.com/media/IKFVtPf8jP6KJH16dB/giphy.gif", want: "https://media0.giphy.com/media/IKFVtPf8jP6KJH16dB/giphy.gif"},
		{name: "unsupported website", rawURL: "https://tenor.com/view/cat-123", wantErr: `unsupported website "tenor.com"`},
		{name: "lookalike host", rawURL: "https://giphy.com.evil.io/gifs/a-b", wantErr: "unsupported website"},
		{name: "unsupported giphy page", rawURL: "https://giphy.com/explore/cats", wantErr: "unsupported Giphy page"},
		{name: "giphy page without id", rawURL: "https://giphy.com/gifs/reze-", wantErr: "could not find a Giphy ID"},
		{name: "webp instead of gif", rawURL: "https://media.giphy.com/media/IKFVtPf8jP6KJH16dB/giphy.webp", wantErr: "is not a .gif file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveGiphyGifURL(tt.rawURL)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDownloadGif(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok.gif":
			w.Header().Set("Content-Type", "image/gif")
			frame := image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{color.Black})
			gif.EncodeAll(w, &gif.GIF{Image: []*image.Paletted{frame}, Delay: []int{0}})
		case "/giphy-placeholder.gif":
			w.Header().Set("Content-Type", "image/gif")
			w.Header().Set("X-Retry-Not-Found-Metric", "1")
		case "/page.gif":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	gifReader, err := downloadGif(server.URL + "/ok.gif")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := gif.DecodeAll(gifReader); err != nil {
		t.Fatalf("downloaded GIF does not decode: %v", err)
	}

	if _, err := downloadGif(server.URL + "/missing.gif"); err == nil || !strings.Contains(err.Error(), "GIF not found") {
		t.Errorf("missing GIF error = %v", err)
	}

	if _, err := downloadGif(server.URL + "/giphy-placeholder.gif"); err == nil || !strings.Contains(err.Error(), "GIF not found") {
		t.Errorf("giphy placeholder error = %v", err)
	}

	if _, err := downloadGif(server.URL + "/page.gif"); err == nil || !strings.Contains(err.Error(), `returned "text/html" instead of a GIF`) {
		t.Errorf("wrong content type error = %v", err)
	}
}
