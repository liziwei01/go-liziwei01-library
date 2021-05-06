package conf
import (
	"bytes"
	"os"
	"regexp"
)

// BeforeFunc 辅助回调方法，在执行ParserFunc前，会先对配置的内容进行解析处理
type BeforeFunc func(conf Conf, content []byte) ([]byte, error)

// beforeHelper 辅助功能
// 在正式解析配置前执行
type beforeHelper struct {
	Name string
	Func BeforeFunc
}

// 构建结构体
func newBeforeHelper(name string, fn BeforeFunc) *beforeHelper {
	return &beforeHelper{
		Name: name,
		Func: fn,
	}
}

// defaultHelpers 默认的helper方法：获取环境变量
var defaultHelpers = []*beforeHelper{
	newBeforeHelper("env", helperOsEnvVars),
}

// 模板变量格式：{env.变量名} 或者 {env.变量名|默认值}
var osEnvVarReg = regexp.MustCompile(`\{env\.([A-Za-z0-9_]+)(\|[^}]+)?\}`)

// helperOsEnvVars 将配置文件中的 {env.xxx} 的内容，从环境变量中读取并替换
func helperOsEnvVars(conf Conf, content []byte) ([]byte, error) {
	contentNew := osEnvVarReg.ReplaceAllFunc(content, func(subStr []byte) []byte {
		// 将 {env.xxx} 中的 xxx 部分取出
		// 或者 将 {env.yyy|val} 中的 yyy|val 部分取出
		keyWithDefaultVal := subStr[len("{env.") : len(subStr)-1] // eg: xxx 或者 yyy|val
		idx := bytes.Index(keyWithDefaultVal, []byte("|"))
		if idx > 0 {
			// {env.变量名|默认值} 有默认值的格式
			key := string(keyWithDefaultVal[:idx])  // eg: yyy
			defaultVal := keyWithDefaultVal[idx+1:] // eg: val
			envVal := os.Getenv(key)
			if envVal == "" {
				return defaultVal
			}
			return []byte(envVal)
		}
		// {env.变量名} 无默认值的部分
		return []byte(os.Getenv(string(keyWithDefaultVal)))
	})
	return contentNew, nil
}
