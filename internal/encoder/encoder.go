package encoder

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type CodeEncoder interface {
	Encode(num int64) string
}

type Base62Encoder struct{}

func NewBase62Encoder() *Base62Encoder {
	return &Base62Encoder{}
}

func (e *Base62Encoder) Encode(num int64) string {
	if num == 0 {
		return string(alphabet[0])
	}
	result := make([]byte, 0)

	for num > 0 {
		remainder := num % 62
		result = append(result, alphabet[remainder])
		num = num / 62
	}

	reverse(result)

	return string(result)
}

func reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
