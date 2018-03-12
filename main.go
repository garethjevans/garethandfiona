package main

import (
  "github.com/magiconair/properties"
  "github.com/gorilla/schema"
  "github.com/gorilla/mux"
  "github.com/newrelic/go-agent"
  "html/template"
  "net/http"
  "encoding/json"
  "fmt"
  "log"
  "os"
  "errors"
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
  Rsvp      *Rsvp
  P         map[string]string
}

type Rsvp struct {
  ID        string   `json:"id,omitempty" schema:"-"`
  Email     string   `json:"email,omitempty" schema:"email"`
  Attendees []Attendee `json:"attendees,omitempty"`
}

type Attendee struct {
  Attending bool   `json:"attending,omitempty"`
  Name      string `json:"name,omitempty"`
  DietryRequirements    string `json:"dietry_requirements,omitempty"`
  Wine      string `json:"wine,omitempty"`
}

var responses []*Rsvp

func ListRsvp(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(responses)
}

func GetRsvpById(id string) (*Rsvp, error) {
  for _, item := range responses {
    if item.ID == id {
	  return item, nil
	}
  }
  return nil, errors.New("Unable to find entry")
}

func GetRsvp(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  item, err := GetRsvpById(params["id"])
  if err != nil {
    log.Print("Invalid reference: ", err)
    http.Error(w, err.Error(), http.StatusNotFound)
	return
  }
  json.NewEncoder(w).Encode(item)
}

func CreateRsvp(w http.ResponseWriter, r *http.Request) {
  log.Print("Create Rsvp")
  params := mux.Vars(r)
  var rsvp Rsvp
  _ = json.NewDecoder(r.Body).Decode(&rsvp)
  rsvp.ID = params["id"]
  responses = append(responses, &rsvp)
  json.NewEncoder(w).Encode(responses)
}


func ShowInvite(w http.ResponseWriter, r *http.Request) {
  showRsvpBase(w,r,"invite")
}

func ShowRsvp(w http.ResponseWriter, r *http.Request) {
  showRsvpBase(w,r,"show_rsvp")
}

func showRsvpBase(w http.ResponseWriter, r *http.Request, v string) {
  log.Print("Show Rsvp")
  params := mux.Vars(r)
  properties := getProperties()
  item, err := GetRsvpById(params["id"])
  if err != nil {
    log.Print("Invalid reference: ", err)
    http.Error(w, err.Error(), http.StatusNotFound)
	return
  }

  t, _ := template.ParseFiles(fmt.Sprintf("templates/%s.tmpl", v))
  err2 := t.Execute(w, Page {Rsvp: item, P: properties})
  if err2 != nil {
    log.Print("Unable to parse template: ", err2)
    http.Error(w, err.Error(), http.StatusInternalServerError)
	return
  }
}

func EditRsvp(w http.ResponseWriter, r *http.Request) {
  log.Print("Edit Rsvp")
  params := mux.Vars(r)
  properties := getProperties()
  item, err := GetRsvpById(params["id"])
  if err != nil {
    log.Print("Invalid reference: ", err)
    http.Error(w, err.Error(), http.StatusNotFound)
	return
  }

  t, _ := template.ParseFiles("templates/edit_rsvp.tmpl")
  err2 := t.Execute(w, Page {Rsvp: item, P: properties})
  if err2 != nil {
    log.Print("Unable to parse template: ", err2)
    http.Error(w, err2.Error(), http.StatusInternalServerError)
	return
  }
}

func SaveRsvp(w http.ResponseWriter, r *http.Request) {
  log.Print("Saving Rsvp")
  params := mux.Vars(r)

  item, err := GetRsvpById(params["id"])
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

  target := "http://" + r.Host + "/rsvp/" + item.ID
  log.Print("Sending Redirect: " + target)
  http.Redirect(w, r, target, http.StatusSeeOther)
}

func main() {
  responses = append(responses, &Rsvp{ID: "1", Email: "mark@bob.com", Attendees: []Attendee { Attendee{Attending: true, Name: "Mark", DietryRequirements: "I Love Eggs"}, Attendee{Attending: true, Name:"KDoggie Dog"} }})
  responses = append(responses, &Rsvp{ID: "2", Email: "rhi@bob.com", Attendees: []Attendee { Attendee{Attending: true, Name: "Rhi", DietryRequirements: "Sausage"} }})
  responses = append(responses, &Rsvp{ID: "3", Email: "ol@bob.com", Attendees: []Attendee { Attendee{Attending: true, Name: "Oli", DietryRequirements: ""},  Attendee{Attending: false, Name: "Plus 1"} }})

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
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/api", handler))
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/invite", handler))
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/ping", handler))

  // web calls
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/invite/{id}", ShowInvite)).Methods("GET")
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/rsvp/{id}", ShowRsvp)).Methods("GET")
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/rsvp/{id}/edit", EditRsvp)).Methods("GET")
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/rsvp/{id}/save", SaveRsvp)).Methods("POST")

  // api calls
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/api/rsvp", ListRsvp)).Methods("GET")
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/api/rsvp/{id}", GetRsvp)).Methods("GET")
  r.HandleFunc(newrelic.WrapHandleFunc(app,"/api/rsvp/{id}", CreateRsvp)).Methods("POST")

  // static calls
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.ListenAndServe(fmt.Sprintf(":%s",port), r)
}
