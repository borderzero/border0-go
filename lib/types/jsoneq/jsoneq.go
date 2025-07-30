package jsoneq

import (
	"bytes"
	"encoding/json"
	"sort"
)

// options is an internal options object used
// for selective pruning of zeroed JSON fields.
type options struct {
	pruneEmptyObjects bool
	pruneEmptySlices  bool
	pruneEmptyStrings bool
}

// Option represents an option for testing JSON object equality.
type Option func(*options)

func PruneEmptyObjects() Option { return func(o *options) { o.pruneEmptyObjects = true } }
func PruneEmptySlices() Option  { return func(o *options) { o.pruneEmptySlices = true } }
func PruneEmptyStrings() Option { return func(o *options) { o.pruneEmptyStrings = true } }

// AreEqual compares two JSON strings for semantic equality, ignoring array ordering.
func AreEqual(a, b string, opts ...Option) bool {
	var ai, bi any
	if err := json.Unmarshal([]byte(a), &ai); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &bi); err != nil {
		return false
	}
	Prune(ai, opts...)
	Prune(bi, opts...)
	Canonicalize(ai)
	Canonicalize(bi)
	ab, err := json.Marshal(ai)
	if err != nil {
		return false
	}
	bb, err := json.Marshal(bi)
	if err != nil {
		return false
	}
	return bytes.Equal(ab, bb)
}

// Prune recursively removes default JSON values given some options.
func Prune(v any, opts ...Option) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	dropEmptyValues(v, o)
}

func dropEmptyValues(v any, opts *options) {
	switch x := v.(type) {
	case []any:
		for _, e := range x {
			dropEmptyValues(e, opts)
		}
	case map[string]any:
		for k, e := range x {
			switch val := e.(type) {
			case string:
				if !opts.pruneEmptyStrings {
					continue
				}

				if val == "" {
					delete(x, k)
					continue
				}
			case nil:
				if !opts.pruneEmptyObjects {
					continue
				}
				delete(x, k)
				continue
			case []any:
				if !opts.pruneEmptySlices {
					continue
				}
				if len(val) == 0 {
					delete(x, k)
					continue
				}
				dropEmptyValues(val, opts)
			case map[string]any:
				if !opts.pruneEmptyObjects {
					continue
				}
				dropEmptyValues(val, opts)
				if len(val) == 0 {
					delete(x, k)
				}
			default:
				dropEmptyValues(e, opts)
			}
		}
	}
}

// Canonicalize recursively sorts arrays within a JSON-like structure.
func Canonicalize(v any) {
	switch x := v.(type) {
	case []any:
		for i := range x {
			Canonicalize(x[i])
		}
		sort.Slice(x, func(i, j int) bool {
			bi, _ := json.Marshal(x[i])
			bj, _ := json.Marshal(x[j])
			return string(bi) < string(bj)
		})
	case map[string]any:
		for _, v := range x {
			Canonicalize(v)
		}
	}
}
