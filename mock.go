package gobmock

import "fmt"

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
// The `mockScript` string will be inserted at the end.
func Mock(name string, mockScript string) GoBMock {
	return &mock{name: name, script: mockScript}
}

type mock struct {
	name   string
	script string
}

func (m *mock) MockContents() string {
	return m.mockFunction() + m.mockExport()
}

func (m *mock) mockExport() string {
	return fmt.Sprintf(exportDefinition, m.name)
}

func (m *mock) mockFunction() string {
	script := scriptStart + spyDefinition + mockDefinition + scriptEnd
	return fmt.Sprintf(script, m.name, m.script)
}
