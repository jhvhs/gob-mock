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
		Expect(spy).To(ContainSubstring("while read -r -t0.1; do"))
	})

	It("does not call through by default", func() {
		spy := Spy("bee").MockContents()
		Expect(spy).NotTo(ContainSubstring("$(which bee)"))
	})

	It("can call through to the executable", func() {
		spy := SpyAndCallThrough("squash").MockContents()
		Expect(spy).To(ContainSubstring("$(which squash)"))
	})

	It("can include a condition for a call through", func() {
		spy := SpyAndConditionallyCallThrough("raspbery", "[[ 3 -ne 4 ]]")
		Expect(spy.MockContents()).To(ContainSubstring("if [[ 3 -ne 4 ]]; then"))
	})
})

