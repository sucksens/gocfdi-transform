package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/models"
)

type ImpuestosLocales10Handler struct {
	// config contiene la configuración del handler, incluyendo flags como SafeNumerics y EscDelimiters.
	*BaseHandler
	builder *ModelBuilder
}

// NewImpuestosLocales10Handler crea una nueva instancia del handler de Impuestos Locales 1.0.
func NewImpuestosLocales10Handler(config HandlerConfig) *ImpuestosLocales10Handler {
	return &ImpuestosLocales10Handler{
		BaseHandler: NewBaseHandler(config),
		builder:     NewModelBuilder(config),
	}
}

// ProcessImpuestosLocalesElement procesa el elemento ImpuestosLocales desde un stream XML existente.
func (h *ImpuestosLocales10Handler) ProcessImpuestosLocalesElement(se xml.StartElement, decoder *xml.Decoder) (*models.ImpuestosLocales10, error) {
	// Validar que la versión sea "1.0" (fija según especificación del SAT)
	if err := h.ValidateVersion(se, "1.0"); err != nil {
		return nil, errors.New("incorrect type of ImpuestosLocales, this handler only supports version 1.0")
	}

	// Inicializar estructura de datos con valores extraídos del elemento principal
	data := &models.ImpuestosLocales10{
		Version:            "1.0",
		TotaldeRetenciones: h.builder.ExtractNumeric(se, "TotaldeRetenciones"),
		TotaldeTraslados:   h.builder.ExtractNumeric(se, "TotaldeTraslados"),
		RetencionesLocales: []models.RetencionLocal10{},
		TrasladosLocales:   []models.TrasladoLocal10{},
	}

	// Iterar sobre los elementos hijos (RetencionesLocales y TrasladosLocales)
	err := ProcessChildElements(decoder, "ImpuestosLocales", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "RetencionesLocales":
			// Procesar retención de impuesto local
			retencion := h.transformRetencionLocalElement(childSE)
			data.RetencionesLocales = append(data.RetencionesLocales, retencion)

		case "TrasladosLocales":
			// Procesar traslado de impuesto local
			traslado := h.transformTrasladoLocalElement(childSE)
			data.TrasladosLocales = append(data.TrasladosLocales, traslado)
		}
		return nil
	})

	return data, err
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

	// Buscar el elemento principal
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "ImpuestosLocales" {
			// Validar versión
			if err := h.ValidateVersion(se, "1.0"); err != nil {
				return nil, errors.New("incorrect type of ImpuestosLocales, this handler only supports version 1.0")
			}

			// Extraer atributos del elemento principal
			data.Version = "1.0"
			data.TotaldeRetenciones = h.builder.ExtractNumeric(se, "TotaldeRetenciones")
			data.TotaldeTraslados = h.builder.ExtractNumeric(se, "TotaldeTraslados")

			err := ProcessChildElements(decoder, "ImpuestosLocales", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
				switch childSE.Name.Local {
				case "RetencionesLocales":
					// Procesar retención de impuesto local
					retencion := h.transformRetencionLocalElement(childSE)
					data.RetencionesLocales = append(data.RetencionesLocales, retencion)

				case "TrasladosLocales":
					// Procesar traslado de impuesto local
					traslado := h.transformTrasladoLocalElement(childSE)
					data.TrasladosLocales = append(data.TrasladosLocales, traslado)
				}
				return nil
			})
			return data, err
		}
	}
	return data, nil
}

// transformRetencionLocalElement extrae los atributos de un elemento RetencionesLocales
func (h *ImpuestosLocales10Handler) transformRetencionLocalElement(se xml.StartElement) models.RetencionLocal10 {
	return models.RetencionLocal10{
		ImpLocRetenido:  h.builder.ExtractCompact(se, "ImpLocRetenido"),
		TasadeRetencion: h.builder.ExtractNumeric(se, "TasadeRetencion"),
		Importe:         h.builder.ExtractNumeric(se, "Importe"),
	}
}

// transformTrasladoLocalElement extrae los atributos de un elemento TrasladosLocales
func (h *ImpuestosLocales10Handler) transformTrasladoLocalElement(se xml.StartElement) models.TrasladoLocal10 {
	return models.TrasladoLocal10{
		ImpLocTrasladado: h.builder.ExtractCompact(se, "ImpLocTrasladado"),
		TasadeTraslado:   h.builder.ExtractNumeric(se, "TasadeTraslado"),
		Importe:          h.builder.ExtractNumeric(se, "Importe"),
	}
}
