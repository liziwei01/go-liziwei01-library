package env

// Default (全局)默认的环境信息
//
// 全局的 RootDir() 、DataDir() 等方法均使用该环境信息
var Default = New(Option{})

// Default 现在为AppEnv接口，现在开始实现接口要求的方法
// RootDir (全局)获取应用根目录
func RootDir() string {
	return Default.RootDir()
}

// DataDir (全局)设置应用数据根目录
func DataDir() string {
	return Default.DataDir()
}

// LogDir (全局)获取应用日志根目录
func LogDir() string {
	return Default.LogDir()
}

// ConfDir (全局)获取应用配置根目录
func ConfDir() string {
	return Default.ConfDir()
}

// AppName (全局)应用的名称
func AppName() string {
	return Default.AppName()
}

// RunMode (全局) 程序运行等级
// 默认是 release(线上发布)，还可选 RunModeDebug、RunModeTest
func RunMode() string {
	return Default.RunMode()
}

// Options 获取当前环境的选项详情
func Options() Option {
	return Default.Options()
}

// CloneWithOption 复制一个新的env对象，并将传入的Option merge进去
func CloneWithOption(opt Option) AppEnv {
	return Default.CloneWithOption(opt)
}
