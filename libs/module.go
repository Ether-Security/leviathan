package libs

type Module struct {
	Name        string
	Description string

	Params []map[string]string

	Reports []string
	PreRun  []string `yaml:"pre_run"`
	Steps   []Step
	PostRun []string `yaml:"post_run"`
}

type Step struct {
	Requirements []string
	Conditions   []string

	// Define specific source file for a step
	Source string

	// Executed if conditions are met
	Commands []string
	Scripts  []string

	// Executed if conditions are not met
	RCommands []string
	RScripts  []string
}
