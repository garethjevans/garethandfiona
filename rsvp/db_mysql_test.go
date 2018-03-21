package rsvp

import (
  "testing"
)

func TestLoadRsvpById(t *testing.T) {
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

    db, err := NewMySQLDB(MySQLConfig{ Username:"root", Password:"password1", Host:"127.0.0.1", Port:3306 })
    if err != nil {
      t.Fatalf("Unable to connect to database: ", err)
    }

	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}

	defer db.Close()

	for _, b := range batch {
		_, err = db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}

	// bob1
	rsvp, err := db.GetRsvp("1")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if (rsvp.RsvpID != "1") {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "1")
	}

	if (rsvp.Email != "bob1@bob.com") {
		t.Fatalf("Expecting rsvp with email %s", "bob1@bob.com")
	}

	if (len(rsvp.Guests) != 1) {
		t.Fatalf("Expecting rsvp with 1 guest")
	}

	// bob2
	rsvp, err = db.GetRsvp("2")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if (rsvp.RsvpID != "2") {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "2")
	}

	if (rsvp.Email != "bob2@bob.com") {
		t.Fatalf("Expecting rsvp with email %s", "bob2@bob.com")
	}

	if (len(rsvp.Guests) != 0) {
		t.Fatalf("Expecting rsvp with 0 guest")
	}

	// bob3
	rsvp, err = db.GetRsvp("3")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if (rsvp.RsvpID != "3") {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "3")
	}

	if (rsvp.Email != "bob3@bob.com") {
		t.Fatalf("Expecting rsvp with email %s", "bob3@bob.com")
	}

	if (len(rsvp.Guests) != 2) {
		t.Fatalf("Expecting rsvp with 2 guests")
	}
}

func TestCanUpdateRsvp(t *testing.T) {
    batch := []string{
	`DELETE FROM rsvp;`,
	`INSERT INTO rsvp (rsvp_id, email, name, comments) VALUES ('1', 'bob1@bob.com','bob1','');`,
	`DELETE FROM guests;`,
	`INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('1',1,'bobs friend','');`,
	}

    db, err := NewMySQLDB(MySQLConfig{ Username:"root", Password:"password1", Host:"127.0.0.1", Port:3306 })
    if err != nil {
      t.Fatalf("Unable to connect to database: ", err)
    }

	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}

	defer db.Close()

	for _, b := range batch {
		_, err = db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}

	// bob1
	rsvp, err := db.GetRsvp("1")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if (rsvp.RsvpID != "1") {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "1")
	}

	if (rsvp.Email != "bob1@bob.com") {
		t.Fatalf("Expecting rsvp with email %s", "bob1@bob.com")
	}

	if (len(rsvp.Guests) != 1) {
		t.Fatalf("Expecting rsvp with 1 guest")
	}

	rsvp.Name = "Bob Full Name"
	rsvp.Email = "bober@bobest.com"

	rsvp.Guests[0].Attending = false

	db.UpdateRsvp(rsvp)

	rsvp, err = db.GetRsvp("1")
	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}

	if (rsvp.RsvpID != "1") {
		t.Fatalf("Expecting rsvp with rsvp_id %s", "1")
	}

	if (rsvp.Email != "bober@bobest.com") {
		t.Fatalf("Expecting rsvp with email %s", "bober@bobest.com")
	}

	if (len(rsvp.Guests) != 1) {
		t.Fatalf("Expecting rsvp with 1 guest")
	}

	if (rsvp.Guests[0].Attending) {
		t.Fatalf("Guest 1 is not supposed to be attending")
	}
}
