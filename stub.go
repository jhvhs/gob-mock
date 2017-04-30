package gobmock

import "fmt"

// Produces a bash function with a given name.
func Stub(name string) GoBMock {
	return &stub{name: name, export: false}
}

// Produces a bash function with a given name
// that is exported to child processes
func ExportedStub(name string) GoBMock {
	return &stub{name: name, export: true}
}

type stub struct {
	name   string
	export bool
}

func (s *stub) MockContents() string {
	if s.export {
		return s.stubFunction() + s.stubExport()
	}
	return s.stubFunction()
}

func (s *stub) stubExport() string {
	return fmt.Sprintf(exportDefinition, s.name)
}

func (s *stub) stubFunction() string {
	script := scriptStart + stubDefinition + scriptEnd
	return fmt.Sprintf(script, s.name)
}
