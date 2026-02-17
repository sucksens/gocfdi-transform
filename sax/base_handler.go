package sax

import (
	"encoding/xml"
	"fmt"
	"io"
)

// BaseHandler dara las funcionalidades comunes para todos los handlers
type BaseHandler struct {
	config HandlerConfig
}

// NewBaseHandler crea un nuevo BaseHandler
func NewBaseHandler(config HandlerConfig) *BaseHandler {
	return &BaseHandler{config: config}
}

// Config retorna la configuracion del handler
func (h *BaseHandler) Config() HandlerConfig {
	return h.config
}

// TransformFromBytes transforma un xml de bytes a String y llama al metodo TransformFromString
func (h *BaseHandler) TransformFromBytes(xmlBytes []byte) (interface{}, error) {
	return h.TransformFromString(string(xmlBytes))
}

// TransformFromString es un metodo abstracto que debe ser implementado por cada handler
func (h *BaseHandler) TransformFromString(xmlStr string) (interface{}, error) {
	return nil, nil
}

// ProcessElement procesa un elemento XML desde un decoder existente
// Este metodo es abstracto y debe ser implementado por cada handler
func (h *BaseHandler) ProcessElement(se xml.StartElement, decoder *xml.Decoder) (interface{}, error) {
	return nil, nil
}

// ParseTokens parsea los tokens de un xml hasta llegar al elemento endElement esperado
func (h *BaseHandler) ParseTokens(decoder *xml.Decoder, endElement string, tokenHandler func(token xml.Token) error) error {
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return fmt.Errorf("unexpected EOF while looking for end element: %s", endElement)
		}
		if err != nil {
			return err
		}
		if err := tokenHandler(token); err != nil {
			return err
		}
		if se, ok := token.(xml.EndElement); ok && se.Name.Local == endElement {
			return nil
		}
	}
}

// ValidateVersion valida que la version de un elemento Xml sea la esperada
func (h *BaseHandler) ValidateVersion(se xml.StartElement, expectedVersion string) error {
	version := getAttrValue(se, "version")

	if version == "" {
		version = getAttrValue(se, "Version")
	}
	if version == "" {
		return fmt.Errorf("missing required attribute: version")
	}
	if version != expectedVersion {
		return fmt.Errorf("incorrect version, expected %s, got %s", expectedVersion, version)
	}
	return nil
}
