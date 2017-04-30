package gobmock

import "fmt"

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
func Spy(name string) GoBMock {
	return &spy{name: name, export: false}
}

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
//
// This function will be exported to the child processes
func ExportedSpy(name string) GoBMock {
	return &spy{name: name, export: true}
}

type spy struct {
	name   string
	export bool
}

func (s *spy) MockContents() string {
	if s.export {
		return s.spyFunction() + s.spyExport()
	}
	return s.spyFunction()
}
func (s *spy) spyExport() string {
	return fmt.Sprintf(exportDefinition, s.name)
}

func (s *spy) spyFunction() string {
	script := scriptStart + spyDefinition + scriptEnd
	return fmt.Sprintf(script, s.name)
}
