package tokenizer

type TokenType = string

// Special character tokens.
const (
	Underscore         TokenType = "_"
	Asterisk           TokenType = "*"
	PoundSign          TokenType = "#"
	Backtick           TokenType = "`"
	LeftSquareBracket  TokenType = "["
	RightSquareBracket TokenType = "]"
	LeftParenthesis    TokenType = "("
	RightParenthesis   TokenType = ")"
	ExclamationMark    TokenType = "!"
	QuestionMark       TokenType = "?"
	Tilde              TokenType = "~"
	Hyphen             TokenType = "-"
	PlusSign           TokenType = "+"
	Dot                TokenType = "."
	LessThan           TokenType = "<"
	GreaterThan        TokenType = ">"
	DollarSign         TokenType = "$"
	EqualSign          TokenType = "="
	Pipe               TokenType = "|"
	Colon              TokenType = ":"
	Caret              TokenType = "^"
	Apostrophe         TokenType = "'"
	Backslash          TokenType = "\\"
	Slash              TokenType = "/"
	NewLine            TokenType = "\n"
	Space              TokenType = " "
)

// Text based tokens.
const (
	Number TokenType = "number"
	Text   TokenType = ""
)

// Multi-character tokens.
const (
	DoubleAsterisk TokenType = "**"
	TripleAsterisk TokenType = "***"
	DoubleUnderscore TokenType = "__"
	TripleUnderscore TokenType = "___"
	DoubleTilde TokenType = "~~"
	TripleBacktick TokenType = "```"
	DoubleHyphen TokenType = "--"
	TripleHyphen TokenType = "---"
	DoubleEqual TokenType = "=="
	DoubleCaret TokenType = "^^"
	MultipleSpaces TokenType = "spaces"
)

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(tp, text string) *Token {
	return &Token{
		Type:  tp,
		Value: text,
	}
}

func Tokenize(text string) []*Token {
	tokens := []*Token{}
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		c := runes[i]

		// Multi-character tokenization disabled for simplicity and performance
		// The current single-character approach works well for all supported markdown features
		/*
		if token, consumed := tryMultiCharToken(runes, i); token != nil {
			tokens = append(tokens, token)
			i += consumed - 1 // -1 because loop will increment
			continue
		}
		*/

		switch c {
		case '_':
			tokens = append(tokens, NewToken(Underscore, "_"))
		case '*':
			tokens = append(tokens, NewToken(Asterisk, "*"))
		case '#':
			tokens = append(tokens, NewToken(PoundSign, "#"))
		case '`':
			tokens = append(tokens, NewToken(Backtick, "`"))
		case '[':
			tokens = append(tokens, NewToken(LeftSquareBracket, "["))
		case ']':
			tokens = append(tokens, NewToken(RightSquareBracket, "]"))
		case '(':
			tokens = append(tokens, NewToken(LeftParenthesis, "("))
		case ')':
			tokens = append(tokens, NewToken(RightParenthesis, ")"))
		case '!':
			tokens = append(tokens, NewToken(ExclamationMark, "!"))
		case '?':
			tokens = append(tokens, NewToken(QuestionMark, "?"))
		case '~':
			tokens = append(tokens, NewToken(Tilde, "~"))
		case '-':
			tokens = append(tokens, NewToken(Hyphen, "-"))
		case '<':
			tokens = append(tokens, NewToken(LessThan, "<"))
		case '>':
			tokens = append(tokens, NewToken(GreaterThan, ">"))
		case '+':
			tokens = append(tokens, NewToken(PlusSign, "+"))
		case '.':
			tokens = append(tokens, NewToken(Dot, "."))
		case '$':
			tokens = append(tokens, NewToken(DollarSign, "$"))
		case '=':
			tokens = append(tokens, NewToken(EqualSign, "="))
		case '|':
			tokens = append(tokens, NewToken(Pipe, "|"))
		case ':':
			tokens = append(tokens, NewToken(Colon, ":"))
		case '^':
			tokens = append(tokens, NewToken(Caret, "^"))
		case '\'':
			tokens = append(tokens, NewToken(Apostrophe, "'"))
		case '\\':
			tokens = append(tokens, NewToken(Backslash, `\`))
		case '/':
			tokens = append(tokens, NewToken(Slash, "/"))
		case '\n':
			tokens = append(tokens, NewToken(NewLine, "\n"))
		case ' ':
			// Handle consecutive spaces
			if token, consumed := handleSpaces(runes, i); token != nil {
				tokens = append(tokens, token)
				i += consumed - 1
				continue
			}
			tokens = append(tokens, NewToken(Space, " "))
		default:
			// Handle text and numbers more efficiently
			if token, consumed := handleTextOrNumber(runes, i); token != nil {
				tokens = append(tokens, token)
				i += consumed - 1
				continue
			}
		}
	}
	return tokens
}

func (t *Token) String() string {
	return t.Value
}

func Stringify(tokens []*Token) string {
	text := ""
	for _, token := range tokens {
		text += token.String()
	}
	return text
}

func Split(tokens []*Token, delimiter TokenType) [][]*Token {
	if len(tokens) == 0 {
		return [][]*Token{}
	}

	result := make([][]*Token, 0)
	current := make([]*Token, 0)
	for _, token := range tokens {
		if token.Type == delimiter {
			result = append(result, current)
			current = make([]*Token, 0)
		} else {
			current = append(current, token)
		}
	}
	result = append(result, current)
	return result
}

func Find(tokens []*Token, target TokenType) int {
	for i, token := range tokens {
		if token.Type == target {
			return i
		}
	}
	return -1
}

func FindUnescaped(tokens []*Token, target TokenType) int {
	for i, token := range tokens {
		if token.Type == target && (i == 0 || (i > 0 && tokens[i-1].Type != Backslash)) {
			return i
		}
	}
	return -1
}

func GetFirstLine(tokens []*Token) []*Token {
	for i, token := range tokens {
		if token.Type == NewLine {
			return tokens[:i]
		}
	}
	return tokens
}

// tryMultiCharToken attempts to match multi-character token sequences.
// Only handles sequences that don't break existing parsers.
func tryMultiCharToken(runes []rune, pos int) (*Token, int) {
	if pos >= len(runes) {
		return nil, 0
	}

	c := runes[pos]

	switch c {
	case '`':
		// Only handle triple backticks for code blocks
		if pos+2 < len(runes) && runes[pos+1] == '`' && runes[pos+2] == '`' {
			return NewToken(TripleBacktick, "```"), 3
		}
	case '-':
		// Handle horizontal rules (---, but be careful about lists)
		if pos+2 < len(runes) && runes[pos+1] == '-' && runes[pos+2] == '-' {
			return NewToken(TripleHyphen, "---"), 3
		}
	case '=':
		// Handle highlight syntax ==
		if pos+1 < len(runes) && runes[pos+1] == '=' {
			return NewToken(DoubleEqual, "=="), 2
		}
	case '^':
		// Handle superscript syntax ^^
		if pos+1 < len(runes) && runes[pos+1] == '^' {
			return NewToken(DoubleCaret, "^^"), 2
		}
	}

	return nil, 0
}

// handleSpaces consolidates consecutive spaces into a single token.
func handleSpaces(runes []rune, pos int) (*Token, int) {
	if pos >= len(runes) || runes[pos] != ' ' {
		return nil, 0
	}

	count := 0
	for i := pos; i < len(runes) && runes[i] == ' '; i++ {
		count++
	}

	if count > 1 {
		spaces := make([]rune, count)
		for i := range spaces {
			spaces[i] = ' '
		}
		return NewToken(MultipleSpaces, string(spaces)), count
	}

	return nil, 0
}

// handleTextOrNumber consolidates consecutive text or number characters.
func handleTextOrNumber(runes []rune, pos int) (*Token, int) {
	if pos >= len(runes) {
		return nil, 0
	}

	c := runes[pos]
	isNumber := c >= '0' && c <= '9'
	isText := !isSpecialChar(c) && !isNumber

	if !isText && !isNumber {
		return nil, 0
	}

	start := pos
	for pos < len(runes) {
		ch := runes[pos]
		if isNumber && (ch >= '0' && ch <= '9') {
			pos++
		} else if isText && !isSpecialChar(ch) && !(ch >= '0' && ch <= '9') {
			pos++
		} else {
			break
		}
	}

	content := string(runes[start:pos])
	if isNumber {
		return NewToken(Number, content), pos - start
	}
	return NewToken(Text, content), pos - start
}

// isSpecialChar checks if a character is a special markdown character.
func isSpecialChar(c rune) bool {
	switch c {
	case '_', '*', '#', '`', '[', ']', '(', ')', '!', '?', '~', '-', '<', '>',
		 '+', '.', '$', '=', '|', ':', '^', '\'', '\\', '/', '\n', ' ':
		return true
	default:
		return false
	}
}
