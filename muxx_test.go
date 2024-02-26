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
		w.Header().Add("X-Test", "true")
		next.ServeHTTP(w, r)
	})
}

func TestGroup_Mount(t *testing.T) {
	mux := http.NewServeMux()

	// add a test handler
	mux.HandleFunc("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	adminGroup := muxx.Mount(mux, "/admin")
	adminGroup.Handle("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make a request to the test handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test", http.NoBody)
	mux.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)

	t.Run("Mount admin routes", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/admin/test", http.NoBody)
		mux.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestGroup_Middleware(t *testing.T) {
	mux := http.NewServeMux()
	g := muxx.Mount(mux, "/test")
	g.Use(testMiddleware)

	g.Handle("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make a request to the test handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test/test", http.NoBody)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "true", recorder.Header().Get("X-Test"))
}

func TestGroup_Handle(t *testing.T) {
	mux, err := muxx.New()
	require.NoError(t, err)

	// add a test handler
	mux.Handle("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make a request to the test handler
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test", http.NoBody)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestMount(t *testing.T) {
	mux := http.NewServeMux()
	basePath := "/mounted"
	g := muxx.Mount(mux, basePath)

	g.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Mounted-Test-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	})

	g.Handle("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, basePath+"/test", http.NoBody)
	mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "true", recorder.Header().Get("X-Mounted-Test-Middleware"))
}
