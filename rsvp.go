package main

import (
	"fmt"
	"time"
)

type Rsvp struct {
	ID          int64      `json:"-" schema:"-"`
	RsvpID      string     `json:"-" schema:"-"`
	RsvpDate    *time.Time `json:"-" schema:"-"`
	ReplyType   string     `json:"-" schema:"-"`
	ReplyStatus string     `json:"status,omitempty" schema:"-"`
	Email       string     `json:"email,omitempty"`
	Guests      []*Guest   `json:"guests,omitempty"`
}

type Guest struct {
	ID        int64  `json:"-" schema:"-"`
	RsvpID    string `json:"-" schema:"-"`
	Name      string `json:"name,omitempty"`
	Attending bool   `json:"attending"`
	Comments  string `json:"comments,omitempty"`
}

func (r *Rsvp) String() string {
	return fmt.Sprintf("Rsvp(ID:%d, RsvpID:%s, Date:%s, Email:%s, Guests:%s)", r.ID, r.RsvpID, r.RsvpDate, r.Email, r.Guests)
}

func (r *Rsvp) WelcomeMessage() string {
	size := len(r.Guests)

	if size == 1 {
		return r.Guests[0].Name
	} else if size == 2 {
		return fmt.Sprintf("%s & %s", r.Guests[0].Name, r.Guests[1].Name)
	} else if size == 3 {
		return fmt.Sprintf("%s, %s & %s", r.Guests[0].Name, r.Guests[1].Name, r.Guests[2].Name)
	} else if size == 4 {
		return fmt.Sprintf("%s, %s, %s & %s", r.Guests[0].Name, r.Guests[1].Name, r.Guests[2].Name,
			r.Guests[3].Name)
	} else if size == 5 {
		return fmt.Sprintf("%s, %s, %s, %s & %s", r.Guests[0].Name, r.Guests[1].Name, r.Guests[2].Name,
			r.Guests[3].Name, r.Guests[4].Name)
	} else if size == 6 {
		return fmt.Sprintf("%s, %s, %s, %s, %s & %s", r.Guests[0].Name, r.Guests[1].Name, r.Guests[2].Name,
			r.Guests[3].Name, r.Guests[4].Name, r.Guests[5].Name)
	} else if size == 7 {
		return fmt.Sprintf("%s, %s, %s, %s, %s, %s & %s", r.Guests[0].Name, r.Guests[1].Name, r.Guests[2].Name,
			r.Guests[3].Name, r.Guests[4].Name, r.Guests[5].Name, r.Guests[6].Name)
	} else if size == 8 {
		return fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s & %s", r.Guests[0].Name, r.Guests[1].Name, r.Guests[2].Name,
			r.Guests[3].Name, r.Guests[4].Name, r.Guests[5].Name, r.Guests[6].Name, r.Guests[7].Name)
	} else if size == 9 {
		return fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s & %s", r.Guests[0].Name, r.Guests[1].Name, r.Guests[2].Name,
			r.Guests[3].Name, r.Guests[4].Name, r.Guests[5].Name, r.Guests[6].Name, r.Guests[7].Name, r.Guests[8].Name)
	}
	return ""
}

func (r *Rsvp) IsAttending() bool {
	for _, g := range r.Guests {
		if g.Attending {
			return true
		}
	}
	return false
}

func (g *Guest) String() string {
	return fmt.Sprintf("Guest(ID:%d, RsvpID:%s, Name:%s, Attending:%t, Comments:%s)", g.ID, g.RsvpID, g.Name, g.Attending, g.Comments)
}
