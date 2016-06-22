package ansi

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Up moves the cursor up n lines.
func Up(n int) string {
	if n <= 0 {
		panic("invalid argument")
	}
	return escape(strconv.Itoa(n) + "A")
}

// Down moves the cursor down n lines.
func Down(n int) string {
	if n <= 0 {
		panic("invalid argument")
	}
	return escape(strconv.Itoa(n) + "B")
}

// Right moves the cursor right n columns.
func Right(n int) string {
	if n <= 0 {
		panic("invalid argument")
	}
	return escape(strconv.Itoa(n) + "C")
}

// Left moves the cursor left n columns.
func Left(n int) string {
	if n <= 0 {
		panic("invalid argument")
	}
	return escape(strconv.Itoa(n) + "D")
}

// Cursor moves the cursor to line y, column x. This indexes from the top left
// corner of the screen.
func Cursor(x, y int) string {
	if x <= 0 || y <= 0 {
		panic("invalid argument")
	}
	return escape(strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H")
}

// Clear clears the screen and move the cursor to the top left corner.
func Clear() string {
	return escape("H\x1B[2J")
}

// Erase erases everything to the left of the cursor on the line the cursor is
// on.
func Erase() string {
	return escape("0K")
}

// Reset resets all colours and text styles to the default.
func Reset() string {
	return escape("0m")
}

// Bold formats text bold.
func Bold(state bool) string {
	if state {
		return escape("1m")
	}
	return escape("22m")
}

// Underline formats text underlined.
func Underline(state bool) string {
	if state {
		return escape("4m")
	}
	return escape("24m")
}

// Blink formats text blinking.
func Blink(state bool) string {
	if state {
		return escape("5m")
	}
	return escape("25m")
}

// Reverse swaps foreground and background colour.
func Reverse(state bool) string {
	if state {
		return escape("7m")
	}
	return escape("27m")
}

// Black colours text black.
func Black() string {
	return escape("30m")
}

// Red colours text red.
func Red() string {
	return escape("31m")
}

// Green colours text green.
func Green() string {
	return escape("32m")
}

// Yellow colours text yellow.
func Yellow() string {
	return escape("33m")
}

// Blue colours text blue.
func Blue() string {
	return escape("34m")
}

// Magenta colours text magenta.
func Magenta() string {
	return escape("35m")
}

// Cyan colours text cyan.
func Cyan() string {
	return escape("36m")
}

// Grey colours text grey.
func Grey() string {
	return escape("90m")
}

// White colours text white.
func White() string {
	return escape("97m")
}

// BrightRed colours text bright red.
func BrightRed() string {
	return escape("91m")
}

// BrightGreen colours text bright green.
func BrightGreen() string {
	return escape("92m")
}

// BrightYellow colours text bright yellow.
func BrightYellow() string {
	return escape("93m")
}

// BrightBlue colours text bright blue.
func BrightBlue() string {
	return escape("94m")
}

// BrightMagenta colours text bright magenta.
func BrightMagenta() string {
	return escape("95m")
}

// BrightCyan colours text bright cyan.
func BrightCyan() string {
	return escape("96m")
}

// BrightGrey colours text bright grey.
func BrightGrey() string {
	return escape("37m")
}

// BlackBg colours text in black background.
func BlackBg() string {
	return escape("40m")
}

// RedBg colours text in red background.
func RedBg() string {
	return escape("41m")
}

// GreenBg colours text in green background.
func GreenBg() string {
	return escape("42m")
}

// YellowBg colours text in yellow background.
func YellowBg() string {
	return escape("43m")
}

// BlueBg colours text in blue background.
func BlueBg() string {
	return escape("44m")
}

// MagentaBg colours text in magenta background.
func MagentaBg() string {
	return escape("45m")
}

// CyanBg colours text in cyan background.
func CyanBg() string {
	return escape("46m")
}

// GreyBg colours text in grey background.
func GreyBg() string {
	return escape("100m")
}

// WhiteBg colours text in white background.
func WhiteBg() string {
	return escape("107m")
}

// BrightRedBg colours text in bright red background.
func BrightRedBg() string {
	return escape("101m")
}

// BrightGreenBg colours text in bright green background.
func BrightGreenBg() string {
	return escape("102m")
}

// BrightYellowBg colours text in bright yellow background.
func BrightYellowBg() string {
	return escape("103m")
}

// BrightBlueBg colours text in bright blue background.
func BrightBlueBg() string {
	return escape("104m")
}

// BrightMagentaBg colours text in bright magenta background.
func BrightMagentaBg() string {
	return escape("105m")
}

// BrightCyanBg colours text in bright cyan background.
func BrightCyanBg() string {
	return escape("106m")
}

// BrightGreyBg colours text in bright grey background.
func BrightGreyBg() string {
	return escape("47m")
}

// RemoveEscapeSequences removes ANSI escape sequences.
func RemoveEscapeSequences() {
	// TODO(uwe): Implement
}

func escape(s string) string {
	return "\x1B[" + s
}

// Colors defines the standard HTML colours for ANSI.
var Colors = []string{
	"#000000", "#Dd0000", "#00CF12", "#C2CB00", "#3100CA", "#E100C6", "#00CBCB", "#C7C7C7",
	"#686868", "#FF5959", "#00FF6B", "#FAFF5C", "#775AFF", "#FF47FE", "#0FFFFF", "#FFFFFF",
}

// ToHTML converts ANSI escape sequences into HTML span tags.
func ToHTML(text []byte) []byte {
	re := regexp.MustCompile("\u001B\\[([0-9A-Za-z;]+)m([^\u001B]+)")
	matches := re.FindAllSubmatch(text, -1)
	if matches == nil {
		return text
	}

	var buf bytes.Buffer

	for _, match := range matches {
		bg, fg := -1, -1
		var bold, underline, negative bool

		codes := bytes.Split(match[1], []byte(";"))
		for _, c := range codes {
			code, _ := strconv.Atoi(string(c))
			if code == 0 {
				bg, fg = -1, -1
				bold, underline, negative = false, false, false
			} else if code == 1 {
				bold = true
			} else if code == 4 {
				underline = true
			} else if code == 7 {
				negative = true
			} else if code == 21 {
				bold = false
			} else if code == 24 {
				underline = false
			} else if code == 27 {
				negative = false
			} else if code >= 30 && code <= 37 {
				fg = code - 30
			} else if code == 39 {
				fg = -1
			} else if code >= 40 && code <= 47 {
				bg = code - 40
			} else if code == 49 {
				bg = -1
			} else if code >= 90 && code <= 97 {
				fg = code - 90 + 8
			} else if code >= 100 && code <= 107 {
				bg = code - 100 + 8
			}
		}

		style := ""
		if negative {
			fg = bg
			bg = fg
		}
		if bold {
			fg = fg | 8
			style += "font-weight: bold;"
		}
		if underline {
			style += "text-decoration:underline"
		}
		if fg >= 0 {
			style += fmt.Sprintf("color: %s;", Colors[fg])
		}
		if bg >= 0 {
			style += fmt.Sprintf("background-color: %s;", Colors[bg])
		}

		html := string(match[2])
		html = strings.Replace(html, "&", "&amp;", -1)
		html = strings.Replace(html, "<", "&lt;", -1)
		html = strings.Replace(html, ">", "&gt;", -1)

		if style == "" {
			buf.WriteString(html)
		} else {
			buf.WriteString(fmt.Sprintf(`<span style="%s">%s</span>`, style, html))
		}
	}

	return buf.Bytes()
}
