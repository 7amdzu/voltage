package mask

// MaskBytes takes a byte slice, replaces every PAN match in-place
// (zero string conversions), and returns the masked bytes.
func MaskBytes(b []byte) []byte {
	return panRegex.ReplaceAllFunc(b, func(match []byte) []byte {
		// match is the raw PAN bytes, e.g. "4111 1111 1111 1111"
		n := len(match)
		if n < 4 {
			// too short to mask; return unchanged
			return match
		}
		// allocate a new slice: prefix + last 4 bytes of the PAN
		// "**** **** **** " is 15 bytes long
		out := make([]byte, 15+4)
		// copy the mask prefix
		copy(out, []byte("**** **** **** "))
		// copy the last 4 digits from the match
		copy(out[15:], match[n-4:])
		return out
	})
}
