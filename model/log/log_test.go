package log_test

import (
	"testing"

	"github.com/liziwei01/go-liziwei01-library/model/log"
	lib "github.com/baidu/go-lib/log"

)

func TestLog(t *testing.T) {
	log.Init("test")
	lib.Logger.Info("shit")
}
