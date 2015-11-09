package router

import (
	"net/http"
	"testing"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

type mockFileSystem struct {
	opened bool
}

type handlerStruct struct {
	handeled *bool
}

func TestRouter(t *testing.T) {
	var get, head, options, post, put, patch, delete bool

	router := New()
	router.Get("/GET", func(w http.ResponseWriter, r *http.Request) {
		get = true
	})
	router.Head("/GET", func(w http.ResponseWriter, r *http.Request) {
		head = true
	})
	router.Options("/GET", func(w http.ResponseWriter, r *http.Request) {
		options = true
	})
	router.Post("/POST", func(w http.ResponseWriter, r *http.Request) {
		post = true
	})
	router.Put("/PUT", func(w http.ResponseWriter, r *http.Request) {
		put = true
	})
	router.Patch("/PATCH", func(w http.ResponseWriter, r *http.Request) {
		patch = true
	})
	router.Del("/DELETE", func(w http.ResponseWriter, r *http.Request) {
		delete = true
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("HEAD", "/GET", nil)
	router.ServeHTTP(w, r)
	if !head {
		t.Error("routing HEAD failed")
	}

	r, _ = http.NewRequest("OPTIONS", "/GET", nil)
	router.ServeHTTP(w, r)
	if !options {
		t.Error("routing OPTIONS failed")
	}

	r, _ = http.NewRequest("POST", "/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest("PUT", "/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest("PATCH", "/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest("DELETE", "/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}
}

func TestRouterWithPrefix(t *testing.T) {
	var get, head, options, post, put, patch, delete bool

	router := New()
	router = router.WithPrefix("/api")
	router.Get("/GET", func(w http.ResponseWriter, r *http.Request) {
		get = true
	})
	router.Head("/GET", func(w http.ResponseWriter, r *http.Request) {
		head = true
	})
	router.Options("/GET", func(w http.ResponseWriter, r *http.Request) {
		options = true
	})
	router.Post("/POST", func(w http.ResponseWriter, r *http.Request) {
		post = true
	})
	router.Put("/PUT", func(w http.ResponseWriter, r *http.Request) {
		put = true
	})
	router.Patch("/PATCH", func(w http.ResponseWriter, r *http.Request) {
		patch = true
	})
	router.Del("/DELETE", func(w http.ResponseWriter, r *http.Request) {
		delete = true
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/api/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("HEAD", "/api/GET", nil)
	router.ServeHTTP(w, r)
	if !head {
		t.Error("routing HEAD failed")
	}

	r, _ = http.NewRequest("OPTIONS", "/api/GET", nil)
	router.ServeHTTP(w, r)
	if !options {
		t.Error("routing OPTIONS failed")
	}

	r, _ = http.NewRequest("POST", "/api/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest("PUT", "/api/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest("PATCH", "/api/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest("DELETE", "/api/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}
}

func TestRouterParam(t *testing.T) {
	var param bool

	router := New()
	router.Get("/user/:username", func(w http.ResponseWriter, r *http.Request) {
		param = true

		ctx := Context(r)
		width := Param(ctx, "username")
		if width != "foo" {
			t.Fatal("Got no param!")
		}
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/user/foo", nil)
	router.ServeHTTP(w, r)
	if !param {
		t.Error(param)
	}
}
