package ctl

func StrPtrToStr(str *string) string {
	if str == nil {
		return "null"
	}

	return *str
}
