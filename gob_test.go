package gobmock

import (
	"io"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/progrium/go-basher"
)

var (
	bash     *basher.Context
	stdout   *gbytes.Buffer
	stderr   *gbytes.Buffer
	bashPath string
)

var _ = Describe("Integration", func() {

	BeforeSuite(func() {
		extractBash()
	})

	AfterSuite(func() {
		os.Remove(bashPath)
	})

	BeforeEach(func() {
		bash, _ = basher.NewContext(bashPath, false)
		stdout = gbytes.NewBuffer()
		stderr = gbytes.NewBuffer()
		bash.Stdout = io.MultiWriter(GinkgoWriter, stdout)
		bash.Stderr = io.MultiWriter(GinkgoWriter, stderr)
		bash.SelfPath = "/bin/echo"
		bash.CopyEnv()
	})

	subShellTest := `main_test() {
			  local sub_shell
			  sub_shell="$(mktemp)"
			  trap "rm '${sub_shell}'" EXIT
			  echo "#!${BASH}" > "${sub_shell}"
			  echo 'echo "My child should bring home $(curl some://nonsense > /dev/null 2>&1; echo $?) bad grades"' > "${sub_shell}"
			  chmod +x "${sub_shell}"
			  "${sub_shell}"
			}`
	Context("Stub", func() {
		It("stubs executables", func() {
			gobs := []GoBMock{Stub("curl")}
			ApplyMocks(bash, gobs)
			status, err := bash.Run("curl", []string{"xyz://nothing.to.see.here"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("works pith pipefail", func() {
			sourceString(`set -o pipefail
				main_test() {
				  echo "Yay!" | curl abs://urdly.namedurl
				}`)

			gobs := []GoBMock{Stub("curl")}
			ApplyMocks(bash, gobs)
			status, err := bash.Run("main_test", []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("can be used in a child process", func() {
			sourceString(subShellTest)

			gobs := []GoBMock{ExportedStub("curl")}
			ApplyMocks(bash, gobs)
			status, err := bash.Run("main_test", []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("My child should bring home 0 bad grades"))
		})
	})

	Context("Spy", func() {
		It("stubs the executable", func() {
			gobs := []GoBMock{Spy("curl")}
			ApplyMocks(bash, gobs)
			status, err := bash.Run("curl", []string{"xyz://nothing.to.see.here"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("reports the arguments", func() {
			gobs := []GoBMock{Spy("curl")}
			ApplyMocks(bash, gobs)
			bash.Run("curl", []string{"xyz://nothing.to.see.here"})
			Expect(stderr).To(gbytes.Say("curl xyz://nothing.to.see.here"))
		})

		It("reports the number of the call", func() {
			gobs := []GoBMock{Spy("curl")}
			ApplyMocks(bash, gobs)
			sourceString("test_main() { curl abc; curl cfg; curl dbc; }")
			bash.Run("test_main", []string{""})
			Expect(stderr).To(gbytes.Say("\\[1\\] curl abc"))
			Expect(stderr).To(gbytes.Say("\\[2\\] curl cfg"))
			Expect(stderr).To(gbytes.Say("\\[3\\] curl dbc"))
		})

		It("reports the standard input", func() {
			sourceString(`
				test_main() {
				  printf "Waves travel with malaria.\nThe self has samadhi\n" | curl
				}`)
			gobs := []GoBMock{Spy("curl")}
			ApplyMocks(bash, gobs)
			bash.Run("test_main", []string{""})
			Expect(stderr).To(gbytes.Say("\\[1 received\\] input:\n"))
			Expect(stderr).To(gbytes.Say("Waves travel with malaria.\n"))
			Expect(stderr).To(gbytes.Say("The self has samadhi"))
			Expect(stderr).To(gbytes.Say("[1 end received]"))
		})
	})

	Context("Mock", func() {
		It("stubs the executable", func() {
			gobs := []GoBMock{Mock("curl", "")}
			ApplyMocks(bash, gobs)
			status, err := bash.Run("curl", []string{"zabi://daba.dooooo"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("can simulate the return code", func() {
			gobs := []GoBMock{Mock("curl", "return 1")}
			ApplyMocks(bash, gobs)
			status, err := bash.Run("curl", []string{"https://google.ie"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(1))
		})

		It("produces predefined output", func() {
			gobs := []GoBMock{Mock("curl", "echo 'Such much wow'")}
			ApplyMocks(bash, gobs)
			status, err := bash.Run("curl", []string{"boop://dodge.for.prezident"})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("Such much wow"))
		})

		It("can be exported to child processes", func() {
			sourceString(subShellTest)

			gobs := []GoBMock{ExportedMock("curl", "")}
			ApplyMocks(bash, gobs)
			status, err := bash.Run("main_test", []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
			Expect(stdout).To(gbytes.Say("My child should bring home 0 bad grades"))
		})
	})

})

func sourceString(script string) {
	bash.Source("", func(string) ([]byte, error) {
		return []byte(script), nil
	})
}

func extractBash() {
	bashDir, err := homedir.Expand("~/.basher")
	if err != nil {
		log.Fatal(err, "1")
	}

	bashPath = bashDir + "/bash"
	if _, err := os.Stat(bashPath); os.IsNotExist(err) {
		err = basher.RestoreAsset(bashDir, "bash")
		if err != nil {
			log.Fatal(err, "1")
		}
	}
}