package gobmock

import "fmt"

// Produces an empty function with a given name
// the function will report the arguments
// and the STDIN if it was passed to it
// All reporting messages are sent to STDERR
func Spy(name string) GoBMock {
	return &spy{name: name, export: false}
}

// Produces an empty function with a given name
// the function will report the arguments
// and the STDIN if it was passed to it
// All reporting messages are sent to STDERR
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
	return fmt.Sprintf("\nexport -f %s\n", s.name)
}

func (s *spy)spyFunction() string {
	return fmt.Sprintf(
		`%[1]s() {
		  local in_line_count=0
		  declare -a in_lines
		  while read -t0.05; do
		    in_lines[in_line_count]="$REPLY"
		    in_line_count=$(expr ${in_line_count} + 1)
		  done
		  callCounter=$(expr ${callCounter} + 1)
		  echo "[$callCounter] %[1]s $@" > /dev/fd/2
		  if [ ${in_line_count} -gt 0 ]; then
		    echo "[$callCounter received] input:" > /dev/fd/2
		    printf '%%s\n' "${in_lines[@]}" > /dev/fd/2
		    echo "[$callCounter end received]" > /dev/fd/2
		  fi
		}`, s.name)
}
