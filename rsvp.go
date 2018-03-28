package main

import (
	"fmt"
	"time"
)

type Rsvp struct {
	ID       int64      `json:"-" schema:"-"`
	RsvpID   string     `json:"-" schema:"-"`
	RsvpDate *time.Time `json:"rsvp_time,omitempty" schema:"-"`
	Status   string     `json:"status,omitempty"`
	Email    string     `json:"email,omitempty"`
	Name     string     `json:"name,omitempty"`
	Comments string     `json:"comments,omitempty"`
	Guests   []*Guest   `json:"guests,omitempty"`
}

type Guest struct {
	ID        int64  `json:"-" schema:"-"`
	RsvpID    string `json:"-" schema:"-"`
	Name      string `json:"name,omitempty"`
	Attending bool   `json:"attending,omitempty"`
	Comments  string `json:"comments,omitempty"`
}

func (r *Rsvp) String() string {
	return fmt.Sprintf("Rsvp(ID:%d, RsvpID:%s, Date:%s, Name:%s, Email:%s, Status:%s, Comments:%s, Guests:%s)", r.ID, r.RsvpID, r.RsvpDate, r.Name, r.Email, r.Status, r.Comments, r.Guests)
}

func (g *Guest) String() string {
	return fmt.Sprintf("Guest(ID:%d, RsvpID:%s, Name:%s, Attending:%t, Comments:%s)", g.ID, g.RsvpID, g.Name, g.Attending, g.Comments)
}
