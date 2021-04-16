package sco

import (
	"regexp"
	"strings"
)

// 删除前后的连续空格
func RemoveLeadingAndTrailingSpace(line string) (string, error) {
	if line == "" {
		return "", nil
	}

	// 初始化正则
	reg, err := regexp.Compile(`((^\s+)|(\s+$))`)
	if err != nil {
		return "", MakeError("初始化正则失败！", err.Error())
	}

	// 替换字符串
	line = reg.ReplaceAllString(line, "")
	if line == "" {
		return "", MakeError("删除前后空格失败！", "")
	}

	return line, nil
}

// 删除第一个等号的前后连续空格
func RemoveFirstEqualsignLeadingAndTrailingSpace(line string) (string, error) {
	if line == "" {
		return "", nil
	}
	// 初始化正则
	reg, err := regexp.Compile(`((\s+=\s*)|(\s*=\s+))`)
	if err != nil {
		return "", MakeError("初始化正则失败！", err.Error())
	}

	keyAndValue := reg.Split(line, 2)
	if len(keyAndValue) != 2 {
		return "", MakeError("分割失败！", "")
	}

	return keyAndValue[0] + "=" + keyAndValue[1], nil
}

// 构造描述字符串
func MakeDescribeStr(describe string) string {
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
			describeStr += "#" + desc + "\n"
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

// 分离属性名和属性值
func ExtractParamNameAndValue(keyAndValueStr string) ([]string, error) {
	keyAndValueStr, err := RemoveFirstEqualsignLeadingAndTrailingSpace(keyAndValueStr)
	if err != nil {
		return nil, MakeError("删除第一个等号前后空格失败!", err.Error())
	}

	// 分割
	keyAndValue := strings.SplitN(keyAndValueStr, "=", 2)
	if len(keyAndValue) != 2 {
		return nil, MakeError("分割属性名和属性值失败！", "")
	}

	return keyAndValue, nil
}
