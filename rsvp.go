package main

import (
	"fmt"
	"github.com/simplereach/timeutils"
)

type Rsvp struct {
	ID       int64          `json:"id,omitempty" schema:"-"`
	RsvpID   string         `json:"rsvp_id,omitempty" schema:"-"`
	RsvpDate timeutils.Time `json:"rsvp_time,omitempty" schema:"-"`
	Email    string         `json:"email,omitempty"`
	Name     string         `json:"name,omitempty"`
	Comments string         `json:"comments,omitempty"`
	Guests   []*Guest       `json:"guests,omitempty"`
}

type Guest struct {
	ID        int64  `json:"id,omitempty" schema:"-"`
	RsvpID    string `json:"rsvp_id,omitempty" schema:"-"`
	Name      string `json:"name,omitempty"`
	Attending bool   `json:"attending,omitempty"`
	Comments  string `json:"comments,omitempty"`
}

func (r *Rsvp) String() string {
	return fmt.Sprintf("Rsvp(ID: %d, RsvpID: %s, Name: %s, Email: %s, Comments: %s, Guests: %s)", r.ID, r.RsvpID, r.Name, r.Email, r.Comments, r.Guests)
}

func (g *Guest) String() string {
	return fmt.Sprintf("Guest(ID: %d, RsvpID: %s, Name: %s, Attending: %t)", g.ID, g.RsvpID, g.Name, g.Attending)
}
