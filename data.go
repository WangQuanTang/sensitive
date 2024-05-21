// data.go
package sensitive

import "regexp"

// DataFilter 敏感词过滤器，用于一组[]string
type DataFilter struct {
	trie  *DataTrie
	noise *regexp.Regexp
}

// NewDataFilter 返回一个敏感词过滤器
func NewDataFilter() *DataFilter {
	return &DataFilter{
		trie:  NewDataTrie(),
		noise: regexp.MustCompile(`[\|\s&%$@*]+`),
	}
}

// UpdateNoisePattern 更新去噪模式
func (filter *DataFilter) UpdateNoisePattern(pattern string) {
	filter.noise = regexp.MustCompile(pattern)
}

// Filter 过滤敏感词
func (filter *DataFilter) Filter(text string) string {
	return filter.trie.Filter(text)
}

// Replace 和谐敏感词
func (filter *DataFilter) Replace(text string, repl rune) string {
	return filter.trie.Replace(text, repl)
}

// FindIn 检测敏感词
func (filter *DataFilter) FindIn(text string) (bool, string) {
	text = filter.RemoveNoise(text)
	return filter.trie.FindIn(text)
}

// FindAll 找到所有匹配词
func (filter *DataFilter) FindAll(text string) []string {
	return filter.trie.FindAll(text)
}

// Validate 检测字符串是否合法
func (filter *DataFilter) Validate(text string) (bool, string) {
	text = filter.RemoveNoise(text)
	return filter.trie.Validate(text)
}

// RemoveNoise 去除空格等噪音
func (filter *DataFilter) RemoveNoise(text string) string {
	return filter.noise.ReplaceAllString(text, "")
}

// FindInSlice 检测敏感词在 []string 中
func (filter *DataFilter) FindInSlice(texts []string) (bool, string) {
	for _, text := range texts {
		text = filter.RemoveNoise(text)
		if found, word := filter.trie.FindIn(text); found {
			return true, word
		}
	}
	return false, ""
}

// FindAllInSlice 找到所有匹配词在 []string 中
func (filter *DataFilter) FindAllInSlice(texts []string) []string {
	var matches []string
	for _, text := range texts {
		text = filter.RemoveNoise(text)
		matches = append(matches, filter.trie.FindAll(text)...)
	}
	return matches
}

// ValidateSlice 检测 []string 中的字符串是否合法
func (filter *DataFilter) ValidateSlice(texts []string) (bool, []string) {
	for _, text := range texts {
		text = filter.RemoveNoise(text)
		if valid, word := filter.trie.Validate(text); !valid {
			return false, []string{word}
		}
	}
	return true, nil
}
