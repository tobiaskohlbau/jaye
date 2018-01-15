package multimedia

import (
	"context"
	"io"
)

// Converter converts and merges multimedia types
type Converter interface {
	Convert(ctx context.Context, src io.Reader, dst io.Writer) error
	Merge(ctx context.Context, video, audio io.Reader, dst io.Writer) error
}
