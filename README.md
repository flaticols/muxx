# muxx

> [!IMPORTANT]
> Use https://github.com/go-pkgz/routegroup, this library is better and does the same thing.

`muxx` added suuport of route groups and middleswares to the standard `http.ServeMux`.

### Install

```bash
go get github.com/flaticols/muxx
```

### Usage

##### Create a new `muxx` instance

```go

func main() {
    mux := muxx.New()
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    http.ListenAndServe(":8080", mux)
}
```
##### or mount to an existing `http.ServeMux`

```go
func main() {
    mux := http.NewServeMux()
    admingGroup := muxx.Mount(mux, "/admin")
    admingGroup.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    http.ListenAndServe(":8080", mux)
}
```

#### middleswares

```go
func main() {
    mux := muxx.New()
    
    mux.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Request-Id", "123")
            next.ServeHTTP(w, r)
        })
    })

    mux.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Request-Actor", "123")
            next.ServeHTTP(w, r)
        })
    })

    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    http.ListenAndServe(":8080", mux)
}
```

##### Apply middleswares to a group

```go
func main() {
    mux := muxx.New()
    adminGroup := mux.Group("/admin")
    adminGroup.Route(func (g *muxx.Group) {
        g.Use(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("X-Request-Id", "123")
                next.ServeHTTP(w, r)
            })
        })

        g.Handle("GET /hey", heyCtrl)
    })

    adminGroup.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    http.ListenAndServe(":8080", mux)
}
```
