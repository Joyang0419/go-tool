package consts

func JSONHeaders() map[string]string {
	return map[string]string{
		"Content-Type": ContentTypeJSON,
		"Accept":       ContentTypeJSON,
	}
}

func FormHeaders() map[string]string {
	return map[string]string{
		"Content-Type": ContentTypeForm,
	}
}
