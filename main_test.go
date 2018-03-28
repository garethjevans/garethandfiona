package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var a App

func TestApp(t *testing.T) {
	// <setup code>
	a = App{}
	a.Initialize(
		"root",
		"password1",
		"app_name",
		os.Getenv("NEWRELIC_LICENSE_KEY"),
	)

	log.Printf("Starting Transaction")

	tx, err := a.DB.DB().Begin()
	if err != nil {
		log.Fatal("Unable to start transaction")
	}

	t.Run("homePageViaWeb", homePageViaWeb)
	t.Run("showRsvpViaWeb", showRsvpViaWeb)
	t.Run("canUpdateRsvpViaWeb", canUpdateRsvpViaWeb)

	t.Run("showRsvpViaRest", showRsvpViaRest)
	t.Run("showMissingRsvpViaRest", showMissingRsvpViaRest)
	t.Run("canUpdateRsvpViaRest", canUpdateRsvpViaRest)

	// <tear-down code>
	log.Printf("Rolling Back Transaction")
	tx.Rollback()
}

func homePageViaWeb(t *testing.T) {
	clearTestData(t, a)
	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); !strings.Contains(body, "<h1>Fiona &amp; Gareth</h1>") {
		t.Errorf("Expected a correct title. Got %s", body)
	}
}

func showRsvpViaWeb(t *testing.T) {
	clearTestData(t, a)
	req, _ := http.NewRequest("GET", "/rsvp/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); !strings.Contains(body, `<form method="post" action="/rsvp/1/save">`) {
		t.Errorf("Expected a correct form. Got %s", body)
	}
}

func canUpdateRsvpViaWeb(t *testing.T) {
	clearTestData(t, a)
	req, _ := http.NewRequest("GET", "/rsvp/1", nil)

	log.Printf("Requesting GET /rsvp/1")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()

	if !strings.Contains(body, `<form method="post" action="/rsvp/1/save">`) {
		t.Errorf("Expected a correct form. Got %s", body)
	}
	if !strings.Contains(body, `<input type="email" class="form-control" name="Email" placeholder="Email" value="bob1@bob.com">`) {
		t.Errorf("Expected a correct email field. Got %s", body)
	}
	if !strings.Contains(body, `<input type="text" class="form-control" name="Name" placeholder="Name" value="bob1">`) {
		t.Errorf("Expected a correct name field. Got %s", body)
	}
	if !strings.Contains(body, `<input type="text" class="form-control" name="Guests.0.Name" value="bobs friend">`) {
		t.Errorf("Expected a correct guest 1 name. Got %s", body)
	}
	if !strings.Contains(body, `<input type="text" class="form-control" name="Guests.0.Comments" value="">`) {
		t.Errorf("Expected a correct guest 1 comments. Got %s", body)
	}
	if !strings.Contains(body, `<label class="radio-inline"><input type="radio" class="form-control" name="Guests.0.Attending" value="true" checked>Yes</label>`) {
		t.Errorf("Expected a correct guest 1 attending. Got %s", body)
	}
	if !strings.Contains(body, `<label class="radio-inline"><input type="radio" class="form-control" name="Guests.0.Attending" value="false">No</label>`) {
		t.Errorf("Expected a correct guest 1 attending. Got %s", body)
	}

	// post update
	postBody := `Status=attending&Name=bobnew&Email=bobnew@bob.com&Guests.0.Name=belinda&Guests.0.Comments=Loves Eggs&Guests.0.Attending=false`
	postBodyReader := strings.NewReader(postBody)

	// follow redirect
	req, _ = http.NewRequest("POST", "/rsvp/1/save", postBodyReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("Requesting POST /rsvp/1/save")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusSeeOther, response.Code)

	// check form
	req, _ = http.NewRequest("GET", "/rsvp/1", nil)
	log.Printf("Requesting GET /rsvp/1")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	body = response.Body.String()

	if !strings.Contains(body, `<form method="post" action="/rsvp/1/save">`) {
		t.Errorf("Expected a correct form. Got %s", body)
	}
	if !strings.Contains(body, `<input type="email" class="form-control" name="Email" placeholder="Email" value="bobnew@bob.com">`) {
		t.Errorf("Expected a correct email field. Got %s", body)
	}
	if !strings.Contains(body, `<input type="text" class="form-control" name="Name" placeholder="Name" value="bobnew">`) {
		t.Errorf("Expected a correct name field. Got %s", body)
	}

	if !strings.Contains(body, `<input type="text" class="form-control" name="Guests.0.Name" value="belinda">`) {
		t.Errorf("Expected a correct guest 1 name. Got %s", body)
	}
	if !strings.Contains(body, `<input type="text" class="form-control" name="Guests.0.Comments" value="Loves Eggs">`) {
		t.Errorf("Expected a correct guest 1 comments. Got %s", body)
	}
	if !strings.Contains(body, `<input type="radio" class="form-control" name="Guests.0.Attending" value="true">`) {
		t.Errorf("Expected a correct guest 1 attending. Got %s", body)
	}
	if !strings.Contains(body, `<input type="radio" class="form-control" name="Guests.0.Attending" value="false" checked>`) {
		t.Errorf("Expected a correct guest 1 attending. Got %s", body)
	}
}

func showRsvpViaRest(t *testing.T) {
	clearTestData(t, a)
	req, _ := http.NewRequest("GET", "/api/rsvp/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != `{"email":"bob1@bob.com","name":"bob1","guests":[{"name":"bobs friend","attending":true}]}` {
		t.Errorf("Expected a correct json body. Got %s", body)
	}
}

func showMissingRsvpViaRest(t *testing.T) {
	clearTestData(t, a)
	req, _ := http.NewRequest("GET", "/api/rsvp/missing", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	body := response.Body.String()
	if body != `{"error":"Unable to find 'missing'"}` {
		t.Errorf("Expected a correct json body. Got %s", body)
	}
}

func canUpdateRsvpViaRest(t *testing.T) {
	clearTestData(t, a)
	req, _ := http.NewRequest("GET", "/api/rsvp/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()

	if body != `{"email":"bob1@bob.com","name":"bob1","guests":[{"name":"bobs friend","attending":true}]}` {
		t.Errorf("Expected a correct json body. Got %s", body)
	}

	// post update
	postBody := `{"status":"attending"}`
	postBodyReader := strings.NewReader(postBody)

	// follow redirect
	req, _ = http.NewRequest("POST", "/api/rsvp/1", postBodyReader)
	req.Header.Set("Content-Type", "application/json")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusSeeOther, response.Code)

	// check form
	req, _ = http.NewRequest("GET", "/api/rsvp/1", nil)
	log.Printf("Requesting GET /api/rsvp/1")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	body = response.Body.String()

	if body != `{"status":"attending","email":"bob1@bob.com","name":"bob1","guests":[{"name":"bobs friend","attending":true}]}` {
		t.Errorf("Expected a correct json body. Got %s", body)
	}
}

func clearTestData(t *testing.T, a App) {
	batch := []string{
		`DELETE FROM rsvp;`,
		`INSERT INTO rsvp (rsvp_id, status, email, name, comments) VALUES ('1', '', 'bob1@bob.com','bob1','');`,
		`INSERT INTO rsvp (rsvp_id, status, email, name, comments) VALUES ('2', '', 'bob2@bob.com','bob2','');`,
		`INSERT INTO rsvp (rsvp_id, status, email, name, comments) VALUES ('3', '', 'bob3@bob.com','bob3','');`,
		`INSERT INTO rsvp (rsvp_id, status, email, name, comments) VALUES ('4', '', 'bob4@bob.com','bob4','');`,
		`INSERT INTO rsvp (rsvp_id, status, email, name, comments) VALUES ('5', '', 'bob5@bob.com','bob5','');`,
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
