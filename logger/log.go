package logger

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

func Config(dir, name string, maxSizeFileMB, maxBackups, maxAgeDays int, compress bool) {
	logFile := path.Join(dir, name+".log")
	// 设置 lumberjack 为日志输出目标
	logHandleHook := &lumberjack.Logger{
		Filename:   logFile,       // 日志文件路径
		MaxSize:    maxSizeFileMB, // 每个日志文件最大10MB
		MaxBackups: maxBackups,    // 最多保留5个旧日志文件
		MaxAge:     maxAgeDays,    // 最多保留30天
		Compress:   compress,      // 是否压缩旧日志
	}

	// 创建 logs 目录（如果不存在）
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		logrus.Errorf("failed to create logs directory: %v", err)
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, logHandleHook)
	logrus.SetOutput(multiWriter)

}
