package color

const (
	RedCode     = "\x1b[31m"
	GreenCode   = "\x1b[32m"
	YellowCode  = "\x1b[33m"
	BlueCode    = "\x1b[34m"
	MagentaCode = "\x1b[35m"
	CyanCode    = "\x1b[36m"
	ResetCode   = "\x1b[0m"
	GrayCode    = "\033[37m"
	WhiteCode   = "\033[97m"
)

func Print(str string) string {
	return (str + ResetCode)
}

func Red(str string) string {
	return Print(RedCode + str)
}

func Green(str string) string {
	return Print(GreenCode + str)
}

func Yellow(str string) string {
	return Print(YellowCode + str)
}

func Blue(str string) string {
	return Print(BlueCode + str)
}

func Magenta(str string) string {
	return Print(MagentaCode + str)
}

func Cyan(str string) string {
	return Print(CyanCode + str)
}

func Gray(str string) string {
	return Print(GrayCode + str)
}

func White(str string) string {
	return Print(WhiteCode + str)
}
