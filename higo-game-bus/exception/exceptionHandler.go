package exception

import (
	"github.com/gin-gonic/gin"
	"higo-game-bus/logger"
	"net/http"
)

func ExceptionStandardHandler(context *gin.Context) {
	defer func() {
		var err = recover()
		if err != nil {
			processStandardError(err, context)
			context.Abort()
		}
	}()
	context.Next()
}

func processStandardError(exception interface{}, context *gin.Context) {
	switch t := exception.(type) {
	case ErrorContextProcess:
		t.printError()
		t.buildStandardContext(context)
	default:
		logger.Logger("errorLogPrint", logger.ERROR, exception.(error), nil)
		context.Writer.Header().Set("Error-Code", "1")
		context.JSON(http.StatusBadRequest, "操作失败")
	}
}

func ExceptionErrHandler(context *gin.Context) {
	defer func() {
		var err = recover()
		if err != nil {
			processErrError(err, context)
			context.Abort()
		}
	}()
	context.Next()
}

func processErrError(exception interface{}, context *gin.Context) {
	switch t := exception.(type) {
	case ErrorContextProcess:
		t.printError()
		t.buildErrContext(context)
	default:
		logger.Logger("errorLogPrint", logger.ERROR, exception.(error), nil)
		context.Writer.Header().Set("Error-Code", "1")
		context.JSON(http.StatusBadRequest, "操作失败")
	}
}

type ErrorContextProcess interface {
	printError()
	buildStandardContext(context *gin.Context)
	buildErrContext(context *gin.Context)
}
