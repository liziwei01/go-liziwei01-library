package conf

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
)

// ParserFunc 针对特定文件后缀的配置解析方法
// 当前已经内置了 .toml  和 .json的解析方法
type ParserFunc func(bf []byte, obj interface{}) error

const (
	// 已支持的文件后缀
	// FileTOML toml
	FileTOML = ".toml"
	// FileJSON  json
	FileJSON = ".json"
)

// stripComment 去除单行的'#'注释
// 只支持单行，不支持行尾
func stripComment(input []byte) (out []byte) {
	var buf bytes.Buffer
	lines := bytes.Split(input, []byte("\n"))
	for _, line := range lines {
		lineN := bytes.TrimSpace(line)
		if !bytes.HasPrefix(lineN, []byte("#")) {
			buf.Write(line)
		}
		buf.WriteString("\n")
	}
	return bytes.TrimSpace(buf.Bytes())
}

// DefaultParserFuncs 所有默认的ParserFunc
var DefaultParserFuncs = map[string]ParserFunc{
	FileJSON: JSONParserFunc,
	FileTOML: TOMLParserFunc,
}

// 若内容以 # 开头，则该为注释
func jsonParserFunc(txt []byte, obj interface{}) error {
	bf := stripComment(txt)
	dec := json.NewDecoder(bytes.NewReader(bf))
	dec.UseNumber()
	return dec.Decode(obj)
}

// JSONParserFunc .json配置文件格式解析函数
var JSONParserFunc ParserFunc = jsonParserFunc

// TOMLParserFunc .toml配置文件格式解析函数
var TOMLParserFunc ParserFunc = toml.Unmarshal
