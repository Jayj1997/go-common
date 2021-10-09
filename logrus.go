/*
 * @Author       : jayj
 * @Date         : 2021-06-23 09:45:12
 * @Description  : 日志，包含日志分割、自动删除、hook等
 */
package common

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type MyHook struct {
}

type Data struct {
	MsgType  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

func InitLog() {

	path := "routinelog"

	writer, err := rotatelogs.New(
		// 文件路径 格式
		path+".%Y%m%d",
		// 为最新日志建立软连接
		rotatelogs.WithLinkName(path),
		// 最多保存的文件数量
		rotatelogs.WithRotationCount(5),
		// 每一天转为新文件
		rotatelogs.WithRotationTime(time.Duration(76400)*time.Second),
	)

	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	// 报警钩子
	// logrus.AddHook(&MyHook{})

	writers := []io.Writer{
		writer,
		os.Stdout,
	}

	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)

	if err == nil {
		logrus.SetOutput(fileAndStdoutWriter)
	} else {
		logrus.Info("failed to log to file.")
	}

	logrus.SetLevel(logrus.InfoLevel)
	// logrus.SetLevel(logrus.WarnLevel)
	// logrus.SetLevel(logrus.ErrorLevel)
	// logrus.SetLevel(logrus.InfoLevel)
}
