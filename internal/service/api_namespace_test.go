package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

func TestTree(t *testing.T) {
	var s = []string{
		"/m/a",
		"/m/b",
		"/m/c",
		"/m/a/z",
		"/m/a/z/b",
		"/m/a/z/d",
		"/m/a/x/y",
		"/m/a/x",
		"/m/a/x/aa",
		"/m/b/z/b",
		"/m/b/z/d",
		"/m/b/x/y",
		"/m/b/x/aa",
		"/m/c/z/b",
		"/m/c/z/d",
		"/m/c/x/y",
		"/m/c/x/aa",
		"/d",
		"/z",
	}
	sort.Strings(s)
	b, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(b))
}

func TestMakeTree(t *testing.T) {
	list := &models.APINamespaceList{
		List: []*models.APINamespace{
			//{FullPath: "/a"},
			{FullPath: "/a/c"},
			{FullPath: "/a/c/c1"},
			{FullPath: "/a/c/c2"},
			{FullPath: "/a/b"},
			{FullPath: "/a/b/b1"},
			{FullPath: "/a/b/b2"},
		},
	}
	for _, v := range list.List {
		v.Parent, v.Namespace = apipath.Split(v.FullPath)
	}
	tree, err := newTreeMaker(list).makeTree()
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(b))
}
