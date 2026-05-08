package get

import (
	"errors"
	"net/http/httptest"
	"testing"
	"url-shotener/internal/lib/api"
	"url-shotener/internal/lib/logger/emptylog"
	"url-shotener/internal/server/handlers/redirect/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "success",
			alias: "test_alias",
			url:   "https://hh.ru",
		},
		{
			name:  "Empty URL",
			alias: "something",
			url:   "",
		},
		{
			name:      "ProvideURL Error",
			alias:     "test",
			url:       "http://google.com",
			respError: "failed to get url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlProviderMock := mocks.NewURLProvider(t)

			if tc.respError == "" || tc.mockError != nil {
				urlProviderMock.On("ProvideURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", New(emptylog.NewEmptyLogger(), urlProviderMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			assert.Equal(t, tc.url, redirectedURL)
		})
	}
}
