package log

import (
	l "log"
)

func Example() {
	name := "Leaf"

	Debug("My name is %v", name)
	Release("My name is %v", name)
	Error("My name is %v", name)
	// log.Fatal("My name is %v", name)

	logger, err := New("release", "", l.LstdFlags)
	if err != nil {
		return
	}
	defer logger.Close()

	logger.Debug("will not print")
	logger.Release("My name is %v", name)

	Export(logger)

	Debug("will not print")
	Release("My name is %v", name)
}
