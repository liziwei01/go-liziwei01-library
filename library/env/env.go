package env

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	// DefaultAppName 默认的app名称
	DefaultAppName = "liziwei01APP"
	// DefaultRunMode 测试默认运行等级
	DefaultRunMode = RunModeRelease
)

// 可以依据不同的运行等级来开启不同的调试功能、接口
const (
	// RunModeDebug 调试
	RunModeDebug = "debug"
	// RunModeTest 测试
	RunModeTest = "test"
	// RunModeRelease 线上发布
	RunModeRelease = "release"
)

// Option 具体的环境信息
//
// 所有的选项都是可选的
type Option struct {
	// AppName 应用名称
	AppName string
	// RunMode 运行模式
	RunMode string
	// RootDir 应用根目录地址
	// 若为空，将通过自动推断的方式获取
	RootDir string
	// DataDir 应用数据根目录地址
	// 默认为 RootDir+"/data/"
	DataDir string
	// LogDir 应用日志根目录地址
	// 默认为 RootDir+"/log/"
	LogDir string
	// ConfDir 应用配置文件根目录地址
	// 默认为RootDir+"/conf/"
	ConfDir string
}

// String 序列化，方便查看
// 目前输出的是一个json
func (opt Option) String() string {
	format := `{"AppName":%q,"RootDir":%q,"DataDir":%q,"LogDir":%q,"ConfDir":%q,"RunMode":%q}`
	return fmt.Sprintf(format, opt.AppName, opt.RootDir, opt.DataDir, opt.LogDir, opt.ConfDir, opt.RunMode)
}

// Merge 合并
// 传入的Option不为空则合并，否则使用老的值
func (opt Option) Merge(newOpt Option) Option {
	return Option{
		AppName: SecondStrFirst(opt.AppName, newOpt.AppName),
		RunMode: SecondStrFirst(opt.RunMode, newOpt.RunMode),
		RootDir: SecondStrFirst(opt.RootDir, newOpt.RootDir),
		DataDir: SecondStrFirst(opt.DataDir, newOpt.DataDir),
		LogDir:  SecondStrFirst(opt.LogDir, newOpt.LogDir),
		ConfDir: SecondStrFirst(opt.ConfDir, newOpt.ConfDir),
	}
}

func SecondStrFirst(v1 string, v2 string) string {
	if v2 != "" {
		return v2
	}
	return v1
}

// RootDirEnv 应用根目录环境信息接口
type RootDirEnv interface {
	RootDir() string
}

// ConfDirEnv 配置环境信息接口
type ConfDirEnv interface {
	ConfDir() string
}

// DataDirEnv 数据目录环境信息接口
type DataDirEnv interface {
	DataDir() string
}

// LogDirEnv 日志目录环境信息接口
type LogDirEnv interface {
	LogDir() string
}

// AppNameEnv 应用名称接口
type AppNameEnv interface {
	AppName() string
}

// RunModeEnv 运行模式/等级接口
type RunModeEnv interface {
	RunMode() string
}

// AppEnv 应用环境信息完整的接口定义
type AppEnv interface {
	// 应用名称
	AppNameEnv
	// 应用根目录
	RootDirEnv
	// 应用配置文件根目录
	ConfDirEnv
	// 应用数据文件根目录
	DataDirEnv
	// 应用日志文件更目录
	LogDirEnv
	// 应用运行情况
	RunModeEnv
	// 获取当前环境的选项详情
	Options() Option
	// 复制一个新的env对象，并将传入的Option merge进去
	CloneWithOption(opt Option) AppEnv
}

// New 创建新的应用环境
func New(opt Option) AppEnv {
	env := &appEnv{}
	if opt.AppName != "" {
		env.setAppName(opt.AppName)
	}
	if opt.RunMode != "" {
		env.setRunMode(opt.RunMode)
	}
	if opt.RootDir != "" {
		env.setRootDir(opt.RootDir)
	}
	if opt.ConfDir != "" {
		env.setConfDir(opt.ConfDir)
	}
	if opt.DataDir != "" {
		env.setDataDir(opt.DataDir)
	}
	if opt.LogDir != "" {
		env.setLogDir(opt.LogDir)
	}
	return env
}

// appEnv 实现了以上接口
type appEnv struct {
	rootDir string
	dataDir string
	confDir string
	logDir  string
	appName string
	runMode string
}

// 所有环境变量设定时都走日志输出
func setValue(addr *string, value string, fieldName string) {
	*addr = value
	_ = log.Output(2, fmt.Sprintf("[env] set %q=%q\n", fieldName, value))
}

// 获取AppName
func (a *appEnv) AppName() string {
	if a.appName != "" {
		return a.appName
	}
	return DefaultAppName
}

// 设定AppName
func (a *appEnv) setAppName(name string) {
	setValue(&a.appName, name, "AppName")
}

// 获取RunMode
func (a *appEnv) RunMode() string {
	if a.runMode != "" {
		return a.runMode
	}
	return DefaultRunMode
}

// 设定RunMode
func (a *appEnv) setRunMode(mod string) {
	setValue(&a.runMode, mod, "RunMode")
}

// 获取RootDir
func (a *appEnv) RootDir() string {
	if a.rootDir != "" {
		return a.rootDir
	}
	return AutoDetectAppRootDir()
}

// 设定RootDir
func (a *appEnv) setRootDir(dir string) {
	setValue(&a.rootDir, dir, "RootDir")
}

// 获取DataDir
func (a *appEnv) DataDir() string {
	return a.chooseDir(a.dataDir, "data")
}

// 设定DataDir
func (a *appEnv) setDataDir(dir string) {
	setValue(&a.dataDir, dir, "DataDir")
}

// 获取LogDir
func (a *appEnv) LogDir() string {
	return a.chooseDir(a.logDir, "log")
}

// 设定LogDir
func (a *appEnv) setLogDir(dir string) {
	setValue(&a.logDir, dir, "LogDir")
}

// 获取ConfDir
func (a *appEnv) ConfDir() string {
	return a.chooseDir(a.confDir, "conf")
}

// 设定ConfDir
func (a *appEnv) setConfDir(dir string) {
	setValue(&a.confDir, dir, "ConfDir")
}

// 获取Dir的时候不直接返回，走个兜底，如果为空就拼接RootDir和文件夹名字
func (a *appEnv) chooseDir(dir string, subDirName string) string {
	if dir != "" {
		return dir
	}
	return filepath.Join(a.RootDir(), subDirName)
}

// 以Option的形式输出现存的环境配置信息
func (a *appEnv) Options() Option {
	return Option{
		AppName: a.AppName(),
		RunMode: a.RunMode(),
		RootDir: a.RootDir(),
		DataDir: a.DataDir(),
		LogDir:  a.LogDir(),
		ConfDir: a.ConfDir(),
	}
}

// 建立新环境配置的时候可以用老环境做base
func (a *appEnv) CloneWithOption(opt Option) AppEnv {
	opts := a.Options().Merge(opt)
	return New(opts)
}

// 为了在编译期即确保appEnv实现了AppEnv接口
var _ AppEnv = (*appEnv)(nil)

// AutoDetectAppRootDir 自动获取应用根目录
// 定义为变量，这样若默认实现不满足，可进行替换
var AutoDetectAppRootDir = autoDetect

// 通过找go.mod文件或含有app.toml的conf文件夹的位置来确定APP现在工作目录
func autoDetect() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	names := []string{
		"go.mod",
		filepath.Join("conf", "app.toml"),
	}
	dir, err := findDirMatch(wd, names)
	if err == nil {
		return dir
	}
	return wd
}

// findDirMatch 在指定目录下，向其父目录查找对应的文件是否存在
// 若存在，则返回匹配到的路径
func findDirMatch(baseDir string, fileNames []string) (dir string, err error) {
	currentDir := baseDir
	// 最多寻找20层父目录
	for i := 0; i < 20; i++ {
		for _, fileName := range fileNames {
			depsPath := filepath.Join(currentDir, fileName)
			// 找到即返回
			if _, err := os.Stat(depsPath); !os.IsNotExist(err) {
				return currentDir, nil
			}
		}
		// 向父目录走
		currentDir = filepath.Dir(currentDir)
		if currentDir == "." {
			break
		}
	}
	return "", fmt.Errorf("cannot found")
}
