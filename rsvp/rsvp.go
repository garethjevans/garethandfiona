package rsvp

import (
  "fmt"
  "github.com/simplereach/timeutils"
)

type Rsvp struct {
  ID        int64            `json:"id,omitempty" schema:"-"`
  RsvpID    string           `json:"rsvp_id,omitempty" schema:"-"`
  RsvpDate  timeutils.Time   `json:"rsvp_time,omitempty" schema:"-"`
  Email     string           `json:"email,omitempty" schema:"Email"`
  Name      string           `json:"name,omitempty" schema:"Name"`
  Comments  string           `json:"comments,omitempty" schema:"Comments"`
  Guests    []*Guest         `json:"guests,omitempty"`
}

type Guest struct {
  ID        int64            `json:"id,omitempty" schema:"-"`
  RsvpID    string           `json:"rsvp_id,omitempty" schema:"-"`
  Name      string           `json:"name,omitempty" schema:"Name"`
  Attending bool             `json:"attending,omitempty" schema:"Attending"`
  Comments  string           `json:"comments,omitempty" schema:"Comments"`
}

func (r *Rsvp) String() string {
  return fmt.Sprintf("Rsvp(ID: %s, RsvpID: %s, Name: %s, Email: %s, Comments: %s, Guests: %s)", r.ID, r.RsvpID, r.Name, r.Email, r.Comments, r.Guests)
}

func (g *Guest) String() string {
  return fmt.Sprintf("Guest(ID: %s, RsvpID: %s, Name: %s, Attending: %s)", g.ID, g.RsvpID, g.Name, g.Attending)
}
