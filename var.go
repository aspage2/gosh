package main

import (
	"os"
	"strings"
	"unicode"
)


const (
	BASE = 0
	ESC = 1
	PARSEVAR = 2
)

func IsVarChar(r rune) bool {
	return unicode.IsLetter(r) ||
		unicode.IsNumber(r) ||
		r == '_'

}

func GetVars(cmd string, vs VarSet) string {
	state := BASE
	var (
		retBuilder strings.Builder
		varBuilder strings.Builder
	)

	for _, chr := range cmd {
		switch state {
		case BASE:
			switch chr {
			case '$':
				state = PARSEVAR
			case '\\':
				state = ESC
			default:
				retBuilder.WriteRune(chr)
			}
		case ESC:
			retBuilder.WriteRune(chr)
			state = BASE
		case PARSEVAR:
			if IsVarChar(chr) {
				varBuilder.WriteRune(chr)
			} else {
				v := vs.Get(varBuilder.String())
				retBuilder.WriteString(v)
				varBuilder.Reset()
				retBuilder.WriteRune(chr)
				state = BASE
			}
		}
	}
	return retBuilder.String()
}


type VarSet interface {
	Get(string) string
	Set(string, string) error
}

type EnvVarSet struct{}

func (es EnvVarSet) Get(key string) string {
	return os.Getenv(key)
}

func (es EnvVarSet) Set(key string, val string) error {
	return os.Setenv(key, val)
}

type MapVarSet map[string]string

func (m MapVarSet) Get(key string) string {
	if v, ok := m[key]; ok {
		return v
	} else {
		return ""
	}
}

func (m MapVarSet) Set(key string, value string) error {
	m[key] = value
	return nil
}
