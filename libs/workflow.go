package libs

type Workflow struct {
	Name        string
	Description string
	Author      string
	Validator   string

	Params   []map[string]string
	Routines []Routine
}

type Routine struct {
	Modules []string
	Params  []map[string]string
}
