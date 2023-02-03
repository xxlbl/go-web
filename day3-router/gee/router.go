package gee

import (
	"net/http"
	"strings"
)

//使用 roots 来存储每种请求方式的Trie树根节点。
//使用 handlers 存储每种请求方式的 HandlerFunc
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
// 分解地址->字段
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	// 查询or创建方法的trie树
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	//向trie树插入节点
	r.roots[method].insert(pattern, parts, 0)
	//建立映射
	r.handlers[key] = handler
}

//查询URL地址是否在trie树，返回终止节点和解析参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 查询的url地址分解
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	// 有无此方法的trie树
	if !ok {
		return nil, nil
	}
	// 在trie树中查找url，返回终止节点
	n := root.search(searchParts, 0)

	// 存在该url
	if n != nil {
		//分解路由树中存储的url地址
		parts := parsePattern(n.pattern)

		for index, part := range parts {
			//通配符解析 + 映射 trie树中url -> 实际查询url
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
		//返回终止节点,解析参数
	}

	return nil, nil
}

// 查找所有完整的url，保存其终止节点到切片中
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	//是否在路由树中存在，拿到解析参数
	if n != nil {
		// 上下文获得url解析参数
		c.Params = params
		key := c.Method + "-" + n.pattern
		// 调用映射对应的处理函数
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
