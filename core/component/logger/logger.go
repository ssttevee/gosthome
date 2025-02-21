package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type Level slog.Level

const (
	LevelNone        = Level(slog.LevelError + 4)
	LevelError       = Level(slog.LevelError)
	LevelWarn        = Level(slog.LevelWarn)
	LevelInfo        = Level(slog.LevelInfo)
	LevelConfig      = Level(slog.LevelInfo - 2)
	LevelDebug       = Level(slog.LevelDebug)
	LevelVerbose     = Level(slog.LevelDebug - 4)
	LevelVeryVerbose = Level(slog.LevelDebug - 8)
)

const (
	_levelNone        string = "none"
	_levelError       string = "error"
	_levelWarn        string = "warn"
	_levelInfo        string = "info"
	_levelConfig      string = "config"
	_levelDebug       string = "debug"
	_levelVerbose     string = "verbose"
	_levelVeryVerbose string = "very_verbose"
)

const (
	_levelIntFirst       = _levelIntNone
	_levelIntNone        = 0
	_levelIntError       = 1
	_levelIntWarn        = 2
	_levelIntInfo        = 3
	_levelIntConfig      = 4
	_levelIntDebug       = 5
	_levelIntVerbose     = 6
	_levelIntVeryVerbose = 7
	_levelIntLast        = _levelIntVeryVerbose
)

var _levelFromInt = [_levelIntLast + 1]Level{
	_levelIntNone:        LevelNone,
	_levelIntError:       LevelError,
	_levelIntWarn:        LevelWarn,
	_levelIntInfo:        LevelInfo,
	_levelIntConfig:      LevelConfig,
	_levelIntDebug:       LevelDebug,
	_levelIntVerbose:     LevelVerbose,
	_levelIntVeryVerbose: LevelVeryVerbose,
}

var _levelToInt = map[Level]int{
	LevelNone:        _levelIntNone,
	LevelError:       _levelIntError,
	LevelWarn:        _levelIntWarn,
	LevelInfo:        _levelIntInfo,
	LevelConfig:      _levelIntConfig,
	LevelDebug:       _levelIntDebug,
	LevelVerbose:     _levelIntVerbose,
	LevelVeryVerbose: _levelIntVeryVerbose,
}

var _levelMap = map[Level]string{
	LevelNone:        _levelNone,
	LevelError:       _levelError,
	LevelWarn:        _levelWarn,
	LevelInfo:        _levelInfo,
	LevelConfig:      _levelConfig,
	LevelDebug:       _levelDebug,
	LevelVerbose:     _levelVerbose,
	LevelVeryVerbose: _levelVeryVerbose,
}

var _levelStringMap = map[string]Level{
	_levelNone:        LevelNone,
	_levelError:       LevelError,
	_levelWarn:        LevelWarn,
	_levelInfo:        LevelInfo,
	_levelConfig:      LevelConfig,
	_levelDebug:       LevelDebug,
	_levelVerbose:     LevelVerbose,
	_levelVeryVerbose: LevelVeryVerbose,
}

var ErrInvalidLevel = errors.New("not a valid Level")

// ParseLevel attempts to convert a string to a Level.
func ParseLevel(name string) (Level, error) {
	if x, ok := _levelStringMap[name]; ok {
		return x, nil
	}
	return Level(0), fmt.Errorf("%s is %w", name, ErrInvalidLevel)
}

func ParseLevelFromInt[T ~int | ~int32 | ~int64](l T) Level {
	return _levelFromInt[min(max(int(l), _levelIntFirst), _levelIntLast)]
}

func (x Level) String() string {
	if str, ok := _levelMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Level(%d)", x)
}
func (x Level) IsValid() bool {
	_, ok := _levelMap[x]
	return ok
}
func (x Level) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}
func (x *Level) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseLevel(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

type Logger struct {
	slog.Logger
}

// Config logs at [LogLevelConfig].
func (l *Logger) Config(msg string, args ...any) {
	l.Log(context.Background(), slog.Level(LevelConfig), msg, args...)
}

// ConfigContext logs at [LogLevelConfig] with the given context.
func (l *Logger) ConfigContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.Level(LevelConfig), msg, args...)
}

// Verbose logs at [LogLevelVerbose].
func (l *Logger) Verbose(msg string, args ...any) {
	l.Log(context.Background(), slog.Level(LevelVerbose), msg, args...)
}

// VerboseContext logs at [LogLevelVerbose] with the given context.
func (l *Logger) VerboseContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.Level(LevelConfig), msg, args...)
}

// VeryVerbose logs at [LogLevelVeryVerbose].
func (l *Logger) VeryVerbose(msg string, args ...any) {
	l.Log(context.Background(), slog.Level(LevelVeryVerbose), msg, args...)
}

// VeryVerboseContext logs at [LogLevelVeryVerbose] with the given context.
func (l *Logger) VeryVerboseContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.Level(LevelVeryVerbose), msg, args...)
}
