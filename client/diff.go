package client

import (
	"path"
	"reflect"
)

type DiffItem struct {
	FromPath  string
	ToPath    string
	Type      ActionType
	FromValue []byte
	ToValue   []byte
}

func Diff(from *Address, to *Address) ([]*DiffItem, error) {
	from.fixupValues()
	to.fixupValues()

	fromVals, err := getValues(from)
	if err != nil {
		return nil, err
	}
	toVals, err := getValues(to)
	if err != nil {
		return nil, err
	}

	results := make([]*DiffItem, 0)
	for fPath, fVal := range fromVals {
		item := &DiffItem{}

		item.Type = ActionAdd
		item.FromPath = path.Join(from.Path, fPath)
		item.FromValue = fVal.Value
		item.ToPath = path.Join(to.Path, fPath)
		item.ToValue = fVal.Value

		if tVal, ok := toVals[fPath]; ok {
			if reflect.DeepEqual(fVal, tVal) {
				// The from and to values are the same, there's nothing
				// to change
				continue
			}
			item.Type = ActionModify

			item.ToPath = path.Join(to.Path, fPath)
			item.ToValue = tVal.Value
		}

		results = append(results, item)
	}
	for tPath, tVal := range toVals {
		item := &DiffItem{}

		if _, ok := fromVals[tPath]; !ok {
			item.Type = ActionRemove
			item.ToPath = path.Join(to.Path, tPath)
			item.ToValue = tVal.Value

			results = append(results, item)
		}
	}

	return results, nil
}
