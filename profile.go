package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const default_cred = "~/.aws/credentials"

func correctHomeDir(path string) (res string) {
	if strings.HasPrefix(path, "~") {
		home := os.Getenv("HOME")
		dir := filepath.Dir(path)
		dir = strings.Replace(dir, "~", home, 1)

		bname := filepath.Base(path)
		res = filepath.Join(dir, bname)
	} else {
		res = path
	}
	return
}

func load_file(file_name string) (ret []string, err error) {
	file_name = correctHomeDir(file_name)
	fp, err := os.Open(file_name)
	if err != nil {
		if file_name == correctHomeDir(default_cred) {
			return []string{""}, nil
		}
		return nil, err
	}
	defer fp.Close()

	sc := bufio.NewScanner(fp)
	d := make(map[string]bool)

	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "[") {
			line = strings.Replace(line, "[profile ", "", 1)
			line = strings.Replace(line, "[", "", 1)
			line = strings.Replace(line, "]", "", 1)
			if !d[line] {
				d[line] = true
				ret = append(ret, line)
			}
		}
	}
	return ret, nil
}
