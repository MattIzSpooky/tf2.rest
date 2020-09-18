package class

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

var Classes = [...]string{
	SCOUT,
	SOLDIER,
	PYRO,
	DEMOMAN,
	HEAVY,
	ENGINEER,
	MEDIC,
	SNIPER,
	SPY,
}

func Contains(value string) bool {
	for _, a := range Classes {
		if a == value {
			return true
		}
	}
	return false
}