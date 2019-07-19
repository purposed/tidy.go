package fsclean

import "unicode"

const (
	// MaxTokenLength represents the maximum number of runes in a single token.
	MaxTokenLength = 64
)

type tokenizeMode int

const (
	modeIdle tokenizeMode = iota
	modeString
	modeOperator
	modeParen
)

// IsSeparator returns whether the rune is a separator.
func IsSeparator(ch rune) bool {
	return ch == ' '
}

// IsOperator returns whether a char is an operator member.
func IsOperator(ch rune) bool {
	operators := []rune{'>', '<', '=', '!', '#', '?', '^', '$'}

	for _, o := range operators {
		if o == ch {
			return true
		}
	}
	return false
}

// IsAlphaNum returns whether a char is alphanumeric.
func IsAlphaNum(ch rune) bool {
	return unicode.IsDigit(ch) || unicode.IsLetter(ch) || ch == '.'
}

func tokenize(str string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		runes := []rune(str)

		currentMode := modeIdle

		var currentToken []rune

		runeCount := len(runes)

		for i := 0; i < runeCount; i++ {
			ch := runes[i]
			if IsSeparator(ch) {
				currentMode = modeIdle
				continue
			}

			if currentMode == modeString {
				if !IsAlphaNum(ch) {
					currentMode = modeIdle
				} else {
					currentToken = append(currentToken, ch)
				}
			}

			if currentMode == modeOperator {
				if !IsOperator(ch) {
					currentMode = modeIdle
				} else {
					currentToken = append(currentToken, ch)
				}
			}

			if currentMode == modeParen {
				currentMode = modeIdle
			}

			if currentMode == modeIdle {
				if len(currentToken) > 0 {
					out <- string(currentToken)
					currentToken = []rune{}
				}

				currentToken = append(currentToken, ch)

				if IsAlphaNum(ch) {
					currentMode = modeString
				} else if IsOperator(ch) {
					currentMode = modeOperator
				}
			}
		}
		if len(currentToken) > 0 {
			out <- string(currentToken)
		}

		out <- ""
	}()

	return out
}

// Tokenize splits the text into tokens & filters out invalid tokens.
func Tokenize(str string) <-chan string {
	return tokenize(str)
}
