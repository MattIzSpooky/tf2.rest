package quotes

type Quote struct {
	Class       string `json:"class"`
	Quote       string `json:"quote"`
	Dominates   string `json:"dominates,omitempty"`
	RevengeKill bool   `json:"revengeKill"`
	AudioFile   string `json:"audioFile"`
}

const (
	SCOUT    = "scout"
	SOLDIER  = "soldier"
	PYRO     = "pyro"
	DEMOMAN  = "demoman"
	HEAVY    = "heavy"
	ENGINEER = "engineer"
	MEDIC    = "medic"
	SNIPER   = "sniper"
	SPY      = "spy"
	NONE     = "none"
)

var ALL []Quote

func Setup() {
	// Attack classes
	ALL = append(ALL, scoutQuotes...)
	ALL = append(ALL, soldierQuotes...)
	ALL = append(ALL, pyroQuotes...)

	// Defense classes
	ALL = append(ALL, demomanQuotes...)
	ALL = append(ALL, heavyQuotes...)
	ALL = append(ALL, engineerQuotes...)

	// Support classes
	ALL = append(ALL, medicQuotes...)
	ALL = append(ALL, sniperQuotes...)
	ALL = append(ALL, spyQuotes...)
}
