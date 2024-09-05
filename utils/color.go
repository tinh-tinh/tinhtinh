package utils

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

type Color struct {
	Code string
	Val  string
}

func Log(colors ...Color) {
	for _, v := range colors {
		print(v.Code + v.Val + ResetCode)
	}
}

func Red(str string) Color {
	return Color{Code: RedCode, Val: str}
}

func Green(str string) Color {
	return Color{Code: GreenCode, Val: str}
}

func Yellow(str string) Color {
	return Color{Code: YellowCode, Val: str}
}

func Blue(str string) Color {
	return Color{Code: BlueCode, Val: str}
}

func Magenta(str string) Color {
	return Color{Code: MagentaCode, Val: str}
}

func Cyan(str string) Color {
	return Color{Code: CyanCode, Val: str}
}

func Gray(str string) Color {
	return Color{Code: GrayCode, Val: str}
}

func White(str string) Color {
	return Color{Code: WhiteCode, Val: str}
}
