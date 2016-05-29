package summer

import (
	"strings"
	"qiniupkg.com/x/log.v7"
)

type tagOption struct {
	auto   bool
	depend bool
	name   string
	path   string
	prefix string
}

func buildTagOptions(tag string) *tagOption {
	to := &tagOption{}
	if tag == "*" {
		to.depend = true
		to.auto = true
		return to
	}
	if len(tag) <= 1 {
		log.Error("bad tag :", tag)
		return to
	}
	if strings.Contains(tag, ".") {
		to.depend = false
		to.prefix = tag[:strings.Index(tag, ".")]
		to.path = tag[strings.Index(tag, ".") + 1:]
		return to
	}
	to.depend = true
	to.name = tag
	return to
}