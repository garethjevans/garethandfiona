package main

import (
  "github.com/magiconair/properties"
  "github.com/gorilla/mux"
  "github.com/newrelic/go-agent"
  "html/template"
  "net/http"
  "fmt"
  "log"
  "os"
)

func getProperties() map[string]string {

  var p map[string]string
  p = properties.MustLoadFile("static/uk.properties", properties.UTF8).Map()
  return p
}

func handler(w http.ResponseWriter, r *http.Request) {
  m := getProperties()
  if r.URL.Path == "/" {
    t, _ := template.ParseFiles("templates/index.tmpl")
    err := t.Execute(w, m)
    if err != nil {
      log.Fatal("Unable to parse template: ", err)
      os.Exit(2)
    }
  } else if r.URL.Path == "/ping" {
    fmt.Fprintf(w, "OK")
  } else {
    t, _ := template.ParseFiles(fmt.Sprintf("templates/%s.tmpl", r.URL.Path))
    err := t.Execute(w, m)
    if err != nil {
      log.Fatal("Unable to parse template: ", err)
      os.Exit(2)
    }
  }
}

func main() {
  newRelicLicenseKey := os.Getenv("NEWRELIC_LICENSE_KEY")
  newRelicApplicationName := os.Getenv("NEWRELIC_APP_NAME")

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  config := newrelic.NewConfig(newRelicApplicationName, newRelicLicenseKey)
  app, err := newrelic.NewApplication(config)

  if err != nil {
    log.Fatal("Unable to create new relic application: ", err)
    os.Exit(1)
  }

  r := mux.NewRouter()
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/", handler))
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/ping", handler))
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.ListenAndServe(fmt.Sprintf(":%s",port), r)
}
