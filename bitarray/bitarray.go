package bitarray

func IsHigh(bs []byte, post int) bool {
	index := post / 8
	flag := post % 8
	if flag > 0 {
		index += 1
	}
	if len(bs)-1 < index {
		return false
	}
	return bs[index]&(1<<flag) > 0
}

func SetHigh(v []byte, post int) []byte {
	var bs = make([]byte, len(v))
	copy(bs, v)
	index := post / 8
	flag := post % 8
	if flag > 0 {
		index += 1
	}
	for len(bs)-1 < index {
		bs = append(bs, byte(0))
	}

	bs[index] |= 1 << flag
	return bs
}

func SetLow(v []byte, post int) []byte {
	var bs = make([]byte, len(v))
	copy(bs, v)
	index := post / 8
	flag := post % 8
	if flag > 0 {
		index += 1
	}
	if len(bs)-1 < index {
		return bs
	}

	bs[index] &= ^(1 << flag)
	last := bs[len(bs)-1]
	for last == 0 && len(bs) > 1 {
		bs = bs[:len(bs)-1]
		last = bs[len(bs)-1]
	}
	return bs
}
