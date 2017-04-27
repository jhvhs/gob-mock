package gobmock

import "fmt"

func Mock(name string, mockScript string) GoBMock {
	return &mock{name: name, script: mockScript, export: false}
}

func ExportedMock(name string, mockScript string) GoBMock {
	return &mock{name: name, script: mockScript, export: true}
}

type mock struct {
	name string
	export bool
	script string
}

func (m *mock)MockContents() string {
	if m.export {
		return m.mockFunction() + m.mockExport()
	}
	return m.mockFunction()
}


func (m *mock)mockExport() string {
	return fmt.Sprintf("\nexport -f %s\n", m.name)
}
func (m *mock)mockFunction() string {
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
		  %[2]s
		}`, m.name, m.script)
}