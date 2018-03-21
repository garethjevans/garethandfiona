package rsvp

import (
  "github.com/simplereach/timeutils"
)

type Rsvp struct {
  ID        int64            `json:"id,omitempty" schema:"-"`
  RsvpID    string           `json:"rsvp_id,omitempty" schema:"-"`
  RsvpDate  timeutils.Time   `json:"rsvp_time,omitempty" schema:"-"`
  Email     string           `json:"email,omitempty" schema:"email"`
  Name      string           `json:"name,omitempty" schema:"name"`
  Comments  string           `json:"comments,omitempty" schema:"comments"`
  Guests    []*Guest         `json:"guests,omitempty"`
}

type Guest struct {
  ID        int64            `json:"id,omitempty" schema:"-"`
  RsvpID    string           `json:"rsvp_id,omitempty" schema:"-"`
  Name      string           `json:"name,omitempty"`
  Attending bool             `json:"attending,omitempty"`
  Comments  string           `json:"comments,omitempty"`
}
