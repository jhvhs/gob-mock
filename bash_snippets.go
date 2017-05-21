package gobmock

const (
	scriptStart    = "%[1]s() {\n"
	stubDefinition = `
          while read -r -t0.1; do
            :
          done
  	`
	spyDefinition = `
          local in_lines
          while read -r -t0.1; do
            in_lines="${in_lines}${REPLY}
"
          done
          callCounter=$((callCounter + 1))
          echo "<${callCounter}> %[1]s $*" > /dev/fd/2
          if [ -n "${in_lines}" ]; then
            in_lines=${in_lines::-1}
            echo "<${callCounter} received> input:" > /dev/fd/2
            echo -n "${in_lines}" > /dev/fd/2
            echo "<${callCounter} end received>" > /dev/fd/2
          fi
        `
	mockDefinition   = "%[2]s\n"
	scriptEnd        = "\n return 0 \n } \n"
	exportDefinition = "export -f %s\n"

	callThroughDefinition = `
	  if [ -n "${in_lines}" ]; then
	    echo -n "${in_lines}" | $(which %[1]s) "$@"
	  else
	    $(which %[1]s) "$@"
	  fi
  	`
)
