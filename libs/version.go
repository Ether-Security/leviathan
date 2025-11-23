package libs

import "fmt"

const (
	VERSION   = "v1.0.0"
	NAME      = "leviathan"
	SHORTNAME = "lvt"
	DESC      = "A workflow engine for OffSec inspired by Osmedeus"
	AUTHOR    = "@y0no"
	DOCS      = "https://leviathan.y0no.fr"
)

var LOGDIR = fmt.Sprintf("/tmp/%s-log", SHORTNAME)
var CONFIGFILE = fmt.Sprintf("~/.config/%s/config.yaml", NAME)
