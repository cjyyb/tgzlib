package tgzlib

import (
	"os"
)

type Rule struct {

}

type pattern struct {
	match func(rule, path string) error
	rule  string
}

func (r *Rule) Ignore(name string, fileInfo os.FileInfo) {

}
