package etcd

import (
	"sort"
	"strconv"
	"strings"
)

const (
	ns     = "namespaces"
	labels = "labels"
)

// keyspace namespace namespaces/user_id:namespace
// example namespaces/username:default
// example data TODO: Think about data... maybe link to resource in other service
// and by link i think key, to o delete when namespace is removed! And labes as well
/*
namespaces/labels/user_id:namespace -> [k:v, k:v]
namespaces/user_id:namespace -> {data}
*/

func newNSKeyspace(userid, namespace string) string {
	userNS := strings.Join([]string{userid, namespace}, ":")
	return strings.Join([]string{ns, userNS}, "/")
}

func newNSLabelsKeyspace(userid, namespace string) string {
	prefix := strings.Join([]string{ns, labels}, "/")
	userNS := strings.Join([]string{userid, namespace}, ":")
	s := []string{prefix, userNS}
	return strings.Join(s, "/")
}

func nsKeyspace(namespace string) string {
	return strings.Join([]string{ns, namespace}, "/")
}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func split(data string) []string {
	return delete_empty(strings.Split(data, ","))
}

func ssplit(data, sep string) []string {
	return delete_empty(strings.Split(data, sep))
}

func join(sep string, parts []string) string {
	return strings.Join(parts, sep)
}

func toString(n int64) string {
	return strconv.FormatInt(n, 10)
}

func NSLabels(userid string) string {
	return strings.Join([]string{ns, labels, userid}, "/")
}

func NSLabelsKey(name string) string {
	prefix := nsKeyspace(labels)
	s := []string{prefix, name}
	return strings.Join(s, "/")
}

func SplitLabels(value string) []string {
	ls := strings.Split(value, ",")
	sort.Strings(ls)
	return ls
}

func Compare(a, b []string, strict bool) bool {
	for _, akv := range a {
		for _, bkv := range b {
			if akv == bkv && !strict {
				return true
			}
		}
	}
	return true
}
