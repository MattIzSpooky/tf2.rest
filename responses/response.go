package responses

type Response struct {
	Class     string `json:"class"`
	Response  string `json:"Response"`
	AudioFile string `json:"audioFile"`
	Type      string `json:"type"`      // example: Kill-related responses
	SubType	  string `json:"subType"` 	// example: Payload-related responses
	Context   string `json:"context"`   // example: Destroying a building
	Condition string `json:"condition"` // example: Melee killing a Heavy
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
)

var ALL []Response

func Setup() {
	// Attack classes
	ALL = append(ALL, scoutResponses...)
	ALL = append(ALL, soldierResponses...)
	ALL = append(ALL, pyroResponses...)

	// Defense classes
	ALL = append(ALL, demomanResponses...)
	ALL = append(ALL, heavyResponses...)
	ALL = append(ALL, engineerResponses...)

	// Support classes
	ALL = append(ALL, medicResponses...)
	ALL = append(ALL, sniperResponses...)
	ALL = append(ALL, spyResponses...)
}
