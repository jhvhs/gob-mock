package gobmock

import "fmt"

// Produces an bash function with a given name.
// The function will silently consume any input
// passed in via STDIN.
func Stub(name string) Gob {
	return &stub{
		name: name,
		shouldSkipReading: false,
	}
}

type stub struct {
	name              string
	shouldSkipReading bool
}

func (s *stub) MockContents() string {
	return s.stubFunction() + s.stubExport()
}

// Produces a stub with the reading portion disabled.
// It will not consume any data passed into it via STDIN.
func (s *stub) WithoutReading() Gob {
	return &stub{
		name: s.name,
		shouldSkipReading: true,
	}
}

func (s *stub) stubExport() string {
	return fmt.Sprintf(exportDefinition, s.name)
}

func (s *stub) stubFunction() string {
	script := scriptStart + s.stubDefinition() + scriptEnd
	return fmt.Sprintf(script, s.name)
}

func (s *stub) stubDefinition() string {
	if s.shouldSkipReading {
		return "# Empty stub"
	}
	return stubDefinition
}
