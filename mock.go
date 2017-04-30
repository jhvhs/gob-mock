package gobmock

import "fmt"

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// The `mockScript` string will be inserted at the end.
func Mock(name string, mockScript string) GoBMock {
	return &mock{name: name, script: mockScript, export: false}
}

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// The `mockScript` string will be inserted at the end.
//
// This function will be exported to the child processes
func ExportedMock(name string, mockScript string) GoBMock {
	return &mock{name: name, script: mockScript, export: true}
}

type mock struct {
	name   string
	export bool
	script string
}

func (m *mock) MockContents() string {
	if m.export {
		return m.mockFunction() + m.mockExport()
	}
	return m.mockFunction()
}

func (m *mock) mockExport() string {
	return fmt.Sprintf(exportDefinition, m.name)
}

func (m *mock) mockFunction() string {
	script := scriptStart + spyDefinition + mockDefinition + scriptEnd
	return fmt.Sprintf(script, m.name, m.script)
}
