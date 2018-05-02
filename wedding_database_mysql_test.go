package main

import (
	"log"
	"os"
	"testing"
)

var b App

func TestDatabase(t *testing.T) {
	// <setup code>
	b = App{}
	b.Initialize(
		"root",
		"password1",
		"app_name",
		os.Getenv("NEWRELIC_LICENSE_KEY"),
	)

	log.Printf("Starting Transaction")

	tx, err := b.DB.DB().Begin()
	if err != nil {
		log.Fatal("Unable to start transaction")
	}

	t.Run("loadRsvpById", loadRsvpById)
	t.Run("canUpdateRsvp", canUpdateRsvp)

	// <tear-down code>
	log.Printf("Rolling Back Transaction")
	tx.Rollback()
}

func loadRsvpById(t *testing.T) {
	clearTestDataForDatabase(t, b)
	// bob1
	rsvp, err := b.DB.GetRsvp("1")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if rsvp.RsvpID != "1" {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "1")
	}

	if rsvp.Email != "bob1@bob.com" {
		t.Fatalf("Expecting rsvp with email %s", "bob1@bob.com")
	}

	if len(rsvp.Guests) != 2 {
		t.Fatalf("Expecting rsvp with 2 guest")
	}

	// bob2
	rsvp, err = b.DB.GetRsvp("2")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if rsvp.RsvpID != "2" {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "2")
	}

	if rsvp.Email != "bob2@bob.com" {
		t.Fatalf("Expecting rsvp with email %s", "bob2@bob.com")
	}

	if len(rsvp.Guests) != 1 {
		t.Fatalf("Expecting rsvp with 1 guest")
	}

	// bob3
	rsvp, err = b.DB.GetRsvp("3")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if rsvp.RsvpID != "3" {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "3")
	}

	if rsvp.Email != "bob3@bob.com" {
		t.Fatalf("Expecting rsvp with email %s", "bob3@bob.com")
	}

	if len(rsvp.Guests) != 3 {
		t.Fatalf("Expecting rsvp with 3 guests")
	}
}

func canUpdateRsvp(t *testing.T) {
	clearTestDataForDatabase(t, b)

	// bob1
	rsvp, err := b.DB.GetRsvp("1")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if rsvp.RsvpID != "1" {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "1")
	}

	if rsvp.RsvpDate != nil {
		t.Fatalf("Expecting rsvp with nill rsvp_date but got %s", rsvp.RsvpDate)
	}

	if rsvp.Email != "bob1@bob.com" {
		t.Fatalf("Expecting rsvp with email %s", "bob1@bob.com")
	}

	if len(rsvp.Guests) != 2 {
		t.Fatalf("Expecting rsvp with 2 guest")
	}

	rsvp.Email = "bober@bobest.com"
	rsvp.Guests[0].Attending = false

	b.DB.UpdateRsvp(rsvp)

	rsvp, err = b.DB.GetRsvp("1")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if rsvp.RsvpID != "1" {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "1")
	}

	if rsvp.RsvpDate == nil {
		t.Fatalf("Expecting rsvp with rsvp_date %s", rsvp.RsvpDate)
	}

	if rsvp.Email != "bober@bobest.com" {
		t.Fatalf("Expecting rsvp with email %s", "bober@bobest.com")
	}

	if len(rsvp.Guests) != 2 {
		t.Fatalf("Expecting rsvp with 2 guest")
	}

	if rsvp.Guests[0].Attending {
		t.Fatalf("Guest 0 is not supposed to be attending")
	}
}

func clearTestDataForDatabase(t *testing.T, a App) {
	batch := []string{
		`DELETE FROM rsvp;`,
		`INSERT INTO rsvp (rsvp_id, reply_type, reply_status, email) VALUES ('1', '', '', 'bob1@bob.com');`,
		`INSERT INTO rsvp (rsvp_id, reply_type, reply_status, email) VALUES ('2', '', '', 'bob2@bob.com');`,
		`INSERT INTO rsvp (rsvp_id, reply_type, reply_status, email) VALUES ('3', '', '', 'bob3@bob.com');`,
		`INSERT INTO rsvp (rsvp_id, reply_type, reply_status, email) VALUES ('4', '', '', 'bob4@bob.com');`,
		`INSERT INTO rsvp (rsvp_id, reply_type, reply_status, email) VALUES ('5', '', '', 'bob5@bob.com');`,
		`DELETE FROM guests;`,
		`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('1',1,'bob1','');`,
		`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('1',1,'bobs friend','');`,
		`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('2',1,'bob2','');`,
		`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('3',1,'bob3','');`,
		`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('3',1,'friend 1','');`,
		`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('3',0,'friend 2','');`,
		`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('4',1,'bob4','');`,
		`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('5',1,'bob5','');`,
	}

	for _, b := range batch {
		_, err := a.DB.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}
}
