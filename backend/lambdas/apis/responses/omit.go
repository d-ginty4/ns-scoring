package responses

import (
	"encoding/json"
	"strings"
)

// omitLoadedRelationIDs removes a foreign key such as regionId when the matching
// relation object is present and has a non-zero id. The foreign key is kept when
// the relation was not loaded.
func omitLoadedRelationIDs(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	stripLoadedRelationIDs(out)
	return out, nil
}

func stripLoadedRelationIDs(v any) {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			stripLoadedRelationIDs(val)
			rel, ok := strings.CutSuffix(k, "Id")
			if !ok || rel == "" {
				continue
			}
			obj, ok := t[rel].(map[string]any)
			if !ok || obj["id"] == nil || obj["id"] == float64(0) {
				continue
			}
			delete(t, k)
		}
	case []any:
		for _, item := range t {
			stripLoadedRelationIDs(item)
		}
	}
}
