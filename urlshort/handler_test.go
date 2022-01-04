package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func Test_parseYAML(t *testing.T) {
	testYaml := `
    - path: /urlshort
      url: https://github.com/gophercises/urlshort
    - path: /urlshort-final
      url: https://github.com/gophercises/urlshort/tree/solution
    `

	type args struct {
		yml []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []yamlFile
		wantErr bool
	}{
		{
			name: "Parse a sample yaml file",
			args: args{
				yml: []byte(testYaml),
			},
			want: []yamlFile{
				{Path: "/urlshort", Url: "https://github.com/gophercises/urlshort"},
				{Path: "/urlshort-final", Url: "https://github.com/gophercises/urlshort/tree/solution"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseYAML(tt.args.yml)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildMap(t *testing.T) {
	type args struct {
		parsedYaml []yamlFile
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Transform correctly a slice into a map",
			args: args{
				parsedYaml: []yamlFile{
					{Path: "/urlshort", Url: "https://github.com/gophercises/urlshort"},
					{Path: "/urlshort-final", Url: "https://github.com/gophercises/urlshort/tree/solution"},
				},
			},
			want: map[string]string{
				"/urlshort":       "https://github.com/gophercises/urlshort",
				"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildMap(tt.args.parsedYaml); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
