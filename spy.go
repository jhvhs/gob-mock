package gobmock

import "fmt"

const unconditionalCallthrough = "📣"

// Produces a bash function with a given name.
// The function will report it's arguments
// as well as any data passed into it via STDIN.
// All reporting messages are sent to STDERR.
func Spy(name string) Gob {
	return &spy{name: name, callThroughCondition: ""}
}

func SpyAndCallThrough(name string) Gob {
	return &spy{name: name, callThroughCondition: unconditionalCallthrough}
}

func SpyAndConditionallyCallThrough(name string, callthroughCondition string) Gob {
	return &spy{name: name, callThroughCondition: callthroughCondition}
}

type spy struct {
	name                 string
	callThroughCondition string
}

func (s *spy) MockContents() string {
	return s.spyFunction() + s.spyExport()
}

func (s *spy) spyExport() string {
	return fmt.Sprintf(exportDefinition, s.name)
}

func (s *spy) spyFunction() string {
	script := scriptStart + spyDefinition
	if s.callThroughCondition == unconditionalCallthrough {
		script = script + callThroughDefinition
	} else if s.callThroughCondition != "" {
		script = script + "if " + s.callThroughCondition + "; then\n" + callThroughDefinition + "\nfi\n"
	}
	return fmt.Sprintf(script+scriptEnd, s.name)
}
