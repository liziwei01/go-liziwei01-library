package env

import (
	"errors"
	"net"
	"os"
	"strconv"
)

var (
	pid       int
	pidString string
	localIP   = "unknown"
)

// PID 得到 PID
func PID() int {
	return pid
}

// PIDString 得到PID 字符串形式
// 如打印日志的场景
func PIDString() string {
	return pidString
}

// LocalIP 本机IP，返回非127域的第一个ipv4 地址
// 极端特殊情况获取失败返回 机器名 或者 unknown
func LocalIP() string {
	return localIP
}
func init() {
	pid = os.Getpid()
	pidString = strconv.Itoa(pid)
	if val, err := localIPV4(); err == nil {
		localIP = val
	}
}
func localIPV4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return os.Hostname()
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("fail to get local ip")
}
