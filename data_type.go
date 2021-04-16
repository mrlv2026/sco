package sco

import "os"

// 参数值
type _parametersValue struct {
	// 描述
	describe string
	// 是否被禁用
	isDisable bool
	// 参数值
	value string
}

// 参数集合
type _parameters map[string]_parametersValue

// 块
type _section struct {
	// 块名
	name string
	// 描述
	describe string
	// 参数集合
	parameters _parameters
}

// 配置
type _config struct {
	// 描述
	describe string
	// 块的集合
	sections []_section
	// 打开的时候是否上锁
	isLock bool
	// 文件指针
	configFile *os.File
}
