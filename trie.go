package jim

import "strings"

//
//  node
//  @Description: 前缀树的结点
//
type node struct {
	pattern  string  // 待匹配的路由 例如 /p/:lang
	part     string  // 路由中的一部分 例如 :lang
	children []*node // 子结点 例如 /p/doc /p/intro /p/tutorial [doc,intro,tutorial]
	isWild   bool    // 是否精确匹配 part 中是否含有 : 或者 * 时为 true
}

//
//  matchChild
//  @Description: 第一个匹配成功的结点 用来插入
//  @receiver n
//  @param part
//  @return *node
//
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//
//  matchChildren
//  @Description: 所有匹配成功的结点 用于查找
//  @receiver n
//  @param part
//  @return []*node
//
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nil
}

// 对于路由来说，最重要的当然是注册与匹配了。
//开发服务时，注册路由规则，映射handler；
//访问时，匹配路由规则，查找到对应的handler。
//因此，Trie 树需要支持节点的插入与查询。

//
//  insert
//  @Description: //插入功能很简单，递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个，
//	有一点需要注意，/p/:lang/doc只有在第三层节点，即doc节点，pattern才会设置为/p/:lang/doc。
//	p和:lang节点的pattern属性皆为空。
//	因此，当匹配结束时，我们可以使用n.pattern == ""来判断路由规则是否匹配成功。
//	例如，/p/python虽能成功匹配到:lang，但:lang的pattern值为空，因此匹配失败。
//  @receiver n
//  @param pattern
//  @param parts
//  @param height
//
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果路由的层数和当前的层数一样才会有pattern
	// /p/:lang/doc  parts 为[p,:lang,doc] height =3 才会给pattern设置
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 查找每一层的结点
	part := parts[height]
	child := n.matchChild(part)
	// 如果没有的话就新建一个
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
	}
	// 否则直接插入
	child.insert(pattern, parts, height+1)
}

//
//  search
//  @Description: 查询node
// 	同样也是递归查询每一层的节点，退出规则是，匹配到了*，匹配失败，
//	或者匹配到了第len(parts)层节点。
//  @receiver n
//  @param parts
//  @param height
//  @return *node
//
func (n *node) search(parts []string, height int) *node {
	// 如果匹配到了第 len(parts)层 或者匹配到 *
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果 pattern 没有 就是匹配失败了
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// 匹配到所有的孩子
	children := n.matchChildren(part)

	for _, child := range children {
		// 如果能找到 node  就返回 node
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
