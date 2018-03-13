package rsvp

type Rsvp struct {
  ID        int64           `json:"id,omitempty" schema:"-"`
  RsvpID    string           `json:"id,omitempty" schema:"-"`
  Email     string           `json:"email,omitempty" schema:"email"`
  Attendees []Attendee       `json:"attendees,omitempty"`
}

type Attendee struct {
  ID        int64            `json:"id,omitempty" schema:"-"`
  RsvpID    string           `json:"id,omitempty" schema:"-"`
  Attending bool             `json:"attending,omitempty"`
  Name      string           `json:"name,omitempty"`
  DietryRequirements  string `json:"dietry_requirements,omitempty"`
  Wine      string           `json:"wine,omitempty"`
}
