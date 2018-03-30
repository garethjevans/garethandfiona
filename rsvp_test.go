package main

import (
	"testing"
)

func TestWelcomeMessage(t *testing.T) {

	var guests []*Guest

	guests = append(guests, &Guest{Name: "Bob"})
	r := Rsvp{Guests: guests}
	assertEquals(t, r.WelcomeMessage(), "Bob")

	guests = append(guests, &Guest{Name: "Belinda"})
	r = Rsvp{Guests: guests}
	assertEquals(t, r.WelcomeMessage(), "Bob & Belinda")

	guests = append(guests, &Guest{Name: "Barbara"})
	r = Rsvp{Guests: guests}
	assertEquals(t, r.WelcomeMessage(), "Bob, Belinda & Barbara")
}

func assertEquals(t *testing.T, actual string, expected string) {
	if expected != actual {
		t.Errorf("Expected '%s' got '%s'", expected, actual)
	}
}
