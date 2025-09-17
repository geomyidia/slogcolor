package slogcolor

type SourceFileMode int

const (
	// Nop does nothing.
	Nop SourceFileMode = iota

	// ShortFile produces only the filename (for example main.go:69).
	ShortFile

	// MediumFile produces the relative file path from project root (for example cmd/server/main.go:69).
	MediumFile

	// LongFile produces the full file path (for example /home/user/go/src/myapp/main.go:69).
	LongFile
)
