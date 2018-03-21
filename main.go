package main

import (
  "github.com/magiconair/properties"
  "github.com/gorilla/schema"
  "github.com/gorilla/mux"
  "github.com/newrelic/go-agent"
  "html/template"
  "net/http"
  "fmt"
  "log"
  "os"
  "garethandfiona/rsvp"
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
      log.Print("Unable to parse template: ", err)
      http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
    }
  } else if r.URL.Path == "/ping" {
    fmt.Fprintf(w, "OK")
  } else {
    t, _ := template.ParseFiles(fmt.Sprintf("templates/%s.tmpl", r.URL.Path))
    err := t.Execute(w, m)
    if err != nil {
      log.Print("Unable to parse template: ", err)
      http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
    }
  }
}

type Page struct {
  Rsvp      *rsvp.Rsvp
  P         map[string]string
}

func ShowRsvp(w http.ResponseWriter, r *http.Request) {
  showRsvpBase(w,r,"show_rsvp")
}

func ShowRsvp2(w http.ResponseWriter, r *http.Request) {
  showRsvpBase(w,r,"show_rsvp2")
}

func db() rsvp.WeddingDatabase {
  dbUser := os.Getenv("DB_USER")
  dbPassword := os.Getenv("DB_PASSWORD")
  dbHost := "127.0.0.1"
  dbPort := 3306

  db, err := rsvp.NewMySQLDB(rsvp.MySQLConfig{ Username:dbUser, Password:dbPassword, Host:dbHost, Port:dbPort })
  if err != nil {
    log.Fatal("Unable to connect to database: ", err)
    os.Exit(1)
  }

  return db
}

func showRsvpBase(w http.ResponseWriter, r *http.Request, v string) {
  params := mux.Vars(r)
  log.Printf("Show Rsvp %s", params["id"])

  db := db()
  log.Print("get db", db)
  item, err := db.GetRsvp(params["id"])
  defer db.Close()

  if err != nil {
    log.Print("Invalid reference: ", err)
    http.Error(w, err.Error(), http.StatusNotFound)
	return
  }

  t, _ := template.ParseFiles(fmt.Sprintf("templates/%s.tmpl", v))
  err2 := t.Execute(w, Page {Rsvp: item, P: getProperties()})
  if err2 != nil {
    log.Print("Unable to parse template: ", err2)
    http.Error(w, err.Error(), http.StatusInternalServerError)
	return
  }
}

func EditRsvp(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  log.Printf("Edit Rsvp %s", params["id"])

  db := db()
  item, err := db.GetRsvp(params["id"])
  defer db.Close()

  if err != nil {
    log.Print("Invalid reference: ", err)
    http.Error(w, err.Error(), http.StatusNotFound)
	return
  }

  t, _ := template.ParseFiles("templates/edit_rsvp.tmpl")
  err2 := t.Execute(w, Page {Rsvp: item, P: getProperties()})
  if err2 != nil {
    log.Print("Unable to parse template: ", err2)
    http.Error(w, err2.Error(), http.StatusInternalServerError)
	return
  }
}

func SaveRsvp(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  log.Printf("Save Rsvp %s", params["id"])

  db := db()
  item, err := db.GetRsvp(params["id"])
  defer db.Close()

  if err != nil {
    log.Print("Invalid reference: ", err)
    http.Error(w, err.Error(), http.StatusNotFound)
	return
  }

  if r.ParseForm() != nil {
    log.Print("Unable to parse form")
  }

  decoder := schema.NewDecoder()
  err2 := decoder.Decode(item, r.PostForm)

  if err2 != nil {
    log.Print("Unable to decode rsvp", err2)
    http.Error(w, err2.Error(), http.StatusInternalServerError)
	return
  }

  target := "http://" + r.Host + "/rsvp/" + item.RsvpID
  log.Print("Sending Redirect: " + target)
  http.Redirect(w, r, target, http.StatusSeeOther)
}

func SaveRsvp2(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  log.Printf("Save Rsvp %s", params["id"])

  db := db()
  item, err := db.GetRsvp(params["id"])
  defer db.Close()

  if err != nil {
    log.Print("Invalid reference: ", err)
    http.Error(w, err.Error(), http.StatusNotFound)
	return
  }

  if r.ParseForm() != nil {
    log.Print("Unable to parse form")
  }

  decoder := schema.NewDecoder()
  err2 := decoder.Decode(item, r.PostForm)

  if err2 != nil {
    log.Print("Unable to decode rsvp", err2)
    http.Error(w, err2.Error(), http.StatusInternalServerError)
	return
  }

  target := "http://" + r.Host + "/rsvp2/" + item.RsvpID
  log.Print("Sending Redirect: " + target)
  http.Redirect(w, r, target, http.StatusSeeOther)
}

func main() {
  newRelicLicenseKey := os.Getenv("NEWRELIC_LICENSE_KEY")
  newRelicApplicationName := os.Getenv("NEWRELIC_APP_NAME")

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  config := newrelic.NewConfig(newRelicApplicationName, newRelicLicenseKey)
  app, err2 := newrelic.NewApplication(config)

  if err2 != nil {
    log.Fatal("Unable to create new relic application: ", err2)
    os.Exit(1)
  }

  r := mux.NewRouter()
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/", handler))
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/ping", handler))

  // web calls
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/rsvp/{id}", ShowRsvp)).Methods("GET")
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/rsvp/{id}/save", SaveRsvp)).Methods("POST")

  r.HandleFunc(newrelic.WrapHandleFunc(app,"/rsvp2/{id}", ShowRsvp2)).Methods("GET")
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/rsvp2/{id}/save", SaveRsvp2)).Methods("POST")

  // api calls
  //r.HandleFunc(newrelic.WrapHandleFunc(app,"/api/rsvp", ListRsvp)).Methods("GET")
  //r.HandleFunc(newrelic.WrapHandleFunc(app,"/api/rsvp/{id}", GetRsvp)).Methods("GET")
  //r.HandleFunc(newrelic.WrapHandleFunc(app,"/api/rsvp/{id}", CreateRsvp)).Methods("POST")

  // static calls
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.ListenAndServe(fmt.Sprintf(":%s",port), r)
}
