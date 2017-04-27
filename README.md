# gob-mock
A simple mocking framework for bash testing with [golang](https://golang.org/) via [go-basher](https://github.com/progrium/go-basher/)

## Quick intro

**gob-mock** provides tree types of mocks for stubbing executables or functions:

### Stubs
Stub is the most simple version of a mock. It will silently drop the call to the original executable. Additionally, it would silently discard any data that was piped through it, and would not cause any problems when the `-o pipefail` option is enabled.

```go
  bash := basher.NewContext("/path/to/bash", false)
  mocks := []GoBMock{Stub("wget")}
  ApplyMocks(bash, mocks)
  status, _ := bash.Run("wget", []string{"zaa://qwee.dooo"})
  Expect(status).To(Equal(0))
```

### Spies
Spy does everything a stub does, but in addition, it will print the function name, the arguments used in a call to STDERR. If any data was piped in, it will also be reported in the same way.

```go
  bash := basher.NewContext("/path/to/bash", false)
  bash.StdErr = gbytes.NewBuffer()
  mocks := []GoBMock{Spy("wget")}
  ApplyMocks(bash, mocks)
  status, _ := bash.Run("wget", []string{"zaa://qwee.dooo", "fus", "ro", "dah"})
  Expect(status).To(Equal(0))
  Expect(stderr).To(gbytes.Say("[1] wget zaa://qwee.dooo fus ro dah"))
```

### Mocks
A mock does everything that a spy does, but also provides an entry point into the mocking function. 

The mock may produce output which depends on the supplied arguments:
```go
  bash := basher.NewContext("/path/to/bash", false)
  bash.StdErr = gbytes.NewBuffer()
  bash.StdOut = gbytes.NewBuffer()
  mocks := []GoBMock{Mock("wget", "if [[ $1 == 'quux' ]]; then echo 'Yes'; else echo 'No'; fi")}
  ApplyMocks(bash, mocks)
  status, _ := bash.Run("wget", []string{"quux", "fus", "ro", "dah"})
  Expect(status).To(Equal(0))
  Expect(stderr).To(gbytes.Say("[1] wget quux fus ro dah"))
  Expect(stdout).To(gbytes.Say("Yes"))
```

An exit code can be simulated by using the `return` keyword with the appropriate number:

```go
  bash := basher.NewContext("/path/to/bash", false)
  mocks := []GoBMock{Mock("wget", "return 12")}
  ApplyMocks(bash, mocks)
  status, _ := bash.Run("wget", []string{"quux", "fus", "ro", "dah"})
  Expect(status).To(Equal(12))
```

If a mock is referring to an executable, it can use `which` to put the call through when needed, while still recording the invocation details:

```go
  bash := basher.NewContext("/path/to/bash", false)
  bash.StdErr = gbytes.NewBuffer()
  bash.StdOut = gbytes.NewBuffer()
  mocks := []GoBMock{Mock("curl", `$(which curl) "$@"`)}
  ApplyMocks(bash, mocks)
  status, _ := bash.Run("curl", []string{"https://www.google.ie/"})
  Expect(status).To(Equal(0))
  Expect(stderr).To(gbytes.Say("[1] curl https://www.google.ie/"))
  Expect(stdout).To(gbytes.Say("<title>Google</title>"))
```

All mocks have an *Exported counterpart that exports the mocking function for all the child processes.

More examples can be found in the [integration tests](./gob_test.go).