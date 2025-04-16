package crawler

type ICrawler interface {
	SetHeader(header map[string]string)
	Visit(url string) ([]byte, error)
}
