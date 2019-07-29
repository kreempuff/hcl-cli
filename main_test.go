package main

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/hcl"
	"io"
	"os"
	"strings"
	"testing"
)

type FileHolder struct {
	json []byte
	hcl  []byte
}

func Test_toJson(t *testing.T) {
	dir, err := os.Open("./tests")
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	infos, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}

	files := map[string]FileHolder{}

	for _, info := range infos {
		split := strings.Split(info.Name(), ".")
		name, ext := split[0], split[1]
		var holder FileHolder
		if h, exists := files[name]; exists {
			holder = h
		}

		var b = &bytes.Buffer{}

		f, err := os.Open(dir.Name() + "/" + info.Name())
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(b, f)
		if err != nil {
			f.Close()
			panic(err)
		}
		f.Close()

		switch ext {
		case "json":
			holder.json = b.Bytes()
		case "hcl":
			holder.hcl = b.Bytes()
		default:
			panic(fmt.Sprintf("unrecognized file type: %s", info.Name()))
		}

		files[name] = holder
	}

	for k, v := range files {
		t.Run(k, func(t *testing.T) {
			tree, err := hcl.ParseBytes(v.hcl)
			if err != nil {
				t.Errorf("hcl.ParseBytes() error = %v", err)
				return
			}

			got, err := toJson(tree)
			if err != nil {
				t.Errorf("toJson() error = %v", err)
				return
			}
			if got != string(v.json) {
				t.Errorf("toJson() got = %v, want %v", got, string(v.json))
			}
		})
	}

	//type args struct {
	//	astFile *ast.File
	//}
	//tests := []struct {
	//	name    string
	//	args    args
	//	want    string
	//	wantErr bool
	//}{
	//
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		got, err := toJson(tt.args.astFile)
	//		if (err != nil) != tt.wantErr {
	//			t.Errorf("toJson() error = %v, wantErr %v", err, tt.wantErr)
	//			return
	//		}
	//		if got != tt.want {
	//			t.Errorf("toJson() got = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}
