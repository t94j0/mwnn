package client

func red(text string) string {
	return "\x1b[0;31m" + text + "\x1b[0;0m"
}

func blue(text string) string {
	return "\x1b[0;34m" + text + "\x1b[0;0m"
}

func green(text string) string {
	return "\x1b[0;32m" + text + "\x1b[0;0m"
}
func purple(text string) string {
	return "\x1b[0;35m" + text + "\x1b[0;0m"
}
