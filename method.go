package sco

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

/********** _config **********/

// 打开一个配置文件
func Open(cfg string, isLock bool) (*_config, error) {
	// 尝试打开
	file, err := os.OpenFile(cfg, os.O_RDWR, 0)
	if err != nil {
		MakeError("打开文件失败！", err.Error())
	}

	// 如果需要上锁
	if isLock {
		fmt.Print()
	}

	// 容器
	config := new(_config)

	// 文件
	config.configFile = file

	// 解析文件
	if err := config.parse(); err != nil {
		// 关闭文件
		config.configFile.Close()

		return nil, MakeError("解析文件失败！", err.Error())
	}

	// 成功
	return config, nil
}

// 关闭一个配置文件
func (config *_config) Close() error {
	var configText string
	var err error

	// 尝试格式化
	if configText, err = config.format(); configText == "" && err != nil {
		return MakeError("格式化失败！", err.Error())
	}

	// 清空文件
	if err = config.configFile.Truncate(0); err != nil {
		return MakeError("清空文件失败！", err.Error())
	}

	byteData := []byte(configText)

	// 把内容写入
	if l, err := config.configFile.WriteAt(byteData, 0); l != len(byteData) && err != nil {
		return MakeError("写入内容失败！", err.Error())
	}

	// 关闭文件
	config.configFile.Close()

	return nil
}

// 解析配置文件
func (config *_config) parse() error {
	// 创建读取对象
	read := bufio.NewReader(config.configFile)
	if read == nil {
		return MakeError("创建 Reader 对象失败！", "")
	}

	// 容器
	var describe string
	var lastSectionPointer *_section

	for {
		// 读取一行
		byteLine, _, err := read.ReadLine()
		if err != nil {
			if err == io.EOF {
				// 读取完毕
				break
			} else {
				// 读取出错
				return MakeError("读取出错！", err.Error())
			}
		}

		// 去除头尾空格
		line := RemoveLeadingAndTrailingSpace(string(byteLine))
		if line == "" {
			// 判断是否是配置文件的注释
			if lastSectionPointer == nil {
				config.describe = describe
				describe = ""
			}
			continue
		}

		switch line[0] {
		// 注释
		case '#':
			if describe == "" {
				describe += line[1:]
			} else {
				describe += "\n" + line[1:]
			}

		// 块名
		case '[':
			// 如果不是第一个块则把上一个块添加
			if lastSectionPointer != nil {
				config.AddSection(*lastSectionPointer)
				lastSectionPointer = MakeSection()
			}

			// 块名
			sectionName := ExtractSectionNameFromSectionNameStrline(line)
			if sectionName == "" {
				continue
			}

			lastSectionPointer = MakeSection()
			lastSectionPointer.SetSectionName(sectionName)
			lastSectionPointer.SetSectionDesc(describe)
			describe = ""

		// 被禁用的属性
		case ';':
			// 分割 key & value
			keyAndValue := SplitKeyAndValue(line[1:])
			if len(keyAndValue) != 2 {
				continue
			}

			var parameterValue _parameterValue
			parameterValue.isDisable = true
			parameterValue.value = keyAndValue[1]
			parameterValue.describe = describe

			// 添加参数
			lastSectionPointer.AddParamete(keyAndValue[0], parameterValue)

			describe = ""

		// 否则就是属性
		default:
			// 分割 key & value
			keyAndValue := SplitKeyAndValue(line[1:])
			if len(keyAndValue) != 2 {
				continue
			}

			var parameterValue _parameterValue
			parameterValue.isDisable = false
			parameterValue.value = keyAndValue[1]
			parameterValue.describe = describe

			// 添加参数
			lastSectionPointer.AddParamete(keyAndValue[0], parameterValue)

			describe = ""
		}
	}

	return nil
}

// 格式化配置
func (config *_config) format() (string, error) {
	// 容器
	var configText string

	// 添加配置注释
	if config.describe != "" {
		desc := MakeDescribeStr(config.describe)
		if desc != "" {
			configText += desc + "\n"
		}
	}

	// 添加块
	for _, section := range config.sections {
		// 容器
		var sectionText string

		// 添加块的描述
		if section.describe != "" {
			desc := MakeDescribeStr(section.describe)
			if desc != "" {
				sectionText += desc
			}
		}

		// 添加块的名称
		sectionText += "[" + section.name + "]\n"

		// 添加属性
		for name, value := range section.parameters {
			if value.isDisable {
				sectionText += fmt.Sprintf("%s;\t%s = %s\n", MakeDescribeStr(value.describe), name, value.value)
			} else {
				sectionText += fmt.Sprintf("%s\t%s = %s\n", MakeDescribeStr(value.describe), name, value.value)
			}
		}

		// 多加一个空行
		configText += sectionText + "\n"
	}

	return configText, nil
}

// 获取块属性
func (config *_config) GetConfigDesc() (string, bool) {
	if config.describe == "" {
		return "", false
	}
	return config.describe, true
}

// 设置块属性
func (config *_config) SetConfigDesc(newDesc string) {
	config.describe = newDesc
}

// 获取块的列表
func (config *_config) GetSectionList() ([]string, error) {
	sectionList := make([]string, 0)

	for _, section := range config.sections {
		sectionList = append(sectionList, section.name)
	}

	return sectionList, nil
}

// 判断块是否存在
func (config *_config) SectionIsExist(name string) bool {
	for _, section := range config.sections {
		if section.name == name {
			return true
		}
	}

	return false
}

// 获取块指针
func (config *_config) GetSectionPointer(name string) (*_section, error) {
	for i, section := range config.sections {
		if section.name == name {
			return &config.sections[i], nil
		}
	}

	return nil, MakeError("块不存在！", "")
}

// 添加块
func (config *_config) AddSection(section _section) error {
	// 判断是否存在
	if config.SectionIsExist(section.name) {
		return MakeError("块已经存在！", "")
	}

	// 添加
	config.sections = append(config.sections, section)
	if !config.SectionIsExist(section.name) {
		return MakeError("添加块失败！", "")
	}
	return nil
}

// 删除块
func (config *_config) DeleteSection(name string) error {
	for index, section := range config.sections {
		if section.name == name {
			config.sections = append(config.sections[:index], config.sections[index+1:]...)
			return nil
		}
	}
	return nil
}

/********** _section **********/

// 获取块名
func (section *_section) GetSectionName() string {
	return section.name
}

// 设置块名
func (section *_section) SetSectionName(newName string) {
	section.name = newName
}

// 获取块属性
func (section *_section) GetSectionDesc() (string, bool) {
	if section.describe == "" {
		return "", false
	}
	return section.describe, true
}

// 设置块属性
func (section *_section) SetSectionDesc(newDesc string) {
	section.describe = newDesc
}

// 获取块的所有参数名
func (section *_section) GetSectionParameterList() ([]string, error) {
	parameterList := make([]string, 0)

	for name := range section.parameters {
		parameterList = append(parameterList, name)
	}

	return parameterList, nil
}

// 禁用属性
func (section *_section) Disable(name string) {
	// 如果不存在
	if !section.ParameterIsExist(name) {
		return
	}

	// 禁用
	temp := section.parameters[name]
	temp.isDisable = true
	section.parameters[name] = temp
}

// 启用属性
func (section *_section) Enable(name string) {
	// 如果不存在
	if !section.ParameterIsExist(name) {
		return
	}

	// 启用
	temp := section.parameters[name]
	temp.isDisable = false
	section.parameters[name] = temp
}

// 判断参数是否被禁用
// 第一个返回参数表示参数是否存在，第二个表示是否被禁用
func (section *_section) IsDisable(name string) (bool, bool) {
	if !section.ParameterIsExist(name) {
		return false, false
	}
	return true, section.parameters[name].isDisable
}

// 设置参数描述
// 返回表示参数是否存在
func (section *_section) SetParameterDesc(name, newDesc string) bool {
	if !section.ParameterIsExist(name) {
		return false
	}

	temp := section.parameters[name]
	temp.describe = newDesc
	section.parameters[name] = temp
	return true
}

// 获取参数描述
// 第一个是描述，第二个参数是否存在
func (section *_section) GetParameterDesc(name string) (string, bool) {
	if !section.ParameterIsExist(name) {
		return "", false
	}
	return section.parameters[name].describe, false
}

// 设置参数值
// 表示参数是否存在
func (section *_section) SetParameterValue(name, newValue string) bool {
	if !section.ParameterIsExist(name) {
		return false
	}

	temp := section.parameters[name]
	temp.value = newValue
	section.parameters[name] = temp
	return true
}

// 获取参数值
// 第一个是值，第二个是参数是否存在
func (section *_section) GetParameterValue(name string) (string, bool) {
	if !section.ParameterIsExist(name) {
		return "", false
	}
	return section.parameters[name].value, false
}

// 判断参数是否存在
func (section *_section) ParameterIsExist(name string) bool {
	for ParamName := range section.parameters {
		if ParamName == name {
			return true
		}
	}
	return false
}

// 添加属性
func (section *_section) AddParamete(name string, value _parameterValue) error {
	// 如果存在就报错
	if !section.ParameterIsExist(name) {
		return MakeError("属性已经存在！", "")
	}

	// 添加
	section.parameters[name] = value

	return nil
}

// 删除属性
func (section *_section) DeleteParamete(name string) error {
	// 如果存在就报错
	if !section.ParameterIsExist(name) {
		return nil
	}

	// 删除
	delete(section.parameters, name)

	if section.ParameterIsExist(name) {
		return MakeError("属性删除失败！", "")
	}

	return nil
}
