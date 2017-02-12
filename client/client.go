package client

import "github.com/hashicorp/consul/api"

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
		Address: addr.Addr,
	})
}

func getValues(addr *Address) (map[string]*kvItem, error) {
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
	for _, kvp := range fKVPairs {
		values[kvp.Key] = &kvItem{
			Path:  kvp.Key,
			Value: kvp.Value,
		}
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
