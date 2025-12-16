// Package helpers nos da utilidades para manejar de manera segura strings y decimales.
package helpers

import (
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

// constantes para manejar decimales
const (
	DefaultSafeNumberZero = "0.00"
	DefaultSafeNumberOne  = "1.00"
)

// CompactString compacta una cadena eliminando los caracteres especificados y colapsando espacios en blanco.
// Tambien colapsa muchos espacios en blanco en uno solo y limpia los espacios en blanco al inicio y final.
func CompactString(delimiters, s string) string {
	if s == "" {
		return s
	}
	// characters a eliminar  \n, \t, \r, y  custom delimitadores
	charsToRemove := "\n\t\r" + delimiters

	// Eliminar los caracteres especificados de la cadena
	result := s
	for _, char := range charsToRemove {
		result = strings.ReplaceAll(result, string(char), "")
	}

	// Colapsar muchos espacios en blanco en uno solo
	result = regexp.MustCompile(`\s{2,}`).ReplaceAllString(result, " ")
	return strings.TrimSpace(result)
}

// SumStrings suma dos cadenas que representan decimales.
// Si alguna de las cadenas esta vacia, retorna la otra cadena.
// Si ambas cadenas estan vacias, retorna una cadena vacia.
func SumStrings(firstVal, secondVal string) string {
	if firstVal == "" {
		return secondVal
	}
	if secondVal == "" {
		return firstVal
	}

	first := TryParseDecimal(firstVal)
	second := TryParseDecimal(secondVal)
	return first.Add(second).String()
}

// TryParseDecimal intenta parsear una cadena en un decimal.Decimal.
// Si la cadena no es un decimal valido, retorna decimal.Zero.
func TryParseDecimal(val string) decimal.Decimal {
	d, err := decimal.NewFromString(val)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// GetOrDefault retorna el valor si no esta vacio, de lo contrario retorna el valor por defecto.
// Si safeNumerics es true, retorna DefaultSafeNumberZero en lugar de emptyChar.
func GetOrDefault(value, emptyChar string, safeNumerics bool) string {
	if value != "" {
		return value
	}
	if safeNumerics {
		return DefaultSafeNumberZero
	}
	return emptyChar
}

// GetOrDefaultOne retorna el valor si no esta vacio, de lo contrario retorna el valor por defecto.
// Si safeNumerics es true, retorna DefaultSafeNumberOne en lugar de emptyChar.
func GetOrDefaultOne(value, emptyChar string, safeNumerics bool) string {
	if value != "" {
		return value
	}
	if safeNumerics {
		return DefaultSafeNumberOne
	}
	return emptyChar
}

// FileToString lee el contenido de un archivo y lo retorna como una cadena.
// Si ocurre un error, retorna una cadena vacia y el error.
func FileToString(filePath string) (string, error) {
	// Este es un placeholder para la logica de lectura del archivo
	// La implementacion final se hara en el handler
	return "", nil
}
