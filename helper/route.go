package helper

import (
	"strings"
)

func SplitSubpaths(uri string) []string {
	paths := strings.SplitAfter(uri, "/")
	var subpaths []string
	for i, path := range paths {
		var subpath string
		if i != 0 {
			subpath = subpaths[i-1]+"/"+strings.TrimSuffix(path, "/")
		} else {
			subpath = ""
		}
		subpaths = append(subpaths, subpath)
	}
	return subpaths[1:]
}
