# gorro
Stupid simple go router for my demo projects

## Usage example

```go
import (
  "github.com/mcustiel/gorro"
  "net/http"
)

func getHandler(w http.ResponseWriter, r *gorro.Request) error {
  w.Write([]byte("GET Hello World"))
}

func postHandler(w http.ResponseWriter, r *gorro.Request) error {
  w.Write([]byte("POST Hello World"))
}



func main() {
  rtr := gorro.NewRouter()

  err := rtr.Register(`/hello`, gorro.HandlersMap{
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
