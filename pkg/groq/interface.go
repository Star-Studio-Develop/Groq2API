package groq

type Headers map[string]string

func (h Headers) Set(key string, value string) {
	h[key] = value
}

func NewHeader() Headers {
	return make(Headers)
}
