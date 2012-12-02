package semver

type itemType int

const (
	itemInvalid itemType = iota // found an invalid character
	itemEOF
	itemDot        // version part separator
	itemDash       // prerelease separator
	itemPlus       // build separator
	itemNumber     // [0-9]
	itemIdentifier // [0-9A-Za-z-]
)

const eof = 255

type item struct {
	itemType
	value string
}

func (i item) String() string {
	switch i.itemType {
	case itemInvalid:
		return "invalid character"
	case itemEOF:
		return "EOF"
	case itemDot:
		return "dot"
	case itemDash:
		return "dash"
	case itemPlus:
		return "plus"
	}
	return i.value
}

type lexer struct {
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	items chan item // channel of scanned items.
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() (next byte) {
	if l.pos >= len(l.input) {
		return eof
	}
	next = l.input[l.pos]
	l.pos += 1
	return
}

func (l *lexer) backup() {
	l.pos -= 1
}

func (l *lexer) peek() byte {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// Lexer state code

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for state := lexFirstPart; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// Part with only numbers and dots
func lexFirstPart(l *lexer) stateFn {
	c := l.next()
	switch {
	case c == eof:
		l.emit(itemEOF)
		return nil
	case c == '.':
		l.emit(itemDot)
	case c == '-':
		l.emit(itemDash)
		return lexSecondPart
	case c == '+':
		l.emit(itemPlus)
		return lexSecondPart
	case '0' <= c && c <= '9':
		return lexNumber
	default:
		l.emit(itemInvalid)
		return nil
	}
	return lexFirstPart
}

// Part after first dash, or plus
func lexSecondPart(l *lexer) stateFn {
	c := l.next()
	switch {
	case c == eof:
		l.emit(itemEOF)
		return nil
	case c == '.':
		l.emit(itemDot)
	case c == '+':
		l.emit(itemPlus)
	case '0' <= c && c <= '9',
		'A' <= c && c <= 'Z',
		'a' <= c && c <= 'z',
		c == '-':
		return lexIdentifier
	default:
		l.emit(itemInvalid)
		return nil
	}
	return lexSecondPart
}

func lexNumber(l *lexer) stateFn {
	c := l.next()
	switch {
	case '0' <= c && c <= '9':
		break // stay in this lexing mode
	case c == '.', c == '-', c == '+':
		l.backup()
		l.emit(itemNumber)
		return lexFirstPart
	case c == eof:
		l.emit(itemNumber)
		l.emit(itemEOF)
		return nil
	case 'A' <= c && c <= 'Z',
		'a' <= c && c <= 'z':
		return lexIdentifier // it's not a number!
	default:
		l.emit(itemInvalid)
		return nil
	}
	return lexNumber
}

func lexIdentifier(l *lexer) stateFn {
	c := l.next()
	switch {
	case '0' <= c && c <= '9',
		'A' <= c && c <= 'Z',
		'a' <= c && c <= 'z',
		c == '-':
		break // stay in this lexing mode
	case c == '.', c == '+':
		l.backup()
		l.emit(itemIdentifier)
		return lexSecondPart
	case c == eof:
		l.emit(itemIdentifier)
		l.emit(itemEOF)
		return nil
	default:
		l.emit(itemInvalid)
		return nil
	}
	return lexIdentifier
}
