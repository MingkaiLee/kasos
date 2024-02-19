package util

import (
	"sort"
	"strings"
)

func ConvertTags(tags map[string]string) string {
	kvs := make([]string, 0, len(tags))
	for t := range tags {
		kvs = append(kvs, t)
	}
	sort.Strings(kvs)
	for i := range kvs {
		kvs[i] = kvs[i] + "=" + tags[kvs[i]]
	}
	tagsStr := strings.Join(kvs, "&")

	return tagsStr
}

func RevertTags(tagsStr string) map[string]string {
	tags := make(map[string]string)
	kvs := strings.Split(tagsStr, "&")
	for i := range kvs {
		kv := strings.Split(kvs[i], "=")
		if len(kv) != 2 {
			continue
		}
		tags[kv[0]] = kv[1]
	}

	return tags
}
