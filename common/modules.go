package common

import (
	"time"
)

// ModulePublic -
type ModulePublic struct {
	Path      string        `json:",omitempty"` // module path
	Version   string        `json:",omitempty"` // module version
	Versions  []string      `json:",omitempty"` // available module versions
	Replace   *ModulePublic `json:",omitempty"` // replaced by this module
	Time      *time.Time    `json:",omitempty"` // time version was created
	Update    *ModulePublic `json:",omitempty"` // available update (with -u)
	Main      bool          `json:",omitempty"` // is this the main module?
	Indirect  bool          `json:",omitempty"` // module is only indirectly needed by main module
	Dir       string        `json:",omitempty"` // directory holding local copy of files, if any
	GoMod     string        `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion string        `json:",omitempty"` // go version used in module
	Error     *ModuleError  `json:",omitempty"` // error loading module
}

// ModuleError -
type ModuleError struct {
	Err string // error text
}
