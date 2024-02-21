package util

import (
	"fmt"
	"sort"
	"strings"
)

// convert tags map to PromQL format
func ConvertTags(tags map[string]string) string {
	kvs := make([]string, 0, len(tags))
	for t := range tags {
		kvs = append(kvs, t)
	}
	sort.Strings(kvs)
	for i := range kvs {
		kvs[i] = fmt.Sprintf(`%s="%s"`, kvs[i], tags[kvs[i]])
	}
	tagsStr := strings.Join(kvs, ",")

	return tagsStr
}

// revert tags map from PromQL format
func RevertTags(tagsStr string) map[string]string {
	tags := make(map[string]string)
	kvs := strings.Split(tagsStr, ",")
	for i := range kvs {
		kv := strings.Split(kvs[i], "=")
		if len(kv) != 2 {
			continue
		}
		tags[kv[0]] = strings.Trim(kv[1], `"`)
	}

	return tags
}
