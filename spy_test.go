package gobmock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spy", func() {

	It("is created", func() {
		Expect(Spy("glass")).NotTo(BeNil())
	})

	It("includes a basic shell spy", func() {
		spy := Spy("chicken-with-a-pulley").MockContents()
		Expect(spy).To(MatchRegexp("^chicken-with-a-pulley\\(\\) {\\s+"))
		Expect(spy).To(MatchRegexp("echo.* > /dev/fd/2"))
	})

	It("includes a pipe consumer", func() {
		spy := Spy("grass").MockContents()
		Expect(spy).To(ContainSubstring("while read -t0.05; do"))
	})

	It("can be exported", func() {
		spy := ExportedSpy("wok").MockContents()
		Expect(spy).To(ContainSubstring("\nexport -f wok\n"))
	})
})
