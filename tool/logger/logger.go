package logger

import (
	"context"
	"fmt"
	"letspay/common/constants"
	"log/slog"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func Info(ctx context.Context, msg string) {
	pc, file, line, _ := runtime.Caller(2)

	funcName := getFuncName(runtime.FuncForPC(pc).Name())
	fileName := filepath.Base(file)
	processId, _ := ctx.Value(constants.PROCESS_ID).(string)

	slog.Info(fmt.Sprintf("%s file=%s:%s func=%s %s=%s",
		msg,
		fileName,
		strconv.Itoa(line),
		funcName,
		constants.PROCESS_ID,
		processId,
	))
}

func Error(ctx context.Context, msg string) {

	pc, file, line, _ := runtime.Caller(2)

	funcName := getFuncName(runtime.FuncForPC(pc).Name())
	fileName := filepath.Base(file)
	processId, _ := ctx.Value(constants.PROCESS_ID).(string)

	slog.Error(fmt.Sprintf("%s file=%s:%s func=%s %s=%s",
		msg,
		fileName,
		strconv.Itoa(line),
		funcName,
		constants.PROCESS_ID,
		processId,
	))
}

func getFuncName(s string) string {
	index := strings.LastIndex(s, ".")
	return s[index+1:]
}
