package muxx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flaticols/muxx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Test-Middleware", "true")
		next.ServeHTTP(w, r)
	})
}

func TestGroup_Middleware(t *testing.T) {
	mux := http.NewServeMux()
	group, err := muxx.New("/test", muxx.WithMux(mux))
	require.NoError(t, err)

	// apply middleware to the group
	group.Use(testMiddleware)

	// add a test handler
	group.Handle("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make a request to the test handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test/test", http.NoBody)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "true", recorder.Header().Get("X-Test-Middleware"))
}
