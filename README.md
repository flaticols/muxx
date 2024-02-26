# muxx
Missing part in Go http.ServerMux


> [!TIP]
> TODO: docs & examples

### Install

```bash
go get github.com/flaticols/muxx
```

### How to use

```go
package main

import (
    "net/http"
    "github.com/flaticols/muxx"
)

func main() {
    mux := muxx.New()
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    http.ListenAndServe(":8080", mux)
}
```
