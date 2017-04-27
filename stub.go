package gobmock

import "fmt"

// Produces an empty function with a given name
func Stub(name string) GoBMock {
	return &stub{name: name, export: false}
}

// Produces an empty function with a given name
// that is exported to child processes
func ExportedStub(name string) GoBMock {
	return &stub{name: name, export: true}
}

type stub struct {
	name string
	export bool
}

func (s *stub) MockContents() string {
	if s.export {
		return s.stubFunction() + s.stubExport()
	}
	return s.stubFunction()
}

func (s *stub) stubExport() string {
	return fmt.Sprintf("\nexport -f %s\n", s.name)
}

func (s *stub) stubFunction() string {
	return fmt.Sprintf(
		`%s() {
		  while read -t0.05; do
		    :
		  done
		}`,
		s.name)
}
