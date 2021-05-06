package conf

// Default 默认的实例,全局的Parse、ParseBytes、Exists等方法均使用该对象
var Default = NewDefault(nil)

// Parse 解析配置，配置文件默认认为在 conf/目录下
//
// 	如配置文件 conf/abc.toml ，则读取时使用 Parse("abc.toml",&xxx)
// 	推荐使用上述相对文件名老读取配置，这样可通过修改全局应用环境信息的env.ConfDir，来调整配置目录
//
// 	也支持传入一个绝对路径 或者 相对路径
// 	如  /tmp/test.toml 或者  ./conf/test.toml
func Parse(confName string, obj interface{}) (err error) {
	return Default.Parse(confName, obj)
}

// ParseBytes 解析bytes
//
// fileExt 是file extension 文件后缀，如.json、.toml
func ParseBytes(fileExt string, content []byte, obj interface{}) error {
	return Default.ParseBytes(fileExt, content, obj)
}

// Exists  判断是否存在
func Exists(confName string) bool {
	return Default.Exists(confName)
}

// RegisterParserFunc 注册一个解析器
func RegisterParserFunc(fileExt string, fn ParserFunc) error {
	return Default.RegisterParserFunc(fileExt, fn)
}

// RegisterBeforeFunc 注册一个在解析前执行辅助回调方法
//
// name 唯一的名字；fn 回调函数
func RegisterBeforeFunc(name string, fn BeforeFunc) error {
	return Default.RegisterBeforeFunc(name, fn)
}