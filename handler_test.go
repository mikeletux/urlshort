package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMapHandler(t *testing.T) {
	testPathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/test":           "",
	}

	type args struct {
		fallback   http.Handler
		requestURL string
	}
	tests := []struct {
		name string
		args args
		want int // Check return codes
	}{
		{
			name: "Check if redirection is done with correct parameters",
			args: args{
				fallback:   sampleTestMux(),
				requestURL: "/urlshort-godoc",
			},
			want: http.StatusMovedPermanently,
		},
		{
			name: "Check if redirection is NOT done with a key that doesn't exist",
			args: args{
				fallback:   sampleTestMux(),
				requestURL: "/yipikaiyei", // This should use the fallback handler as it doesn't exist
			},
			want: http.StatusOK,
		},
		{
			name: "Check if redirection is NOT done with a key that does exist but has no redirect URL",
			args: args{
				fallback:   sampleTestMux(),
				requestURL: "/test", // This should use the fallback handler as the redirect URL is empty
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapHandler(testPathsToUrls, tt.args.fallback)
			w := httptest.NewRecorder()
			got(w, httptest.NewRequest(http.MethodGet, tt.args.requestURL, nil))

			// Check
			if w.Result().StatusCode != tt.want {
				t.Errorf("Error, Got %d wanted %d", w.Result().StatusCode, http.StatusMovedPermanently)
			}
		})
	}
}

func sampleTestMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", sampleTestFallback)
	return mux
}

func sampleTestFallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}
