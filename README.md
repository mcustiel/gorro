# gorro
Stupid simple go router for my demo projects

## Usage example

```go
import (
  "github.com/mcustiel/gorro"
  "net/http"
  "fmt"
)

func getHandler(w http.ResponseWriter, r *gorro.Request) error {
  fmt.Fprintf(w, "GET Hello %s", r.NamedParams["name"])
  return nil
}

func postHandler(w http.ResponseWriter, r *gorro.Request) error {
  fmt.Fprint(w, "POST Hello World")
  return nil
}



func main() {
  rtr := gorro.NewRouter()

  err := rtr.Register(`/hello/(?P<name>[a-z]+)`, gorro.HandlersMap{
    http.MethodPost: postHandler,
    http.MethodGet:  getHandler})

  if err != nil {
    panic(err)
  }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    rtr.Route(w, r)
  })

  http.ListenAndServe(":8080", nil)
}
```
