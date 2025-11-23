package libs

type Options struct {
	Scan        Scan
	Log         Log
	ConfigFile  string
	Environment Environment
}

type Environment struct {
	Workspaces string
	Workflows  string
	Modules    string
	Binaries   string
}

type Scan struct {
	Flow    string
	Threads int
	Targets []string
	Params  map[string]string
	Output  string
	NoClean bool
	Resume  bool
}

type Log struct {
	JSON      bool
	Directory string
	Debug     bool
	Quiet     bool
}
