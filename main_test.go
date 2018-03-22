package main

import (
  "os"
  "strings"
  "testing"
  "net/http"
  "net/http/httptest"
)

var a App

func TestMain(m *testing.M) {
  a = App{}
  a.Initialize(
    "root",
    "password1",
	"app_name",
	os.Getenv("NEWRELIC_LICENSE_KEY"),
  )

  code := m.Run()
  os.Exit(code)
}

func TestHomePage(t *testing.T) {
  clearTestData(t)
  req, _ := http.NewRequest("GET", "/", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)

  if body := response.Body.String(); !strings.Contains(body, "<h1>Fiona &amp; Gareth</h1>") {
    t.Errorf("Expected a correct title. Got %s", body)
  }
}

func TestShowRsvp(t *testing.T) {
  clearTestData(t)
  req, _ := http.NewRequest("GET", "/rsvp/1", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)

  if body := response.Body.String(); !strings.Contains(body, `<form method="post" action="/rsvp/1/save">`) {
    t.Errorf("Expected a correct form. Got %s", body)
  }
}

func TestCanUpdateRsvp(t *testing.T) {
  clearTestData(t)
  req, _ := http.NewRequest("GET", "/rsvp/1", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)

  body := response.Body.String()

  if !strings.Contains(body, `<form method="post" action="/rsvp/1/save">`) {
    t.Errorf("Expected a correct form. Got %s", body)
  }
  if !strings.Contains(body, `<input type="email" class="form-control" id="Email" placeholder="Email" value="bob1@bob.com">`) {
    t.Errorf("Expected a correct email field. Got %s", body)
  }
  if !strings.Contains(body, `<input type="text" class="form-control" id="Name" placeholder="Name" value="bob1">`) {
    t.Errorf("Expected a correct name field. Got %s", body)
  }

  // post update
  // follow redirect
  // check form
}

func clearTestData(t *testing.T) {
  batch := []string{
	`DELETE FROM rsvp;`,
	`INSERT INTO rsvp (rsvp_id, email, name, comments) VALUES ('1', 'bob1@bob.com','bob1','');`,
	`INSERT INTO rsvp (rsvp_id, email, name, comments) VALUES ('2', 'bob2@bob.com','bob2','');`,
	`INSERT INTO rsvp (rsvp_id, email, name, comments) VALUES ('3', 'bob3@bob.com','bob3','');`,
	`INSERT INTO rsvp (rsvp_id, email, name, comments) VALUES ('4', 'bob4@bob.com','bob4','');`,
	`INSERT INTO rsvp (rsvp_id, email, name, comments) VALUES ('5', 'bob5@bob.com','bob5','');`,
	`DELETE FROM guests;`,
	`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('1',1,'bobs friend','');`,
	`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('3',1,'friend 1','');`,
	`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('3',0,'friend 2','');`,
  }

  for _, b := range batch {
    _, err := a.DB.Exec(b)
    if err != nil {
      t.Fatalf("sql.Exec: Error: %s\n", err)
    }
  }
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
  rr := httptest.NewRecorder()
  a.Router.ServeHTTP(rr, req)
  return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
  if expected != actual {
    t.Errorf("Expected response code %d. Got %d\n", expected, actual)
  }
}
