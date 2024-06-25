package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var stats = make(map[rune]int)

func isInclude(filePath string, extsString *string) bool {
	fileExt := path.Ext(filePath)
	if fileExt == "" {
		return false
	}

	// remove dot
	fileExt = fileExt[1:]

	exts := strings.Split(*extsString, ",")

	for _, v := range exts {
		if fileExt == v {
			return true
		}
	}

	return false
}

func fileStat(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		stats[c]++
	}
}

func main() {
	var xFlag = flag.String(
		"exts",
		"go,mod",
		"File extensions to include; separated by commas")

	flag.Parse()

	root := flag.Args()

	if len(root) == 0 {
		root = append(root, ".")
	}

	for _, v := range root {
		filepath.WalkDir(v, fs.WalkDirFunc(func(filePath string, info fs.DirEntry, err error) error {
			if err != nil {
				panic(err)
			}

			if !info.IsDir() && isInclude(filePath, xFlag) {
				fileStat(filePath)
			}

			return nil
		}))
	}

	FormatMap(stats)
}

func FormatMap(m map[rune]int) {
	n := map[int][]rune{}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, k := range a {
		for _, s := range n[k] {
			fmt.Printf("%q %d\n", s, k)
		}
	}
}
