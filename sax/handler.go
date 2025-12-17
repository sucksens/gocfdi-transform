// Package sax contiene los manejadores SAX-like para parsear documentos XML CFDI.
package sax

import (
	"encoding/xml"

	"github.com/sucksens/gocfdi-transform/models"
)

// HandlerConfig contiene la configuración para el manejador SAX.
type HandlerConfig struct {
	EmptyChar             string
	SafeNumerics          bool
	EscDelimiters         string
	ParseConcepts         bool
	ParseRelatedCFDIs     bool
	ParseConceptsTaxes    bool
	ParsePagos20          bool
	ParseVentaVehiculos11 bool
}

// NewDefaultConfig retorna una configuración por defecto para el manejador SAX.
func NewDefaultConfig() HandlerConfig {
	return HandlerConfig{
		EmptyChar:             "",
		SafeNumerics:          false,
		EscDelimiters:         "",
		ParseConcepts:         false,
		ParseRelatedCFDIs:     false,
		ParseConceptsTaxes:    false,
		ParsePagos20:          false,
		ParseVentaVehiculos11: false,
	}
}

// Handler es una interfaz para manejadores CFDI.
type Handler interface {
	TransformFromFile(path string) (interface{}, error)
	TransformFromString(xml string) (interface{}, error)
}

type ComplementHandler interface {
	TransformFromBytes(xml []byte) (interface{}, error)
}

type ComplementFactory func(config HandlerConfig) ComplementHandler

type ComplementRegistry map[string]ComplementFactory

func DefaultCFDI40Complements() ComplementRegistry {
	return ComplementRegistry{
		"{http://www.sat.gob.mx/TimbreFiscalDigital}TimbreFiscalDigital": func(config HandlerConfig) ComplementHandler {
			return NewTFD11Handler(config)
		},
	}
}

// getAttrValue gets the value of an attribute from a xml.StartElement.
func getAttrValue(se xml.StartElement, name string) string {
	for _, attr := range se.Attr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

// getAttrValueOrDefault gets the value of an attribute or a default value if not found or empty.
func getAttrValueOrDefault(se xml.StartElement, name string, defaultValue string) string {
	val := getAttrValue(se, name)
	if val == "" {
		return defaultValue
	}
	return val
}

// initCFDI40Data creates a new CFDI40Data with default values.
func initCFDI40Data(config HandlerConfig) *models.CFDI40Data {
	emptyOrZero := config.EmptyChar
	emptyOrOne := config.EmptyChar
	if config.SafeNumerics {
		emptyOrZero = "0.00"
		emptyOrOne = "1.00"
	}

	return &models.CFDI40Data{
		CFDI40: models.CFDI40{
			Version:         config.EmptyChar,
			Serie:           config.EmptyChar,
			Folio:           config.EmptyChar,
			Fecha:           config.EmptyChar,
			NoCertificado:   config.EmptyChar,
			SubTotal:        emptyOrZero,
			Descuento:       emptyOrZero,
			Total:           emptyOrZero,
			Moneda:          config.EmptyChar,
			TipoCambio:      emptyOrOne,
			TipoComprobante: config.EmptyChar,
			MetodoPago:      config.EmptyChar,
			FormaPago:       config.EmptyChar,
			CondicionesPago: config.EmptyChar,
			LugarExpedicion: config.EmptyChar,
			Exportacion:     config.EmptyChar,
			Sello:           config.EmptyChar,
			Certificado:     config.EmptyChar,
			Confirmacion:    config.EmptyChar,
			Emisor: models.Emisor40{
				RFC:              config.EmptyChar,
				Nombre:           config.EmptyChar,
				RegimenFiscal:    config.EmptyChar,
				FacAtrAdquirente: config.EmptyChar,
			},
			Receptor: models.Receptor40{
				RFC:                     config.EmptyChar,
				Nombre:                  config.EmptyChar,
				DomicilioFiscalReceptor: config.EmptyChar,
				ResidenciaFiscal:        config.EmptyChar,
				NumRegIdTrib:            config.EmptyChar,
				RegimenFiscalReceptor:   config.EmptyChar,
				UsoCFDI:                 config.EmptyChar,
			},
			Conceptos: []models.Concepto40{},
			Impuestos: models.Impuestos{
				TotalImpuestosTrasladados: emptyOrZero,
				TotalImpuestosRetenidos:   emptyOrZero,
				Traslados:                 []models.Traslado{},
				Retenciones:               []models.Retencion{},
			},
			Complementos: config.EmptyChar,
			Addendas:     config.EmptyChar,
		},
		TFD11: []models.TFD11{},
	}
}
