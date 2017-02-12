package main

import "net/url"

func fixupURL(srcURL *url.URL) *url.URL {
	updatedURL := *srcURL

	if updatedURL.Path == "" || updatedURL.Path[len(updatedURL.Path)-1] != '/' {
		updatedURL.Path += "/"
	}

	return &updatedURL
}
