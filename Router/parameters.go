package Router

import "strings"

type Values struct {
	*node
	path     string
	wildcard int
}

func (vv Values) Index(name string) int {
	if i, has := vv.namesidx[name]; has {
		return i
	}
	return -1
}

func (vv Values) Get(name string) string {
	if i, has := vv.namesidx[name]; has {
		return vv.findParam(i)
	}
	return ""
}

func (vv Values) findParam(idx int) (param string) {
	curIndex := len(vv.names) - 1
	urlPath := vv.path
	pathLen := len(vv.path)
	_node := vv.node

	if _node.text == "*" {
		pathLen -= vv.wildcard
		if curIndex == idx {
			param = urlPath[pathLen:]
			return
		}
		curIndex--
		_node = _node.parent
	}

	for ; _node != nil; _node = _node.parent {
		if _node.text == ":" {

			ctn := strings.LastIndexByte(urlPath, '/')
			if ctn == -1 {
				break
			}

			pathLen = ctn + 1

			if curIndex == idx {
				param = urlPath[pathLen:]
				break
			}
			curIndex--
		} else {
			pathLen -= len(_node.text)
		}

		urlPath = urlPath[0:pathLen]
	}
	return
}
