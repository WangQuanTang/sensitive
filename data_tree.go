package sensitive

// DataTrie 短语组成的Trie树.
type DataTrie struct {
	Root *DataNode
}

// DataNode Trie树上的一个节点.
type DataNode struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*DataNode
}

// NewDataTrie 新建一棵Trie
func NewDataTrie() *DataTrie {
	return &DataTrie{
		Root: NewDataRootNode(0),
	}
}

// Add 添加若干个词
func (tree *DataTrie) Add(words ...string) {
	for _, word := range words {
		tree.add(word)
	}
}

func (tree *DataTrie) add(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewDataNode(r)
			current.Children[r] = newNode
			current = newNode
		}
		if position == len(runes)-1 {
			current.isPathEnd = true
		}
	}
}

func (tree *DataTrie) Del(words ...string) {
	for _, word := range words {
		tree.del(word)
	}
}

func (tree *DataTrie) del(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; !ok {
			return
		} else {
			current = next
		}

		if position == len(runes)-1 {
			current.SoftDel()
		}
	}
}

// Replace 词语替换
func (tree *DataTrie) Replace(text string, character rune) string {
	var (
		parent  = tree.Root
		current *DataNode
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			for i := left; i <= position; i++ {
				runes[i] = character
			}
		}

		parent = current
	}

	return string(runes)
}

// Filter 直接过滤掉字符串中的敏感词
func (tree *DataTrie) Filter(text string) string {
	var (
		parent      = tree.Root
		current     *DataNode
		left        = 0
		found       bool
		runes       = []rune(text)
		length      = len(runes)
		resultRunes = make([]rune, 0, length)
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			resultRunes = append(resultRunes, runes[left])
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() {
			left = position + 1
			parent = tree.Root
		} else {
			parent = current
		}

	}

	resultRunes = append(resultRunes, runes[left:]...)
	return string(resultRunes)
}

// Validate 验证字符串是否合法，如不合法则返回false和检测到的第一个敏感词
func (tree *DataTrie) Validate(text string) (bool, string) {
	const (
		Empty = ""
	)
	var (
		parent  = tree.Root
		current *DataNode
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			return false, string(runes[left : position+1])
		}

		parent = current
	}

	return true, Empty
}

// FindIn 判断text中是否含有词库中的词
func (tree *DataTrie) FindIn(text string) (bool, string) {
	validated, first := tree.Validate(text)
	return !validated, first
}

// FindAll 找有所有包含在词库中的词
func (tree *DataTrie) FindAll(text string) []string {
	var matches []string
	var (
		parent  = tree.Root
		current *DataNode
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			matches = append(matches, string(runes[left:position+1]))
		}

		if position == length-1 {
			parent = tree.Root
			position = left
			left++
			continue
		}

		parent = current
	}

	var i = 0
	if count := len(matches); count > 0 {
		set := make(map[string]struct{})
		for i < count {
			_, ok := set[matches[i]]
			if !ok {
				set[matches[i]] = struct{}{}
				i++
				continue
			}
			count--
			copy(matches[i:], matches[i+1:])
		}
		return matches[:count]
	}

	return nil
}

// NewDataNode 新建子节点
func NewDataNode(character rune) *DataNode {
	return &DataNode{
		Character: character,
		Children:  make(map[rune]*DataNode, 0),
	}
}

// NewDataRootNode 新建根节点
func NewDataRootNode(character rune) *DataNode {
	return &DataNode{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*DataNode, 0),
	}
}

// IsLeafNode 判断是否叶子节点
func (node *DataNode) IsLeafNode() bool {
	return len(node.Children) == 0
}

// IsRootNode 判断是否为根节点
func (node *DataNode) IsRootNode() bool {
	return node.isRootNode
}

// IsPathEnd 判断是否为某个路径的结束
func (node *DataNode) IsPathEnd() bool {
	return node.isPathEnd
}

// SoftDel 置软删除状态
func (node *DataNode) SoftDel() {
	node.isPathEnd = false
}
