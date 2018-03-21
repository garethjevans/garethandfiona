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
  req, _ := http.NewRequest("GET", "/", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)

  if body := response.Body.String(); !strings.Contains(body, "<h1>Fiona &amp; Gareth</h1>") {
    t.Errorf("Expected a correct title. Got %s", body)
  }
}

func TestShowRsvp(t *testing.T) {
  req, _ := http.NewRequest("GET", "/rsvp/1", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)

  if body := response.Body.String(); !strings.Contains(body, `<form method="post" action="/rsvp/1/save">`) {
    t.Errorf("Expected a correct title. Got %s", body)
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
