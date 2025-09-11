package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPokedexAPI_Success(t *testing.T) {
	want := `{"count":1,"next":null,"previous":null,"results":[{"name":"kanto","url":"https://example.com"}]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("received %s %s", r.Method, r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(want))
	}))
	defer ts.Close()

	got, err := GetPokedexAPI(ts.URL)
	if err != nil {
		t.Fatalf("GetPokedexAPI() unexpected error: %v", err)
	}

	if string(got) != want {
		t.Fatalf("GetPokedexAPI() got body = %q, want %q", string(got), want)
	}
}

func TestGetPokedexAPI_InvalidURL(t *testing.T) {
	// Invalid URL should cause http.Get to fail and propagate an error.
	if _, err := GetPokedexAPI(""); err == nil {
		t.Fatalf("GetPokedexAPI() expected error for invalid URL, got nil")
	}
}
