package mask

import (
	"io"
)

const tailSize = 15

// StreamMasker wraps an io.Reader, keeps a small rune-buffer
// to catch PANs split across Read() calls, and masks them.
type StreamMasker struct {
	src io.Reader // the underlying reader (e.g. HTTP body)
	buf []rune    // trailing runes from previous Read
}

func NewStreamMasker(r io.Reader) *StreamMasker {
	return &StreamMasker{src: r}
}

// Read reads up to len(p) bytes, masks any full PANs,
// and preserves the last tailSize runes in buf.
func (m *StreamMasker) Read(p []byte) (int, error) {
	// Read raw bytes into tmp
	tmp := make([]byte, len(p))
	n, err := m.src.Read(tmp)
	if n == 0 {
		return 0, err
	}

	// Convert raw bytes to runes and prepend leftover buf
	runes := append(m.buf, []rune(string(tmp[:n]))...)

	// Split into “safe” (everything except last tailSize runes)
	// and new buf (last tailSize runes)
	var safe []rune
	if len(runes) > tailSize {
		safe, m.buf = runes[:len(runes)-tailSize], runes[len(runes)-tailSize:]
	} else {
		// not enough data yet—keep all in buf, nothing to process
		m.buf = runes
		return 0, err
	}

	// Mask PANs in the safe part
	masked := panRegex.ReplaceAllStringFunc(string(safe), maskPAN)

	// Copy masked bytes into p and return
	written := copy(p, []byte(masked))
	return written, err
}
