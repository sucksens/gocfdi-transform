package sax

import (
	"encoding/xml"
	"fmt"
	"io"
)

// ElementProcessor es una funcion que procesa un elemento StartElement
// Recibe el elemento y el decoder XML, permitiendo procesamiento personalizado.
// Si retorna un error, el procesamiento se detiene.
type ElementProcessor func(se xml.StartElement, decoder *xml.Decoder)

// ProcessElementUntil procesa un elemento StartElement hasta llegar al elemento endElement esperado
func ProcessElementUntil(decoder *xml.Decoder, endElement string, processor ElementProcessor) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("unexpected EOF while looking for end element: %s", endElement)
			}
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if err := processor(t, decoder); err != nil {
				return err
			}
		case xml.EndElement:
			if t.Name.Local == endElement {
				return nil
			}
		}
	}
}

// ProcessChildElements procesa todos los elementos hijos hasta el EndElement del elemento padre
func ProcessChildElements(decoder *xml.Decoder, parentElement string, processor ElementProcessor) error {
	return ProcessElementUntil(decoder, parentElement, processor)
}

// FindElement encuentra el elemento StartElement con el nombre especificado
func FindElement(decoder *xml.Decoder, elementName string) (*xml.StartElement, error) {
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		if se, ok := token.(xml.StartElement); ok && se.Name.Local == elementName {
			return &se, nil
		}
	}
}

// SkipElement salta todo el contenido de un elemento XML, incluyendo todos sus elementos hijos,
// hasta encontrar el elemento de cierre correspondiente.
// Útil para ignorar elementos que no necesitan ser procesados
func SkipElement(decoder *xml.Decoder, elementName string) error {
	depth := 1
	for depth > 0 {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch token.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
	}
	return nil
}
