package sco

import (
	"regexp"
	"strings"
)

// 删除前后的连续空格
// 如果失败返回空字符串
func RemoveLeadingAndTrailingSpace(line string) string {
	if line == "" {
		return ""
	}

	// 初始化正则
	reg, err := regexp.Compile(`((^\s+)|(\s+$))`)
	if err != nil {
		return ""
	}

	// 替换字符串
	line = reg.ReplaceAllString(line, "")
	if line == "" {
		return ""
	}

	return line
}

// 分割 key & value
func SplitKeyAndValue(line string) []string {
	if line == "" {
		return nil
	}

	var key, value string

	// 遍历字符串
	for index, char := range line {
		if char == '=' {
			key = RemoveLeadingAndTrailingSpace(line[:index])
			value = RemoveLeadingAndTrailingSpace(line[index+1:])
			break
		}
	}

	if key == "" || value == "" {
		return nil
	}

	keyAndValue := make([]string, 0)
	keyAndValue = append(keyAndValue, key, value)

	return keyAndValue
}

// 构造描述字符串
func MakeDescribeStr(describe, prefix string) string {
	if describe == "" {
		return ""
	}
	var describeStr string

	// 分割描述
	descs := strings.Split(describe, "\n")

	// 如果失败
	if len(descs) == 0 {
		return ""
	}

	// 遍历
	for _, desc := range descs {
		if desc != "" {
			describeStr += prefix + "#" + desc + "\n"
		}
	}

	return describeStr
}

// 提取块名
func ExtractSectionNameFromSectionNameStrline(sectionNameStrline string) string {
	// 初始化正则
	reg, err := regexp.Compile(`(\[\s*)|(\s*\])`)
	if err != nil {
		return ""
	}

	return reg.ReplaceAllString(sectionNameStrline, "")
}

// 创建一个块
func MakeSection() *_section {
	section := new(_section)
	if section.parameters == nil {
		section.parameters = make(_parameters)
	}
	return section
}
