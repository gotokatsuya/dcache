package dcache

import (
	"os"
	"sort"
)

type byModTimeAsc []os.FileInfo

func (s byModTimeAsc) Len() int {
	return len(s)
}

func (s byModTimeAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byModTimeAsc) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() < s[j].ModTime().UnixNano()
}

func SortFileInfosByModTimeAsc(in []os.FileInfo) []os.FileInfo {
	out := byModTimeAsc(in)
	sort.Sort(out)
	return out
}
