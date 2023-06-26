package command

import "bytes"

func NextLine(s []byte) []byte {
	i := bytes.IndexByte(s, '\n')
	if i == -1 {
		return nil
	}
	return s[i+1:]
}
func ReadLine(s []byte) []byte {
	i := bytes.IndexByte(s, '\n')
	if i == -1 {
		return nil
	}
	return s[:i]
}
func NextField(s []byte) ([]byte, []byte) {
	// Skip whitespace.
	for i, b := range s {
		if b != ' ' && b != 10 && b != 0 && b != 9 {
			s = s[i:]
			break
		}
		if i == len(s)-1 {
			s = s[i+1:]
			break
		}
	}
	// Up until the next whitespace field.
	for i, b := range s {
		if b == ' ' || b == 10 || b == 0 || b == 9 {
			return s[:i], s[i:]
		}
	}
	return s, nil
}

func ReadField(s []byte, idx int) (val []byte, remain []byte) {
	if idx == 0 {
		return []byte{}, s
	}
	remain = s
	for i := 1; i <= idx; i++ {
		if len(remain) == 0 {
			return
		}
		val, remain = NextField(remain)
	}
	return
}
