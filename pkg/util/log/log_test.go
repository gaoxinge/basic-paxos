package log

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestLog(t *testing.T) {
	Output(os.Stderr)
	With()
	Level(zerolog.DebugLevel)
	Sample(zerolog.Often)
	Hook(zerolog.NewLevelHook())
	Err(nil)
	Trace()
	Debug()
	Info()
	Warn()
	Error()
	Fatal()
	Panic()
	WithLevel(zerolog.DebugLevel)
	Log()
	Print()
	Printf("printf")
	Ctx(context.TODO())
}
