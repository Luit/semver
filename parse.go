package semver

import (
	"errors"
	"strconv"
)

var (
	expectedNumber      = errors.New("semver.Parse: unexpected character, expected number")
	expectedDashPlusEOF = errors.New("semver.Parse: unexpected character, expected dash, plus or end of string")
	expectedDot         = errors.New("semver.Parse: unexpected character, expected dot")
	expectedDotEOF      = errors.New("semver.Parse: unexpected character, expected dot or end of string")
	expectedDotPlusEOF  = errors.New("semver.Parse: unexpected character, expected dot, plus or end of string")
	expectedIdentifier  = errors.New("semver.Parse: unexpected character, expected identifier")
	unexpectedEOF       = errors.New("semver.Parse: unexpected end of string")
	invalidCharacter    = errors.New("semver.Parse: invalid character in string")
)

func parse(version string) (v Version, err error) {
	l := lex(version)

	if v.Major, err = parseNumber(l); err != nil {
		return
	}
	if err = parseDot(l); err != nil {
		return
	}

	if v.Minor, err = parseNumber(l); err != nil {
		return
	}
	if err = parseDot(l); err != nil {
		return
	}

	if v.Patch, err = parseNumber(l); err != nil {
		return
	}

	v.PreRelease, v.Build, err = parseSecondPart(l)

	return
}

func parseNumber(l *lexer) (n uint, err error) {
	i := <-l.items
	switch i.itemType {
	case itemEOF:
		return 0, unexpectedEOF
	case itemInvalid:
		return 0, invalidCharacter
	case itemNumber:
		n, err := strconv.ParseUint(i.value, 10, 64) // TODO: hard-code 64-bit here?
		return uint(n), err
	}
	return 0, expectedNumber

}

func parseDot(l *lexer) error {
	i := <-l.items
	switch i.itemType {
	case itemEOF:
		return unexpectedEOF
	case itemInvalid:
		return invalidCharacter
	case itemDot:
		return nil
	}
	return expectedDot
}

func parseIdentifier(l *lexer) (identifier []byte, err error) {
	i := <-l.items
	switch i.itemType {
	case itemEOF:
		return nil, unexpectedEOF
	case itemInvalid:
		return nil, invalidCharacter
	case itemIdentifier:
		return []byte(i.value), nil
	}
	return nil, expectedIdentifier
}

func parseSecondPart(l *lexer) (pre, build [][]byte, err error) {
	i := <-l.items
	switch i.itemType {
	case itemEOF:
		return
	case itemDash:
	parsepreloop:
		for {
			i := <-l.items
			switch i.itemType {
			case itemEOF:
				err = unexpectedEOF
				return
			case itemIdentifier:
				pre = append(pre, []byte(i.value))
				i = <-l.items
				switch i.itemType {
				case itemEOF:
					return
				case itemDot:
					continue parsepreloop
				case itemPlus:
					build, err = parseBuild(l)
					return
				default:
					break parsepreloop
				}
			default:
				err = expectedIdentifier
				return
			}

			if i.itemType == itemEOF {
				break
			}

			if i.itemType == itemDot {
				continue
			}

			break
		}
		err = expectedDotEOF
		return
	case itemPlus:
		build, err = parseBuild(l)
		return
	}

	err = expectedDashPlusEOF
	return
}

func parseBuild(l *lexer) (build [][]byte, err error) {
parsebuildloop:
	for {
		i := <-l.items
		switch i.itemType {
		case itemEOF:
			err = unexpectedEOF
			return
		case itemIdentifier:
			build = append(build, []byte(i.value))
			i = <-l.items
			switch i.itemType {
			case itemEOF:
				return
			case itemDot:
				continue parsebuildloop
			default:
				break parsebuildloop
			}
		default:
			err = expectedIdentifier
			return
		}

		if i.itemType == itemEOF {
			break
		}

		if i.itemType == itemDot {
			continue
		}

		break
	}
	err = expectedDotEOF
	return
}
