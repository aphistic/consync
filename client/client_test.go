package client

import (
	"testing"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	sweet.T(func(s *sweet.S) {
		RegisterFailHandler(sweet.GomegaFail)

		s.RunSuite(t, &ClientSuite{})
	})
}

type ClientSuite struct{}

func (s *ClientSuite) TestStripPrefix(t *testing.T) {
	Expect(stripPrefix("/", "")).To(Equal(""))
	Expect(stripPrefix("/testing", "")).To(Equal("testing"))
	Expect(stripPrefix("/testing/a/long/path/", "")).To(Equal("testing/a/long/path/"))
	Expect(stripPrefix("/testing/a/long/path/", "testing")).To(Equal("a/long/path/"))
}
