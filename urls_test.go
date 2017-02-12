package main

import (
	"testing"

	"net/url"

	. "github.com/onsi/gomega"
)

type URLSuite struct{}

func (s *URLSuite) TestStripTrailingSlash(t *testing.T) {
	startURL, err := url.Parse("http://10.10.10.10:1234")
	Expect(err).To(BeNil())
	Expect(startURL.Path).To(Equal(""))
	updatedURL := fixupURL(startURL)
	Expect(updatedURL.Path).To(Equal("/"))

	startURL, err = url.Parse("http://10.10.10.10:1234/some/long/path")
	Expect(err).To(BeNil())
	Expect(startURL.Path).To(Equal("/some/long/path"))
	updatedURL = fixupURL(startURL)
	Expect(updatedURL.Path).To(Equal("/some/long/path/"))
}
