package models

type Player struct {
	ID                  string
	Nickname            string
	Guess               int
	GuessDiffWithSecret int
	Trophy              int
}

type PlayerResult struct {
	Rank        int    `json:"rank"`
	Player      string `json:"player"`
	Guess       int    `json:"guess,omitempty"`
	DeltaTrophy int    `json:"deltaTrophy"`
}
