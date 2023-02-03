package gee

import (
	"fmt"
	"strings"
)

//实现动态路由，即一条路由规则可以匹配某一类型而非某一条固定的路由。
//例如/hello/:name，可以匹配/hello/geektutu、hello/jack等。
//HTTP请求的路径恰好是由/分隔的多段构成的，因此，每一段可以作为前缀树的一个节点。

//目前实现功能：
//参数匹配:。例如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc。
//通配*。例如 /static/*filepath，可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js，
//这种模式常用于静态服务器，能够递归地匹配子路径。

type node struct {
	pattern  string  //待匹配路由，例如 /p/:lang，是否一个完整的url，不是则为空字符串
	part     string  //URL块值，用/分割的部分，比如/abc/123，abc和123就是2个part
	children []*node //子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否模糊匹配，part含有 : 或 * 时为true
}

//为了实现动态路由匹配，加上了isWild这个参数。
//即当我们匹配 /p/go/doc/这个路由时，第一层节点，p精准匹配到了p，
//第二层节点，go模糊匹配到:lang，那么将会把lang这个参数赋值为go，继续下一层匹配

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

//为某个方法的trie树插入节点，一边匹配一边插入的方法
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height { //当前已填height个字段
		// 如果已经匹配完了，那么将pattern赋值给该node，表示它是一个完整的url，终止标志
		// 递归终止条件
		n.pattern = pattern
		return
	}

	part := parts[height]
	//寻找第一个匹配成功的节点
	child := n.matchChild(part)
	if child == nil {
		// 没有匹配上，那么进行生成，放到n节点的子列表中
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 接着插入下一个part节点
	child.insert(pattern, parts, height+1)
}

//查找某个方法的trie树，查询某个url返回所在节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 递归终止条件，找到末尾了或者通配符
		if n.pattern == "" {
			// pattern为空字符串表示它不是一个完整的url，匹配失败
			return nil
		}
		return n
	}

	part := parts[height]
	// 获取所有可能的子路径
	children := n.matchChildren(part)

	for _, child := range children {
		// 对于每条路径接着用下一part字段去查找
		result := child.search(parts, height+1)
		if result != nil {
			// 找到了即返回
			return result
		}
	}

	return nil
}

// 遍历路由树，查找所有完整的url，保存到切片中
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		// 一层一层的递归找pattern是非空的节点
		child.travel(list)
	}
}

//找到匹配的子节点，场景是用在插入时使用，找到1个匹配的就立即返回
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 返回所有匹配的子节点，原因是它的场景是用以查找
// 它必须返回所有可能的子节点来进行遍历查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
