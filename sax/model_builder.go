package sax

import (
	"encoding/xml"
	"strings"

	"github.com/sucksens/gocfdi-transform/helpers"
)

// ModelBuilder proporciona métodos para construir modelos desde atributos XML.
type ModelBuilder struct {
	config HandlerConfig
}

// NewModelBuilder crea un nuevo ModelBuilder
func NewModelBuilder(cfg HandlerConfig) *ModelBuilder {
	return &ModelBuilder{config: cfg}
}

// ExtractString extrae un atributo como string
func (b *ModelBuilder) ExtractString(se xml.StartElement, attrName string) string {
	return getAttrValue(se, attrName)
}

// ExtractStringOrDefault extrae un atributo como string o retorna un valor por defecto
func (b *ModelBuilder) ExtractStringOrDefault(se xml.StartElement, attrName string) string {
	return getAttrValueOrDefault(se, attrName, b.config.EmptyChar)
}

// ExtractCompact extrae un atributo como string y lo compacta
func (b *ModelBuilder) ExtractCompact(se xml.StartElement, attrName string) string {
	return helpers.CompactString(b.config.EscDelimiters, getAttrValueOrDefault(se, attrName, b.config.EmptyChar))
}

// ExtractNumeric extrae un atributo como string y lo convierte a numerico
func (b *ModelBuilder) ExtractNumeric(se xml.StartElement, attrName string) string {
	return helpers.GetOrDefault(getAttrValue(se, attrName), b.config.EmptyChar, b.config.SafeNumerics)
}

// ExtractNumericOne extrae un atributo como string y lo convierte a numerico
func (b *ModelBuilder) ExtractNumericOne(se xml.StartElement, attrName string) string {
	return helpers.GetOrDefaultOne(getAttrValue(se, attrName), b.config.EmptyChar, b.config.SafeNumerics)
}

// ExtractUpper extrae un atributo como string y lo convierte a mayusculas
func (b *ModelBuilder) ExtractUpper(se xml.StartElement, attrName string) string {
	return strings.ToUpper(getAttrValue(se, attrName))
}

// ExtractTrim extrae un atributo como string y lo trimea
func (b *ModelBuilder) ExtractTrim(se xml.StartElement, attrName string) string {
	return strings.TrimSpace(getAttrValue(se, attrName))
}
