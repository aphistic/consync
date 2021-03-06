package client

import (
	"path"

	"strings"

	"github.com/hashicorp/consul/api"
)

type ActionType int

const (
	ActionAdd ActionType = iota
	ActionModify
	ActionRemove
)

type kvItem struct {
	Path  string
	Value []byte
}

func getClient(addr *Address) (*api.Client, error) {
	return api.NewClient(&api.Config{
		Address:    addr.Addr,
		Scheme:     addr.Scheme,
		Datacenter: addr.DataCenter,
		Token:      addr.ACLToken,
	})
}

func getValues(addr *Address, recursive bool) (map[string]*kvItem, error) {
	client, err := getClient(addr)
	if err != nil {
		return nil, err
	}

	fKV := client.KV()
	fKVPairs, _, err := fKV.List(addr.Path, nil)
	if err != nil {
		return nil, err
	}

	values := make(map[string]*kvItem)
	folders := make(map[string]bool)
	for _, kvp := range fKVPairs {
		key := stripPrefix(kvp.Key, addr.Path)

		if key == "" {
			// Skip the root folder
			continue
		}

		// If we're not syncing recursively then just skip
		// any paths that include a folder
		if !recursive && strings.Contains(key, "/") {
			continue
		}

		// Keep track of any folders we find in the keys
		folders[path.Dir(key)] = true

		values[key] = &kvItem{
			Path:  key,
			Value: kvp.Value,
		}
	}

	// Remove any folders that have been added as values.
	// For some reason the api doesn't give any way to determine if
	// something is a folder or a value when you get it in the list.
	for folder := range folders {
		delete(values, folder)
		delete(values, folder+"/")
	}

	return values, nil
}

func fixPath(path string) string {
	if path != "" && path[0] == '/' {
		path = path[1:]
	}
	if path != "" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}

func stripPrefix(path, prefix string) string {
	newPath := path
	if newPath != "" && newPath[0] == '/' {
		newPath = newPath[1:]
	}
	if len(newPath) > len(prefix) {
		if newPath[:len(prefix)] == prefix {
			newPath = newPath[len(prefix):]
			if newPath != "" && newPath[0] == '/' {
				newPath = newPath[1:]
			}
		}
	}
	return newPath
}
