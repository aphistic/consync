package client

import (
	"reflect"

	"path"

	"github.com/hashicorp/consul/api"
)

func Sync(from *Address, to *Address) error {
	from.Path = fixPath(from.Path)
	to.Path = fixPath(to.Path)

	items, err := SyncPreview(from, to)
	if err != nil {
		return err
	}

	tClient, err := getClient(to)
	tKV := tClient.KV()
	for _, item := range items {
		switch item.Type {
		case ActionAdd:
			_, err = tKV.Put(&api.KVPair{
				Key:   item.Path,
				Value: item.Value,
			}, nil)
			if err != nil {
				return err
			}
		case ActionModify:
			_, err = tKV.Put(&api.KVPair{
				Key:   item.Path,
				Value: item.Value,
			}, nil)
			if err != nil {
				return err
			}
		case ActionRemove:
			_, err = tKV.Delete(item.Path, nil)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

type SyncPreviewItem struct {
	Path  string
	Value []byte
	Type  ActionType
}

func SyncPreview(from *Address, to *Address) ([]*SyncPreviewItem, error) {
	from.Path = fixPath(from.Path)
	to.Path = fixPath(to.Path)

	fromVals, err := getValues(from)
	if err != nil {
		return nil, err
	}
	toVals, err := getValues(to)
	if err != nil {
		return nil, err
	}

	results := make([]*SyncPreviewItem, 0)
	for fPath, fVal := range fromVals {
		item := &SyncPreviewItem{
			Type:  ActionAdd,
			Path:  path.Join(to.Path, fPath),
			Value: fVal.Value,
		}

		if tVal, ok := toVals[fPath]; ok {
			if reflect.DeepEqual(fVal, tVal) {
				continue
			}
			item.Type = ActionModify
			item.Value = fVal.Value
		}

		results = append(results, item)
	}
	for tPath := range toVals {
		if _, ok := fromVals[tPath]; !ok {
			item := &SyncPreviewItem{
				Type: ActionRemove,
				Path: path.Join(to.Path, tPath),
			}
			results = append(results, item)
		}
	}

	return results, nil
}
