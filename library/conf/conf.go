package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go-liziwei01-library/library/env"
)

// Conf 配置解析定义接口
type Conf interface {
	// 读取并解析配置文件
	// confName 支持相对路径和绝对路径
	Parse(confName string, obj interface{}) error
	// 解析bytes内容
	ParseBytes(fileExt string, content []byte, obj interface{}) error
	// 配置文件是否存在
	Exists(confName string) bool
	// 注册一个指定后缀的配置的parser
	// 如要添加 .ini 文件的支持，可在此注册对应的解析函数即可
	RegisterParserFunc(fileExt string, fn ParserFunc) error
	// 注册一个在解析前执行辅助回调方法
	// 先注册的先执行，不能重复
	RegisterBeforeFunc(name string, fn BeforeFunc) error
	// 配置的环境信息
	Env() env.AppEnv
}

// New 创建一个新的配置解析实例
// 返回的实例是没有注册任何解析能力的
func New(e env.AppEnv) Conf {
	conf := &conf{
		parsers: map[string]ParserFunc{},
		env:     e,
	}
	return conf
}

// NewDefault 创建一个新的配置解析实例
// 会注册默认的配置解析方法和辅助方法
func NewDefault(e env.AppEnv) Conf {
	conf := New(e)
	for name, fn := range DefaultParserFuncs {
		_ = conf.RegisterParserFunc(name, fn)
	}
	for _, h := range defaultHelpers {
		if err := conf.RegisterBeforeFunc(h.Name, h.Func); err != nil {
			panic(fmt.Sprintf("RegisterHelper(%q) err=%s", h.Name, err))
		}
	}
	return conf
}

// conf实现Conf接口
type conf struct {
	env     env.AppEnv
	parsers map[string]ParserFunc
	helpers []*beforeHelper
}

// 传入文件名和接收obj实现解析配置文件，配置文件默认认为在 conf/ 目录下
func (c *conf) Parse(confName string, obj interface{}) (err error) {
	confAbsPath := c.confFileRealPath(confName)
	return c.parseByAbsPath(confAbsPath, obj)
}

var relPathPre = "." + string(filepath.Separator)

// 将文件名组装为文件实际所在目录
func (c *conf) confFileRealPath(confName string) string {
	// 若文件名已经是绝对路径或以./开头或以./开头，视为找到了绝对路径
	if filepath.IsAbs(confName) ||
		strings.HasPrefix(confName, "./") ||
		strings.HasPrefix(confName, relPathPre) {
		return confName
	}
	// 将文件名加上环境变量里面的confdir的前缀
	return filepath.Join(c.Env().ConfDir(), confName)
}

// 通过绝对路径找到文件，确保文件名不空并解析
func (c *conf) parseByAbsPath(confAbsPath string, obj interface{}) (err error) {
	if len(c.parsers) == 0 {
		return fmt.Errorf("no parser found")
	}
	return c.readConfDirect(confAbsPath, obj)
}

// 开始读取配置文件并解析
func (c *conf) readConfDirect(confPath string, obj interface{}) error {
	content, errIO := os.ReadFile(confPath)
	if errIO != nil {
		return errIO
	}
	// 读取文件扩展名，现在支持.toml .json
	fileExt := filepath.Ext(confPath)
	return c.ParseBytes(fileExt, content, obj)
}

// 配置里面如果设置了环境就返回设置好的，没有就返回default环境
func (c *conf) Env() env.AppEnv {
	if c.env == nil {
		return env.Default
	}
	return c.env
}

// 开始按照文件扩展名分配解析函数解析配置文件
func (c *conf) ParseBytes(fileExt string, content []byte, obj interface{}) error {
	parserFn, hasParser := c.parsers[fileExt]
	if fileExt == "" || !hasParser {
		return fmt.Errorf("%w, fileExt %q is not supported yet", fmt.Errorf("no parser found"), fileExt)
	}
	contentNew, errHelper := c.executeBeforeHelpers(content, c.helpers)
	if errHelper != nil {
		return fmt.Errorf("%w, content=\n%s", errHelper, string(contentNew))
	}
	if errParser := parserFn(contentNew, obj); errParser != nil {
		return fmt.Errorf("%w, content=\n%s", errParser, string(contentNew))
	}
	return nil
}

// executeBeforeHelpers 执行
func (c *conf) executeBeforeHelpers(input []byte, helpers []*beforeHelper) (output []byte, err error) {
	if len(helpers) == 0 {
		return input, nil
	}
	output = input
	for _, helper := range helpers {
		output, err = helper.Func(c, output)
		if err != nil {
			return nil, fmt.Errorf("beforeHelper=%q has error:%w", helper.Name, err)
		}
	}
	return output, err
}

// 检查该配置文件是否存在
func (c *conf) Exists(confName string) bool {
	info, err := os.Stat(c.confFileRealPath(confName))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// 注册解析能力
func (c *conf) RegisterParserFunc(fileExt string, fn ParserFunc) error {
	if _, has := c.parsers[fileExt]; has {
		return fmt.Errorf("parser=%q already exists", fileExt)
	}
	c.parsers[fileExt] = fn
	return nil
}

func (c *conf) RegisterBeforeFunc(name string, fn BeforeFunc) error {
	if name == "" {
		return fmt.Errorf("name is empty, not allow")
	}
	for _, h1 := range c.helpers {
		if name == h1.Name {
			return fmt.Errorf("beforeHelper=%q already exists", name)
		}
	}
	c.helpers = append(c.helpers, newBeforeHelper(name, fn))
	return nil
}

// 为了在编译期即确保实现了接口
var _ Conf = (*conf)(nil)
