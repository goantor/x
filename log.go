package x

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Roboter interface {
	Send(message string) error
}

type ILogger interface {
	IContextData
	TakeContextData() IContextData
	GiveContextData(data IContextData)
	WithRobot() ILogger
	GiveRobot(robot Roboter) ILogger
	Info(message string, data H)
	Trace(message string, data H)
	Debug(message string, data H)
	Warn(message string, data H)
	Fatal(message string, data H)
	Panic(message string, data H)
	Error(message string, err error, data H)
}

func NewLoggerWithData(log *logrus.Logger, data IContextData) ILogger {
	return &logger{
		log:          log,
		IContextData: data,
	}
}

func NewLogger(log *logrus.Logger) ILogger {
	return &logger{
		log:          log,
		IContextData: NewContextData(),
	}
}

type logger struct {
	IContextData
	log      *logrus.Logger
	robot    Roboter
	masker   *MaskReflect
	useRobot bool
}

func (l *logger) takeMasker() *MaskReflect {
	if l.masker == nil {
		l.masker = &MaskReflect{}
	}

	return l.masker
}

func (l *logger) GiveContextData(data IContextData) {
	l.IContextData = data
}

func (l *logger) TakeContextData() IContextData {
	return l.IContextData
}

func (l *logger) GiveRobot(robot Roboter) ILogger {
	l.robot = robot
	return l
}

func (l *logger) WithRobot() ILogger {
	if l.robot != nil {
		l.useRobot = true
	}

	return l
}

func (l *logger) resetRobot() {
	if l.useRobot {
		l.useRobot = false
	}
}

func (l *logger) makeFields(data H) (fields logrus.Fields) {
	fields = make(logrus.Fields)
	//fields["data"] = l.takeMasker().MakeMask(data)
	if data == nil {
		data = H{}
	}

	fields["data"] = data
	fields["context"] = l.IContextData
	return
}

func (l *logger) doLog(level logrus.Level, message string, data H) {
	go l.log.WithFields(l.makeFields(data)).Log(level, message)
}

func (l *logger) Info(message string, data H) {
	l.doLog(logrus.InfoLevel, message, data)
}

func (l *logger) Trace(message string, data H) {
	l.doLog(logrus.TraceLevel, message, data)
}

func (l *logger) Debug(message string, data H) {
	l.doLog(logrus.DebugLevel, message, data)
}

func (l *logger) Warn(message string, data H) {
	l.doLog(logrus.WarnLevel, message, data)
}

func (l *logger) Fatal(message string, data H) {
	l.doLog(logrus.FatalLevel, message, data)
}

func (l *logger) Panic(message string, data H) {
	l.doLog(logrus.PanicLevel, message, data)
}

func (l *logger) Error(message string, err error, data H) {
	fields := l.makeFields(data)
	fields["error"] = err
	go l.log.WithFields(fields).Error(message)
}

type ILoggerOption interface {
	TakeStdout() bool
	TakeHooks() []logrus.Hook
	TakeLevel() string
	TakeFormatter() logrus.Formatter
	TakeReportCaller() bool
	ModifyName(name string)
	Clone() ILoggerOption
}

func NewLogBuilder(opt ILoggerOption) *LogBuilder {
	return &LogBuilder{opt: opt}
}

type LogBuilder struct {
	opt ILoggerOption
}

func (b LogBuilder) makeHooks() (hooks logrus.LevelHooks) {
	hooks = make(logrus.LevelHooks)
	for _, hook := range b.opt.TakeHooks() {
		hooks.Add(hook)
	}

	return
}

func (b LogBuilder) makeLevel() (level logrus.Level) {
	level, _ = logrus.ParseLevel(b.opt.TakeLevel())
	return
}

func (b LogBuilder) makeStdout() io.Writer {
	if b.opt.TakeStdout() {
		return os.Stderr
	}

	file, _ := os.OpenFile(os.DevNull, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	return file
}

func (b LogBuilder) Make() (entity *logrus.Logger) {
	entity = &logrus.Logger{
		Out:          b.makeStdout(),
		Formatter:    b.opt.TakeFormatter(),
		Hooks:        b.makeHooks(),
		Level:        b.makeLevel(),
		ExitFunc:     os.Exit,
		ReportCaller: b.opt.TakeReportCaller(),
	}

	return
}
