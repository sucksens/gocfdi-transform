package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/helpers"
	"github.com/sucksens/gocfdi-transform/models"
)

type ImpuestosLocales10Handler struct {
	// config contiene la configuración del handler, incluyendo flags como SafeNumerics y EscDelimiters.
	config HandlerConfig
}

// NewImpuestosLocales10Handler crea una nueva instancia del handler de Impuestos Locales 1.0.
func NewImpuestosLocales10Handler(cfg HandlerConfig) *ImpuestosLocales10Handler {
	return &ImpuestosLocales10Handler{config: cfg}
}

// ProcessImpuestosLocalesElement procesa el elemento ImpuestosLocales desde un stream XML existente.
func (h *ImpuestosLocales10Handler) ProcessImpuestosLocalesElement(se xml.StartElement, decoder *xml.Decoder) (*models.ImpuestosLocales10, error) {
	// Validar que la versión sea "1.0" (fija según especificación del SAT)
	version := strings.TrimSpace(getAttrValue(se, "version"))
	if version != "1.0" {
		return nil, errors.New("incorrect type of ImpuestosLocales, this handler only supports version 1.0")
	}

	// Inicializar estructura de datos con valores extraídos del elemento principal
	data := &models.ImpuestosLocales10{
		Version:            version,
		TotaldeRetenciones: helpers.GetOrDefault(getAttrValue(se, "TotaldeRetenciones"), h.config.EmptyChar, h.config.SafeNumerics),
		TotaldeTraslados:   helpers.GetOrDefault(getAttrValue(se, "TotaldeTraslados"), h.config.EmptyChar, h.config.SafeNumerics),
		RetencionesLocales: []models.RetencionLocal10{},
		TrasladosLocales:   []models.TrasladoLocal10{},
	}

	// Iterar sobre los elementos hijos (RetencionesLocales y TrasladosLocales)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "RetencionesLocales":
				// Procesar retención de impuesto local
				retencion := h.transformRetencionLocalElement(t)
				data.RetencionesLocales = append(data.RetencionesLocales, retencion)

			case "TrasladosLocales":
				// Procesar traslado de impuesto local
				traslado := h.transformTrasladoLocalElement(t)
				data.TrasladosLocales = append(data.TrasladosLocales, traslado)
			}

		case xml.EndElement:
			// Terminar el procesamiento cuando se cierra el elemento ImpuestosLocales
			if t.Name.Local == "ImpuestosLocales" {
				return data, nil
			}
		}
	}

	return data, nil
}

// TransformFromBytes parsea una cadena de bytes que contiene un XML de ImpuestosLocales.
func (h *ImpuestosLocales10Handler) TransformFromBytes(xmlBytes []byte) (*models.ImpuestosLocales10, error) {
	return h.TransformFromString(string(xmlBytes))
}

// TransformFromString parsea una cadena XML que contiene un complemento de ImpuestosLocales.
func (h *ImpuestosLocales10Handler) TransformFromString(xmlStr string) (*models.ImpuestosLocales10, error) {
	// Inicializar estructura de datos
	data := &models.ImpuestosLocales10{
		RetencionesLocales: []models.RetencionLocal10{},
		TrasladosLocales:   []models.TrasladoLocal10{},
	}

	// Crear decoder XML para parsear la cadena
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	// Iterar sobre los tokens del XML
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "ImpuestosLocales":
				// Encontrar el elemento principal y extraer sus atributos
				version := strings.TrimSpace(getAttrValue(se, "version"))
				if version != "1.0" {
					return nil, errors.New("incorrect type of ImpuestosLocales, this handler only supports version 1.0")
				}
				data.Version = version
				data.TotaldeRetenciones = helpers.GetOrDefault(getAttrValue(se, "TotaldeRetenciones"), h.config.EmptyChar, h.config.SafeNumerics)
				data.TotaldeTraslados = helpers.GetOrDefault(getAttrValue(se, "TotaldeTraslados"), h.config.EmptyChar, h.config.SafeNumerics)

			case "RetencionesLocales":
				// Procesar retención de impuesto local
				retencion := h.transformRetencionLocalElement(se)
				data.RetencionesLocales = append(data.RetencionesLocales, retencion)

			case "TrasladosLocales":
				// Procesar traslado de impuesto local
				traslado := h.transformTrasladoLocalElement(se)
				data.TrasladosLocales = append(data.TrasladosLocales, traslado)
			}
		}
	}

	return data, nil
}

// transformRetencionLocalElement extrae los atributos de un elemento RetencionesLocales
func (h *ImpuestosLocales10Handler) transformRetencionLocalElement(se xml.StartElement) models.RetencionLocal10 {
	return models.RetencionLocal10{
		// CompactString elimina caracteres especiales (\n, \t, \r) y colapsa espacios múltiples
		ImpLocRetenido: helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "ImpLocRetenido")),
		// GetOrDefault retorna el valor o "0.00" si SafeNumerics=true y el valor está vacío
		TasadeRetencion: helpers.GetOrDefault(getAttrValue(se, "TasadeRetencion"), h.config.EmptyChar, h.config.SafeNumerics),
		// GetOrDefault retorna el valor o "0.00" si SafeNumerics=true y el valor está vacío
		Importe: helpers.GetOrDefault(getAttrValue(se, "Importe"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

// transformTrasladoLocalElement extrae los atributos de un elemento TrasladosLocales
func (h *ImpuestosLocales10Handler) transformTrasladoLocalElement(se xml.StartElement) models.TrasladoLocal10 {
	return models.TrasladoLocal10{
		// CompactString elimina caracteres especiales (\n, \t, \r) y colapsa espacios múltiples
		ImpLocTrasladado: helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "ImpLocTrasladado")),
		// GetOrDefault retorna el valor o "0.00" si SafeNumerics=true y el valor está vacío
		TasadeTraslado: helpers.GetOrDefault(getAttrValue(se, "TasadeTraslado"), h.config.EmptyChar, h.config.SafeNumerics),
		// GetOrDefault retorna el valor o "0.00" si SafeNumerics=true y el valor está vacío
		Importe: helpers.GetOrDefault(getAttrValue(se, "Importe"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}
