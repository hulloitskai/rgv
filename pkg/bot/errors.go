package bot

// ErrorCode identifies a type of error.
type ErrorCode uint8

const (
	// InvalidConfig is an error that describes some invalid Bot configuration.
	InvalidConfig ErrorCode = iota
)

// Error is the common structure for errors created by package bot.
type Error struct {
	Code ErrorCode
	msg  string
}

func (e *Error) Error() string {
	return e.msg
}
