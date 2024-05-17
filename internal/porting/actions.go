package porting

type ActionList struct {
	Modules  []ModuleAction
	Packages []PackageAction
}

type ModuleAction struct {
	Path     string
	Version  string
	Fixed    string
	Dir      string `json:",omitempty"`
	Imported bool
}

type PackageAction struct {
	Path    string
	Module  string
	Dir     string `json:",omitempty"`
	Actions []map[string]any
	Files   []FileAction
	Error   string `json:",omitempty"`
}

type FileAction struct {
	Name     string
	Build    bool
	BaseFile string `json:",omitempty"`
	Lines    []LineAction `json:",omitempty"`
}

type LineAction struct {
	Line     uint
	Original string
	Fixed    string
	Changes  []ChangeAction
}

type ChangeAction struct {
	Column      uint
	Original    string
	Replacement string
}

type Error struct {
	IsPortError bool
	Error       string
}
