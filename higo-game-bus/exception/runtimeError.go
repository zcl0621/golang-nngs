package exception

import (
	"github.com/gin-gonic/gin"
	"higo-game-bus/logger"
	"higo-game-bus/responses"
	"net/http"
)

func StandardRuntimeWarnError() *StandardRuntimeWarn {
	return &StandardRuntimeWarn{BaseCodeError: BaseCodeError{ErrorCode: error_code_unkonw_error}}
}

func StandardRuntimeBadError() *StandardRuntimeError {
	return &StandardRuntimeError{StandardRuntimeWarn: StandardRuntimeWarn{BaseCodeError: BaseCodeError{ErrorCode: error_code_unkonw_error}}}
}

type BaseCodeError struct {
	ErrorCode     int
	OutPutMessage string
}

func (runtime *BaseCodeError) Error() string {
	return runtime.OutPutMessage
}

func (runtime *BaseCodeError) SetErrorCode(code int) *BaseCodeError {
	runtime.ErrorCode = code
	return runtime
}

func (runtime *BaseCodeError) SetOutPutMessage(msg string) *BaseCodeError {
	runtime.OutPutMessage = msg
	return runtime
}

func (runtime *BaseCodeError) printError() {}

func (runtime *BaseCodeError) buildErrContext(context *gin.Context) {
	context.Writer.Header().Set("Error-Code", "1")
	context.JSON(http.StatusBadRequest, &responses.ErrorResponse{
		Err: runtime.OutPutMessage,
	})
}

func (runtime *BaseCodeError) buildStandardContext(context *gin.Context) {
	context.Writer.Header().Set("Error-Code", "1")
	context.JSON(http.StatusOK, &responses.StandardResponse{Code: runtime.ErrorCode, Msg: runtime.OutPutMessage})
}

type StandardRuntimeWarn struct {
	BaseCodeError
	Parameter    []interface{}
	FunctionName string
}

func (runtime *StandardRuntimeWarn) SetErrorCode(code int) *StandardRuntimeWarn {
	runtime.ErrorCode = code
	return runtime
}

func (runtime *StandardRuntimeWarn) SetOutPutMessage(msg string) *StandardRuntimeWarn {
	runtime.OutPutMessage = msg
	return runtime
}

func (runtime *StandardRuntimeWarn) SetFunctionName(functionName string) *StandardRuntimeWarn {
	runtime.FunctionName = functionName
	return runtime
}

func (runtime *StandardRuntimeWarn) SetParameter(parameter ...interface{}) *StandardRuntimeWarn {
	runtime.Parameter = parameter
	return runtime
}

type RuntimeWarn struct {
	StandardRuntimeWarn
}

func (runtime *RuntimeWarn) SetErrorCode(code int) *RuntimeWarn {
	runtime.ErrorCode = code
	return runtime
}

func (runtime *RuntimeWarn) SetOutPutMessage(msg string) *RuntimeWarn {
	runtime.OutPutMessage = msg
	return runtime
}

func (runtime *RuntimeWarn) SetFunctionName(functionName string) *RuntimeWarn {
	runtime.FunctionName = functionName
	return runtime
}

func (runtime *RuntimeWarn) SetParameter(parameter ...interface{}) *RuntimeWarn {
	runtime.Parameter = parameter
	return runtime
}

func (runtime *RuntimeWarn) buildErrContext(context *gin.Context) {
	context.Writer.Header().Set("Error-Code", "1")
	context.JSON(http.StatusBadRequest, &responses.ErrorResponse{Err: runtime.OutPutMessage})
}

func (runtime *RuntimeWarn) buildStandardContext(context *gin.Context) {
	context.Writer.Header().Set("Error-Code", "1")
	context.JSON(http.StatusOK, &responses.StandardResponse{Code: runtime.ErrorCode, Msg: runtime.OutPutMessage})
}

type StandardRuntimeError struct {
	StandardRuntimeWarn
	OriginalError error
}

func (runtime *StandardRuntimeError) SetErrorCode(code int) *StandardRuntimeError {
	runtime.ErrorCode = code
	return runtime
}

func (runtime *StandardRuntimeError) SetOutPutMessage(msg string) *StandardRuntimeError {
	runtime.OutPutMessage = msg
	return runtime
}

func (runtime *StandardRuntimeError) SetFunctionName(functionName string) *StandardRuntimeError {
	runtime.FunctionName = functionName
	return runtime
}

func (runtime *StandardRuntimeError) SetParameter(parameter ...interface{}) *StandardRuntimeError {
	runtime.Parameter = parameter
	return runtime
}

func (runtime *StandardRuntimeError) SetOriginalError(originalError error) *StandardRuntimeError {
	runtime.OriginalError = originalError
	return runtime
}

func (runtime *StandardRuntimeError) printError() {
	logger.Logger(runtime.FunctionName, logger.ERROR, runtime, runtime.Parameter)
}

type RuntimeError struct {
	StandardRuntimeError
}

func (runtime *RuntimeError) SetErrorCode(code int) *RuntimeError {
	runtime.ErrorCode = code
	return runtime
}

func (runtime *RuntimeError) SetOutPutMessage(msg string) *RuntimeError {
	runtime.OutPutMessage = msg
	return runtime
}

func (runtime *RuntimeError) SetFunctionName(functionName string) *RuntimeError {
	runtime.FunctionName = functionName
	return runtime
}

func (runtime *RuntimeError) SetParameter(parameter ...interface{}) *RuntimeError {
	runtime.Parameter = parameter
	return runtime
}

func (runtime *RuntimeError) SetOriginalError(originalError error) *RuntimeError {
	runtime.OriginalError = originalError
	return runtime
}

func (runtime *RuntimeError) printError() {
	logger.Logger(runtime.FunctionName, logger.ERROR, runtime, runtime.Parameter)
}

func (runtime *RuntimeError) buildErrContext(context *gin.Context) {
	context.Writer.Header().Set("Error-Code", "1")
	context.JSON(http.StatusBadRequest, &responses.ErrorResponse{Err: runtime.OutPutMessage})
}

func (runtime *RuntimeError) buildStandardContext(context *gin.Context) {
	context.Writer.Header().Set("Error-Code", "1")
	context.JSON(http.StatusOK, &responses.StandardResponse{Code: runtime.ErrorCode, Msg: runtime.OutPutMessage})
}
