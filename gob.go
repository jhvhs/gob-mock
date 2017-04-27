package gobmock

import "github.com/progrium/go-basher"

type GoBMock interface {
	MockContents() string
}

type BasherContext interface {
	Source(string, func(string) ([]byte, error)) error
}

func ApplyMocks(bash *basher.Context, mocks []GoBMock) {
	bash.Source("", func(string) ([]byte, error) {
		return []byte("export callCounter=0"), nil
	})
	for _, mock := range mocks {
		bash.Source("", func(string) ([]byte, error) {
			return []byte(mock.MockContents()), nil
		})
	}
}
