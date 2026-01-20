package easylog

import (
	"fmt"
	"io"
	"log"
	"time"
)

var SRVLOG_BEGIN = "BEGIN"
var SRVLOG_END = "END"

// application log: to log the request and response infomation of a request
// service log: to log the running status of on step in one operate
//	output: the output media, etc. log file, stdout, stdin ...
//	reqtime: the time the request recieved. in us
//	reqid: request id. every user request is asigned a unique request,
//		and passed through all over the service call between service modules.
//	module: the module processes the request
//	operate: the operate that processing the request
//	params: the parameters of this request
//	rsptime: the time the request begin to response
//	rescode: the result code of this request
//	elapse: the elapse of processing this request,
//		include data transfer time of the response. in us
//	comment: comment of this log
func AppLog(output io.Writer, reqtime, reqid, source, module, operate, params, rsptime string, rspcode, elapse int, comment string) {
	fmt.Fprintf(output, "%s %s %s %s %s %s %s %d %d %s\n", reqtime, reqid, source, module, operate, params, rsptime, rspcode, elapse, comment)
}

// service log: to log the running status of on step in one operate
//	output: the output media, etc. log file, stdout, stdin ...
//	reqid: request id. every user request is asigned a unique request,
//		and passed through all over the service call between service modules.
//	module: the module processes the request
//	operate: the operate of current step that processing the request
//	level: the log level. include:
//		1. INFO: used in normal situation, to log the status
//			of processing
//		2. WARN: used in warning situation, means error occurs
//			but the process can continue
//		3. ERROR: used in error situation, means error occurs
//			and the bisness logic interrupted,
//			but the system process continue running
//		4. FATAL: used in error situation, means error occurs
//			and the system process cannot continue,
//			and the program is going to exit
//	rescode: the result code of this processing step
//	elapse: the elapse of current processing step. in us
//	comment: comment of this log
func SrvLog(output io.Writer, reqid, module, operate, level string, rescode, elapse int, comment string) {
	now := time.Now()
	fmt.Fprintf(output, "%s %s %s %s %s %d %d %s\n", now.Format("2006-01-02 15:04:05.999999"),
		reqid, module, operate, level, rescode, elapse, comment)
}

// wrappers of service log
// Info service log
func SrvInfo(output io.Writer, reqid, module, operate string, rescode, elapse int, comment string) {
	SrvLog(output, reqid, module, operate, "INFO", rescode, elapse, comment)
}

// Warning service log
func SrvWarn(output io.Writer, reqid, module, operate string, rescode, elapse int, comment string) {
	SrvLog(output, reqid, module, operate, "WARN", rescode, elapse, comment)
}

// Error service log
func SrvError(output io.Writer, reqid, module, operate string, rescode, elapse int, comment string) {
	SrvLog(output, reqid, module, operate, "ERROR", rescode, elapse, comment)
}

// Fatal service log
func SrvFatal(output io.Writer, reqid, module, operate string, rescode, elapse int, comment string) {
	SrvLog(output, reqid, module, operate, "FATAL", rescode, elapse, comment)
}

// wrapper struct of service log, in order to easy the use of service log.
//	wraps the reqid, module, and output.
// only used in on goroutine and passed by argument throughout the function
//	call.
type SrvLogger struct {
	logger *log.Logger
	reqid  string
	module string
}

// creator
func NewSrvLogger(output io.Writer, reqid, module string) *SrvLogger {
	return &SrvLogger{log.New(output, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile), reqid, module}
}

// NewSrvLoggerFromLogger creator
func NewSrvLoggerFromLogger(logger *log.Logger, reqid, module string) *SrvLogger {
	return &SrvLogger{logger, reqid, module}
}

func (this *SrvLogger) output(calldepth int, operate, level string, rescode, elapse int, comment string) {
	this.logger.Output(calldepth+1, fmt.Sprintf("%s %s %s %s %d %d %s",
		this.reqid, this.module, operate, level, rescode, elapse, comment))
}

// general log
func (this *SrvLogger) Log(operate, level string, rescode, elapse int, comment string) {
	this.output(2, operate, level, rescode, elapse, comment)
}

// Info
func (this *SrvLogger) Info(operate string, rescode, elapse int, comment string) {
	this.output(2, operate, "INFO", rescode, elapse, comment)
}

// warning
func (this *SrvLogger) Warn(operate string, rescode, elapse int, comment string) {
	this.output(2, operate, "WARN", rescode, elapse, comment)
}

// error
func (this *SrvLogger) Error(operate string, rescode, elapse int, comment string) {
	this.output(2, operate, "ERROR", rescode, elapse, comment)
}

// fatal
func (this *SrvLogger) Fatal(operate string, rescode, elapse int, comment string) {
	this.output(2, operate, "FATAL", rescode, elapse, comment)
}
