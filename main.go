package main
import (
    "fmt"
    _ "io"
    "net/http"
    "log"
    "runtime"
    "encoding/json"
    "github.com/gorilla/mux"

)

 func indexPageHandler(response http.ResponseWriter, request *http.Request) {
     const indexPage = `
     <html>
     <body>
     Main Page
     </body>
     </html>
     `
     fmt.Fprintf(response, indexPage)
 }

 func PageHandler404(response http.ResponseWriter, request *http.Request) {
     const error404 = `
     <h1>4XX error so you've been 307'd to a Page expressing our appologies</h1>
     `
     fmt.Fprintf(response, error404)
 }

func memStat(response http.ResponseWriter, request *http.Request) {
  var mem runtime.MemStats
         cpu := runtime.NumCPU()
        runtime.ReadMemStats(&mem)
        //fmt.Fprintln(response, "CPUs:", cpu)
        //fmt.Fprintln(response, "Cpu used: ", mem.GCCPUFraction*100)
        //fmt.Fprintln(response, "Alloc: ", mem.Alloc)
        //fmt.Fprintln(response, "Total System memmory: ", mem.Sys)
        //  fmt.Fprintln(response, "memused: ", (mem.Sys/mem.Alloc),"%")
        //marshall the mem structure into a json document
        b, err := json.MarshalIndent(&mem, "", "    ")
        if err != nil {
            fmt.Printf("Error: %s", err)
            return;
        }
        fmt.Fprintln(response, string(b))
        runtime.GOMAXPROCS(cpu)
        log.Println("API TOUCHED")
}

/// create a router with the gorilla mux router and handle the requests
var router = mux.NewRouter()
func main() {
  //configure NotFoundHandler
  router.NotFoundHandler = http.HandlerFunc(PageHandler404)
  //set a main root router
  //set a GET method router to serve up some json
  router.HandleFunc("/memStat.json", memStat).Methods("GET")
  router.Handle("/", http.FileServer(http.Dir("./static/")))

  http.Handle("/", router)
    err := http.ListenAndServeTLS(":8000", "server.crt", "server.key", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
