package summer

import (
	"strings"
)

type tagOption struct {
	auto   bool
	depend bool
	name   string
	path   string
	prefix string
}

func parseTagOption(tag string) *tagOption {
	option := &tagOption{}
	if tag == "*" {
		option.depend = true
		option.auto = true
		return option
	}
	if len(tag) <= 1 {
		logger.Error("bad tag :", tag)
		return option
	}
	if strings.Contains(tag, ".") {
		option.depend = false
		option.prefix = tag[:strings.Index(tag, ".")]
		option.path = tag[strings.Index(tag, ".")+1:]
		return option
	}
	option.depend = true
	option.name = tag
	return option
}
