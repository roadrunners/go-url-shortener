package key

var (
	keyChar   = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	decodeMap = makeDecodeMap()
)

func makeDecodeMap() map[byte]int64 {
	m := make(map[byte]int64)
	for i, b := range keyChar {
		m[b] = int64(i)
	}
	return m
}

func GenKey(n int64) string {
	if n == 0 {
		return string(keyChar[0])
	}
	l := int64(len(keyChar))
	s := make([]byte, 20)
	i := int64(len(s))
	for n > 0 && i >= 0 {
		i--
		j := n % l
		n = (n - j) / l
		s[i] = keyChar[j]
	}
	return string(s[i:])
}

func GenId(key string) int64 {
	l := int64(len(keyChar))
	n := int64(0)
	for _, b := range key {
		n *= l
		n += decodeMap[byte(b)]
	}
	return n
}
