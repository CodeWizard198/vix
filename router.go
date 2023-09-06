package vix

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	START = "*"
	ROD   = "/"
)

const (
	SWEAT  = ':'
	SYMBOL = '@'
)

// router 路由
type router struct {
	trees map[string]*node
}

// node 路由树
type node struct {
	route string
	path  string
	// 静态匹配node
	children map[string]*node
	// 通配符匹配node
	startChild *node
	// 参数匹配
	paramChild *node
	// 正则匹配
	regexpChild *node
	// 正则匹配解析器
	regexpResolver *regexp.Regexp
	handler        HandleFunc
}

func newRouter() *router {
	return &router{trees: make(map[string]*node)}
}

// addRouter 添加路由到路由树
func (r *router) addRouter(method, path string, handleFunc HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: ROD,
		}
		r.trees[method] = root
	}
	checkPath(method, path)
	if path == ROD {
		root.handler = handleFunc
		root.route = ROD
		return
	}
	path = strings.Trim(path, ROD)
	locs := strings.Split(path, ROD)
	for _, loc := range locs {
		children := root.searchOrCreate(loc)
		if children == nil {
			panic(fmt.Sprintf("路由[%s-%s]注册失败, error:已存在相同路径参数匹配或通配符匹配", method, path))
		}
		root = children
	}
	if root.handler != nil {
		panic(fmt.Sprintf("注册路由[%s-%s]失败,error:路由重复注册", method, path))
	}
	root.handler = handleFunc
	root.route = path
}

// checkRouter 检查路由树中是否有指定路由
// 优先考虑静态匹配，匹配不上考虑正则匹配，再考虑通参数匹配，最后考虑通配符匹配
func (r *router) checkRouter(method, path string) (bool, *matchInfo) {
	root, ok := r.trees[method]
	var param map[string]string
	if !ok {
		return false, nil
	}
	if path == ROD {
		return true, &matchInfo{pod: root}
	}
	locs := strings.Split(path[1:], ROD)
	for _, loc := range locs {
		nod, isChild, isParma := root.childOf(loc)
		if !isChild {
			return false, nil
		}
		if isParma {
			if param == nil {
				param = make(map[string]string)
			}
			param[nod.path[1:]] = loc
		}
		root = nod
	}
	return true, &matchInfo{pod: root, param: param}
}

func (n *node) childOf(loc string) (*node, bool, bool) {
	if n.children == nil {
		return n.makeDecision(loc)
	}
	nod, ok := n.children[loc]
	if !ok {
		return n.makeDecision(loc)
	}
	return nod, ok, false
}

// searchOrCreate 查询或者创建路由树子节点
func (n *node) searchOrCreate(loc string) *node {
	if loc == START {
		if n.paramChild != nil {
			return nil
		}
		n.startChild = &node{
			path: loc,
		}
		return n.startChild
	}

	switch loc[0] {
	case SYMBOL:
		regexpTemplate := loc[1:]
		compile, err := regexp.Compile(regexpTemplate)
		if err != nil {
			panic(err)
		}
		n.regexpResolver = compile
		n.regexpChild = &node{
			path: loc,
		}
		return n.regexpChild
	case SWEAT:
		if n.startChild != nil {
			return nil
		}
		n.paramChild = &node{
			path: loc,
		}
		return n.paramChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	nod, ok := n.children[loc]
	if !ok {
		nod = &node{
			path: loc,
		}
		n.children[loc] = nod
	}
	return nod
}

// makeDecision 决策
// 选择node
// 正则 > 参数 > 通配符
func (n *node) makeDecision(loc string) (*node, bool, bool) {
	if n.regexpChild != nil {
		return n.regexpChild, n.regexpResolver.MatchString(loc), false
	}
	if n.paramChild != nil {
		return n.paramChild, true, true
	}
	return n.startChild, n.startChild != nil, false
}
