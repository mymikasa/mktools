package mlog

import (
	"testing"
)

func TestInfo(t *testing.T) {

	// SetLogLevel(logrus.DebugLevel.String())
	o := WithOutPutPath("aa.log")
	MustSetUp(o)
	// var log defaultConfig
	mLog.Info("the producer group has been created, specify another one", map[string]interface{}{
		"test": "rsss",
	})
	Info("the producer group has been created, specify another one", map[string]interface{}{
		"test": "rsss",
	})
}
