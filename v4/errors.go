package gosql

import (
	"errors"
)

var (
	ErrShardPicker          = errors.New("pick shard error")
	ErrShardPickerOverLimit = errors.New("pick shard over limit")
)
