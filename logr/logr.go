package logr

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

const (
	// ColourReset code
	ColourReset = "\x1B[0m"
)

type (
	// Message used to send log message to logger goroutine
	Message struct {
		T      Type
		Time   time.Time
		Args   []interface{}
		chDone chan bool
	}
)

var (
	// Level sets the default log filter level.
	Level = Critical
	// NewMessageFunc factory function for log message
	NewMessageFunc = func() interface{} {
		return &Message{}
	}

	writers  = make(map[io.Writer]Type, 0)
	pool     = &sync.Pool{}
	messages = make(chan *Message, 100)
)

// Colour setting for colourful log output
var (
	Colour  = true
	ColourP = "\x1B[38;5;124m"
	ColourE = "\x1B[38;5;124m"
	ColourW = "\x1B[38;5;208m"
	ColourI = "\x1B[38;5;33m"
	ColourD = "\x1B[38;5;153m"
	ColourS = "\x1B[38;5;34m"
)

func init() {
	pool.New = NewMessageFunc
	go golog()
}

// Wait for log messages to be processed
func Wait() {
	time.Sleep(time.Nanosecond * 1000)
}

// Output logs matching the given type filter to the given writers.
func Output(t Type, to ...io.Writer) {
	for _, w := range to {
		writers[w] = t
	}
}

// Panic logs inputs as panics and panics
func Panic(v ...interface{}) {
	code := log(P, true, v...)
	panic(code)
}

// Panicf logs a formatted message as a panic and panics
func Panicf(msg string, v ...interface{}) {
	code := logf(P, true, msg, v...)
	panic(code)
}

// Error logs inputs as errors
func Error(v ...interface{}) string {
	return log(E, false, v...)
}

// Errorf logs a formatted message as an error
func Errorf(msg string, v ...interface{}) string {
	return logf(E, false, msg, v...)
}

// Warn logs inputs as warnings
func Warn(v ...interface{}) string {
	return log(W, false, v...)
}

// Warnf logs a formatted message as a warning
func Warnf(msg string, v ...interface{}) string {
	return logf(W, false, msg, v...)
}

// Info logs inputs as info messages
func Info(v ...interface{}) string {
	return log(I, false, v...)
}

// Infof logs a formatted message as an info message
func Infof(msg string, v ...interface{}) string {
	return logf(I, false, msg, v...)
}

// Debug logs inputs as debug messages
func Debug(v ...interface{}) string {
	return log(D, false, v...)
}

// Debugf logs a formatted message as a debug message
func Debugf(msg string, v ...interface{}) string {
	return logf(D, false, msg, v...)
}

// Success logs inputs as success messages
func Success(v ...interface{}) string {
	return log(S, false, v...)
}

// Successf logs a formatted message as a success message
func Successf(msg string, v ...interface{}) string {
	return logf(S, false, msg, v...)
}

// log inputs to given type
func log(t Type, wait bool, v ...interface{}) string {
	done := make(chan bool)
	m := pool.Get().(*Message)
	m.T = t
	m.Time = time.Now()
	m.Args = v
	m.chDone = done
	c := m.Code()
	messages <- m

	if wait {
		<-done
	}

	return c
}

// Concurrently work through the logs buffered channel
func golog() {
	for {
		m := <-messages
		t := m.T
		s := m.String()
		ch := m.chDone
		m.Reset()
		pool.Put(m)

		for w, l := range writers {
			if t&l != t {
				continue
			}
			w.Write([]byte(s))
		}

		close(ch)
	}
}

// format a msg and log as given type
func logf(t Type, wait bool, msg string, v ...interface{}) string {
	return log(t, wait, fmt.Sprintf(msg, v...))
}

// Code returns a message code for later tracking
func (m *Message) Code() string {
	return strconv.FormatInt(m.Time.UnixNano(), 36)
}

// String implements fmt.Stringer
func (m *Message) String() string {
	msg := fmt.Sprintf("%-25s | %s | %s | %s\n", m.Time.Format("Jan 02 2006 15:04:05.9999"), m.Code(), m.T, fmt.Sprint(m.Args...))
	if !Colour {
		return msg
	}

	switch m.T {
	case P:
		msg = ColourP + msg + ColourReset
	case E:
		msg = ColourE + msg + ColourReset
	case W:
		msg = ColourW + msg + ColourReset
	case I:
		msg = ColourI + msg + ColourReset
	case D:
		msg = ColourD + msg + ColourReset
	case S:
		msg = ColourS + msg + ColourReset
	}
	return msg
}

// Reset the message object for later reuse
func (m *Message) Reset() {
	m.T = 0
	m.Time = time.Time{}
	m.Args = make([]interface{}, 0)
}
