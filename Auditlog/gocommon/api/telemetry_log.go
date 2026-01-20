package api

// Package telemetry telemetry sdk
// @File logger.go
// @Description  telemetry sdk for log and trace
import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/resource"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/encoder"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/log"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/open_standard"
	truntime "github.com/kweaver-ai/TelemetrySDK-Go/span/v2/runtime"
)

const (
	sameLogTemp    = "samelog start %s end %s total %d"
	attrPosition   = "Position"
	attrExtMessage = "ExtMessage"
	fieldCap       = 2
)

var (
	lvlMap = map[int]field.StringField{
		log.DebugLevel: log.DebugLevelString,
		log.InfoLevel:  log.InfoLevelString,
		log.WarnLevel:  log.WarnLevelString,
		log.ErrorLevel: log.ErrorLevelString,
	}
)

// TLogger Telemetry日志记录器
type TLogger struct {
	ctx        context.Context
	LogLevel   int
	calldepth  int
	lastLog    string
	lastLevel  int
	times      int64
	start      time.Time
	mu         sync.Mutex
	Runtime    *truntime.Runtime
	tagMarkers []TagMarker
}

// TagMarker 标签标记器
type TagMarker func(level field.Field, typ string, message field.Field) []string

// position 日志位置
type position struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

type messageTags struct {
	Tags []string `json:"tags"`
}

// AddTagMarker 添加标签标记器
func (t *TLogger) AddTagMarker(marker TagMarker) {
	t.tagMarkers = append(t.tagMarkers, marker)
}

func (t *TLogger) isSame(lvl int, msg *string) (isSame bool, lastLevel int, samelog string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.lastLog == *msg {
		t.times++
		isSame = true
		return
	}
	now := time.Now()
	if t.times > 1 {
		samelog = fmt.Sprintf(sameLogTemp, t.start.Format(time.RFC3339), now.Format(time.RFC3339), t.times)
		lastLevel = t.lastLevel
	}
	t.lastLog = *msg
	t.lastLevel = lvl
	t.start = time.Now()
	t.times = 1
	return
}

// LogOptions 日志选项
type LogOptions struct {
	output      io.Writer
	level       int
	serviceInfo *LogOptionServiceInfo
	callDepth   int
}

// LogOption 日志选项
type LogOption interface {
	Apply(*LogOptions)
}

// LogOptionLevel 日志等级
type LogOptionLevel int

// Apply 选项配置应用
func (op LogOptionLevel) Apply(ops *LogOptions) {
	ops.level = int(op)
}

// LogOptionCallDepth 日志堆栈显示初始等级
type LogOptionCallDepth int

// Apply 选项配置应用
func (op LogOptionCallDepth) Apply(ops *LogOptions) {
	ops.callDepth = int(op)
}

// LogOptionServiceInfo 服务信息
type LogOptionServiceInfo struct {
	Name     string
	Version  string
	Instance string
}

// Apply 选项配置应用
func (op *LogOptionServiceInfo) Apply(ops *LogOptions) {
	ops.serviceInfo = op
}

// DefaultLogOptions 默认log选项
func DefaultLogOptions() *LogOptions {
	return &LogOptions{
		output:    os.Stdout,
		level:     log.InfoLevel,
		callDepth: 3,
	}
}

// NewTelemetryLogger 新建日志
func NewTelemetryLogger(opts ...interface{}) Logger {
	ops := DefaultLogOptions()
	for _, op := range opts {
		switch v := op.(type) {
		case io.Writer:
			ops.output = v // 兼容原函数
		case int:
			ops.level = v // 兼容原函数
		case LogOption:
			v.Apply(ops)
		}
	}
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(ops.output),
	}
	if ops.serviceInfo != nil {
		resource.SetServiceName(ops.serviceInfo.Name)
		resource.SetServiceVersion(ops.serviceInfo.Version)
		resource.SetServiceInstance(ops.serviceInfo.Instance)
	}
	writer.Resource = resource.LogResource()
	run := truntime.NewRuntime(writer, field.NewSpanFromPool)
	// start runtime
	go run.Run()
	return &TLogger{
		LogLevel:   ops.level,
		calldepth:  ops.callDepth,
		Runtime:    run,
		tagMarkers: []TagMarker{},
	}
}

func newRecord(typ string, message field.Field) field.Field {
	record := field.MallocStructField(fieldCap)
	record.Set(typ, message)
	record.Set("Type", field.StringField(typ))
	return record
}

func (t *TLogger) writeLogField(typ string, message, level field.Field, options ...field.LogOptionFunc) {
	span := t.Runtime.Children(context.Background())
	defer span.Signal()
	span.SetLogLevel(level)
	record := newRecord(typ, message)
	span.SetRecord(record)
	if t.ctx != nil {
		span.SetOption(field.WithContext(t.ctx))
	}
	span.SetOption(t.getOptionTag(level, typ, message))
	span.SetOption(options...)
}

func (t *TLogger) writeLog(message string, level field.Field, options ...field.LogOptionFunc) {
	span := t.Runtime.Children(context.Background())
	defer span.Signal()
	span.SetLogLevel(level)
	typ := "Message"
	msg := field.StringField(message)
	record := field.MallocStructField(1)
	record.Set(typ, msg)
	span.SetRecord(record)
	if t.ctx != nil {
		span.SetOption(field.WithContext(t.ctx))
	}
	span.SetOption(t.getOptionTag(level, typ, msg))
	span.SetOption(options...)
}

func (t *TLogger) log(lvl int, v ...interface{}) {
	if t.LogLevel > lvl {
		return
	}
	msg := fmt.Sprint(v...)
	same, lastLvl, sameLog := t.isSame(lvl, &msg)
	if same {
		return
	}
	if lastLvl > 0 {
		t.writeLog(sameLog, lvlMap[lastLvl])
	}
	t.writeLog(msg, lvlMap[lvl], t.getOptionPosition(t.calldepth))
}

func isField(v []interface{}) (f field.Field, ok bool) {
	if len(v) > 0 {
		f, ok = v[0].(field.Field)
	}
	return
}

func (t *TLogger) logf(lvl int, format string, v ...interface{}) {
	if t.LogLevel > lvl {
		return
	}
	if f, ok := isField(v); ok {
		t.writeLogField(format, f, lvlMap[lvl])
		return
	}

	msg := fmt.Sprintf(format, v...)
	same, lastLvl, sameLog := t.isSame(lvl, &msg)
	if same {
		return
	}
	if lastLvl > 0 {
		t.writeLog(sameLog, lvlMap[lastLvl])
	}
	t.writeLog(msg, lvlMap[lvl], t.getOptionPosition(t.calldepth))
}

func (t *TLogger) Debugln(v ...interface{}) {
	t.log(log.DebugLevel, v...)
}
func (t *TLogger) Infoln(v ...interface{}) {
	t.log(log.InfoLevel, v...)
}
func (t *TLogger) Warnln(v ...interface{}) {
	t.log(log.WarnLevel, v...)
}
func (t *TLogger) Errorln(v ...interface{}) {
	t.log(log.ErrorLevel, v...)
}

func (t *TLogger) Debugf(format string, v ...interface{}) {
	t.logf(log.DebugLevel, format, v...)
}
func (t *TLogger) Infof(format string, v ...interface{}) {
	t.logf(log.InfoLevel, format, v...)
}
func (t *TLogger) Warnf(format string, v ...interface{}) {
	t.logf(log.WarnLevel, format, v...)
}
func (t *TLogger) Errorf(format string, v ...interface{}) {
	t.logf(log.ErrorLevel, format, v...)
}

// Output 输出Infor日志
func (t *TLogger) Output(calldepth int, msg string) (err error) {
	if t.LogLevel > log.InfoLevel {
		return
	}
	same, lastLvl, sameLog := t.isSame(log.InfoLevel, &msg)
	if same {
		return
	}
	if lastLvl > 0 {
		t.writeLog(sameLog, lvlMap[lastLvl])
	}
	t.writeLog(msg, lvlMap[log.InfoLevel], t.getOptionPosition(calldepth))
	return
}

// WithContext 获取上下文
func (t *TLogger) WithContext(ctx context.Context) *TLogger {
	t.ctx = ctx
	return t
}

func (t *TLogger) getOptionTag(level field.Field, typ string, message field.Field) field.LogOptionFunc {
	var tags []string
	for _, tagMarker := range t.tagMarkers {
		tags = append(tags, tagMarker(level, typ, message)...)
	}

	attr := field.NewAttribute(attrExtMessage, field.MallocJsonField(messageTags{
		Tags: tags,
	}))

	return field.WithAttribute(attr)
}

func (t *TLogger) getOptionPosition(calldepth int) field.LogOptionFunc {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "-"
		line = 0
	}
	attr := field.NewAttribute(attrPosition, field.MallocJsonField(position{
		File: file,
		Line: line,
	}))
	return field.WithAttribute(attr)
}
