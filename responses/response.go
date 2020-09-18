package responses

type Response struct {
	Class     string `json:"class"`
	Response  string `json:"response"`
	AudioFile string `json:"audioFile"`
	Type      string `json:"type"`                // example: Kill-related responses
	SubType   string `json:"subType,omitempty"`   // example: Payload-related responses
	Context   string `json:"context,omitempty"`   // example: Destroying a building
	Condition string `json:"condition,omitempty"` // example: Melee killing a Heavy
}

var All []Response

func Setup() {
	// Attack classes
	All = append(All, scoutResponses...)
	All = append(All, soldierResponses...)
	All = append(All, pyroResponses...)

	// Defense classes
	All = append(All, demomanResponses...)
	All = append(All, heavyResponses...)
	All = append(All, engineerResponses...)

	// Support classes
	All = append(All, medicResponses...)
	All = append(All, sniperResponses...)
	All = append(All, spyResponses...)
}

func FilterByClass(class string) []Response {
	var responses []Response

	for _, value := range All {
		if value.Class == class {
			responses = append(responses, value)
		}
	}

	return responses
}
