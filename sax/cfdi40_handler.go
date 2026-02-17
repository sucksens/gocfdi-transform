package sax

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sucksens/gocfdi-transform/models"
)

// CFDI40Handler handles parsing of CFDI 4.0 XML documents.
type CFDI40Handler struct {
	*BaseHandler
	complements ComplementRegistry
	builder     *ModelBuilder
}

// NewCFDI40Handler creates a new CFDI40Handler with the given configuration.
func NewCFDI40Handler(config HandlerConfig) *CFDI40Handler {
	return &CFDI40Handler{
		BaseHandler: NewBaseHandler(config),
		complements: DefaultCFDI40Complements(),
		builder:     NewModelBuilder(config),
	}
}

// UseConcepts enables parsing of concepts.
func (h *CFDI40Handler) UseConcepts() *CFDI40Handler {
	h.config.ParseConcepts = true
	return h
}

// UseConceptsWithTaxes enables parsing of concepts with their taxes.
func (h *CFDI40Handler) UseConceptsWithTaxes() *CFDI40Handler {
	h.config.ParseConceptsTaxes = true
	return h
}

// UseRelatedCFDIs enables parsing of related CFDIs.
func (h *CFDI40Handler) UseRelatedCFDIs() *CFDI40Handler {
	h.config.ParseRelatedCFDIs = true
	return h
}

// UseNomina12 enables parsing of Nomina 1.2 complement.
func (h *CFDI40Handler) UseNomina12() *CFDI40Handler {
	h.config.ParseNomina12 = true
	return h
}

// UsePagos20 enables parsing of Pagos 2.0 complement.
func (h *CFDI40Handler) UsePagos20() *CFDI40Handler {
	h.config.ParsePagos20 = true
	return h
}

// UseVentaVehiculos11 enables parsing of Venta Vehículos 1.1 complement.
func (h *CFDI40Handler) UseVentaVehiculos11() *CFDI40Handler {
	h.config.ParseVentaVehiculos11 = true
	return h
}

// UseImpuestosLocales habilita el parsing del complemento de Impuestos Locales 1.0.
func (h *CFDI40Handler) UseImpuestosLocales() *CFDI40Handler {
	h.config.ParseImpuestosLocales = true
	return h
}

// TransformFromFile parses a CFDI 4.0 XML file.
func (h *CFDI40Handler) TransformFromFile(path string) (*models.CFDI40Data, error) {
	if !strings.HasSuffix(strings.ToLower(path), ".xml") {
		return nil, errors.New("incorrect type of document, only support XML files")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return h.TransformFromString(string(content))
}

// TransformFromString parses a CFDI 4.0 XML string.
func (h *CFDI40Handler) TransformFromString(xmlStr string) (*models.CFDI40Data, error) {
	data := initCFDI40Data(h.config)
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	var insideConcepts bool
	var currentConcept *models.Concepto40
	var complementNames []string

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing XML: %w", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "Comprobante":
				if err := h.transformComprobante(se, data); err != nil {
					return nil, err
				}

			case "Emisor":
				h.transformEmisor(se, data)

			case "Receptor":
				h.transformReceptor(se, data)

			case "Conceptos":
				insideConcepts = true
				if !h.config.ParseConcepts {
					// Skip the entire Conceptos subtree
					if err := decoder.Skip(); err != nil {
						return nil, err
					}
					insideConcepts = false
				}

			case "Concepto":
				if insideConcepts && h.config.ParseConcepts {
					currentConcept = h.transformConcepto(se)
				}

			case "Impuestos":
				if !insideConcepts {
					h.transformImpuestos(se, decoder, data)
				} else if h.config.ParseConcepts && h.config.ParseConceptsTaxes && currentConcept != nil {
					h.transformImpuestosConcepto(se, decoder, currentConcept)
				}

			case "CfdiRelacionados":
				if h.config.ParseRelatedCFDIs {
					h.transformCFDIsRelacionados(se, decoder, data)
				}

			case "Complemento":
				h.transformComplemento(decoder, data, &complementNames)

			case "Addenda":
				h.transformAddenda(decoder, data)
			}

		case xml.EndElement:
			switch se.Name.Local {
			case "Conceptos":
				insideConcepts = false

			case "Concepto":
				if currentConcept != nil {
					data.CFDI40.Conceptos = append(data.CFDI40.Conceptos, *currentConcept)
					currentConcept = nil
				}
			}
		}
	}

	if len(complementNames) > 0 {
		data.CFDI40.Complementos = strings.Join(complementNames, " ")
	}

	return data, nil
}

func (h *CFDI40Handler) transformComprobante(se xml.StartElement, data *models.CFDI40Data) error {
	if err := h.ValidateVersion(se, "4.0"); err != nil {
		return errors.New("incorrect type of CFDI, this handler only supports CFDI version 4.0")
	}

	data.CFDI40.Version = "4.0"
	data.CFDI40.Serie = h.builder.ExtractCompact(se, "Serie")
	data.CFDI40.Folio = h.builder.ExtractCompact(se, "Folio")
	data.CFDI40.Fecha = h.builder.ExtractString(se, "Fecha")
	data.CFDI40.NoCertificado = h.builder.ExtractString(se, "NoCertificado")
	data.CFDI40.SubTotal = h.builder.ExtractString(se, "SubTotal")
	data.CFDI40.Descuento = h.builder.ExtractNumeric(se, "Descuento")
	data.CFDI40.Total = h.builder.ExtractString(se, "Total")
	data.CFDI40.Moneda = h.builder.ExtractString(se, "Moneda")
	data.CFDI40.TipoCambio = h.builder.ExtractNumericOne(se, "TipoCambio")
	data.CFDI40.TipoComprobante = h.builder.ExtractString(se, "TipoDeComprobante")
	data.CFDI40.MetodoPago = h.builder.ExtractCompact(se, "MetodoPago")
	data.CFDI40.FormaPago = h.builder.ExtractCompact(se, "FormaPago")
	data.CFDI40.CondicionesPago = h.builder.ExtractCompact(se, "CondicionesDePago")
	data.CFDI40.LugarExpedicion = h.builder.ExtractString(se, "LugarExpedicion")
	data.CFDI40.Exportacion = h.builder.ExtractString(se, "Exportacion")
	data.CFDI40.Sello = h.builder.ExtractCompact(se, "Sello")
	data.CFDI40.Certificado = h.builder.ExtractCompact(se, "Certificado")
	data.CFDI40.Confirmacion = h.builder.ExtractStringOrDefault(se, "Confirmacion")

	return nil
}

func (h *CFDI40Handler) transformEmisor(se xml.StartElement, data *models.CFDI40Data) {
	data.CFDI40.Emisor.RFC = h.builder.ExtractString(se, "Rfc")
	data.CFDI40.Emisor.Nombre = h.builder.ExtractCompact(se, "Nombre")
	data.CFDI40.Emisor.RegimenFiscal = h.builder.ExtractString(se, "RegimenFiscal")
	data.CFDI40.Emisor.FacAtrAdquirente = h.builder.ExtractCompact(se, "FacAtrAdquirente")
}

func (h *CFDI40Handler) transformReceptor(se xml.StartElement, data *models.CFDI40Data) {
	data.CFDI40.Receptor.RFC = h.builder.ExtractString(se, "Rfc")
	data.CFDI40.Receptor.Nombre = h.builder.ExtractCompact(se, "Nombre")
	data.CFDI40.Receptor.DomicilioFiscalReceptor = h.builder.ExtractString(se, "DomicilioFiscalReceptor")
	data.CFDI40.Receptor.ResidenciaFiscal = h.builder.ExtractStringOrDefault(se, "ResidenciaFiscal")
	data.CFDI40.Receptor.NumRegIdTrib = h.builder.ExtractCompact(se, "NumRegIdTrib")
	data.CFDI40.Receptor.RegimenFiscalReceptor = h.builder.ExtractString(se, "RegimenFiscalReceptor")
	data.CFDI40.Receptor.UsoCFDI = h.builder.ExtractString(se, "UsoCFDI")
}

func (h *CFDI40Handler) transformConcepto(se xml.StartElement) *models.Concepto40 {
	return &models.Concepto40{
		ClaveProdServ:    h.builder.ExtractString(se, "ClaveProdServ"),
		NoIdentificacion: h.builder.ExtractCompact(se, "NoIdentificacion"),
		Cantidad:         h.builder.ExtractString(se, "Cantidad"),
		ClaveUnidad:      h.builder.ExtractString(se, "ClaveUnidad"),
		Unidad:           h.builder.ExtractCompact(se, "Unidad"),
		Descripcion:      h.builder.ExtractCompact(se, "Descripcion"),
		ValorUnitario:    h.builder.ExtractString(se, "ValorUnitario"),
		Importe:          h.builder.ExtractString(se, "Importe"),
		Descuento:        h.builder.ExtractNumeric(se, "Descuento"),
		ObjetoImp:        h.builder.ExtractString(se, "ObjetoImp"),
	}
}

func (h *CFDI40Handler) transformImpuestos(se xml.StartElement, decoder *xml.Decoder, data *models.CFDI40Data) {
	data.CFDI40.Impuestos.TotalImpuestosTrasladados = h.builder.ExtractNumeric(se, "TotalImpuestosTrasladados")
	data.CFDI40.Impuestos.TotalImpuestosRetenidos = h.builder.ExtractNumeric(se, "TotalImpuestosRetenidos")

	ProcessChildElements(decoder, "Impuestos", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "Traslado":
			traslado := models.Traslado{
				Base:       h.builder.ExtractNumeric(childSE, "Base"),
				Impuesto:   h.builder.ExtractCompact(childSE, "Impuesto"),
				TipoFactor: h.builder.ExtractCompact(childSE, "TipoFactor"),
				TasaOCuota: h.builder.ExtractNumeric(childSE, "TasaOCuota"),
				Importe:    h.builder.ExtractNumeric(childSE, "Importe"),
			}
			data.CFDI40.Impuestos.Traslados = append(data.CFDI40.Impuestos.Traslados, traslado)

		case "Retencion":
			retencion := models.Retencion{
				Impuesto: h.builder.ExtractCompact(childSE, "Impuesto"),
				Importe:  h.builder.ExtractNumeric(childSE, "Importe"),
			}
			data.CFDI40.Impuestos.Retenciones = append(data.CFDI40.Impuestos.Retenciones, retencion)
		}
		return nil
	})
}

func (h *CFDI40Handler) transformImpuestosConcepto(se xml.StartElement, decoder *xml.Decoder, concept *models.Concepto40) {
	ProcessChildElements(decoder, "Impuestos", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "Traslado":
			traslado := models.TrasladoConcepto{
				Base:       h.builder.ExtractNumeric(childSE, "Base"),
				Impuesto:   h.builder.ExtractCompact(childSE, "Impuesto"),
				TipoFactor: h.builder.ExtractCompact(childSE, "TipoFactor"),
				TasaOCuota: h.builder.ExtractNumeric(childSE, "TasaOCuota"),
				Importe:    h.builder.ExtractNumeric(childSE, "Importe"),
			}
			concept.Traslados = append(concept.Traslados, traslado)

		case "Retencion":
			retencion := models.RetencionConcepto{
				Impuesto: h.builder.ExtractCompact(childSE, "Impuesto"),
				Importe:  h.builder.ExtractNumeric(childSE, "Importe"),
			}
			concept.Retenciones = append(concept.Retenciones, retencion)
		}
		return nil
	})
}

func (h *CFDI40Handler) transformCFDIsRelacionados(se xml.StartElement, decoder *xml.Decoder, data *models.CFDI40Data) {
	tipoRelacion := h.builder.ExtractString(se, "TipoRelacion")

	ProcessChildElements(decoder, "CfdiRelacionados", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		if childSE.Name.Local == "CfdiRelacionado" {
			cfdiRel := models.CFDIRelacionado{
				UUID:         strings.ToUpper(h.builder.ExtractString(childSE, "UUID")),
				TipoRelacion: tipoRelacion,
			}
			data.CFDI40.CFDIsRelacionados = append(data.CFDI40.CFDIsRelacionados, cfdiRel)
		}
		return nil
	})
}

func (h *CFDI40Handler) transformComplemento(decoder *xml.Decoder, data *models.CFDI40Data, complementNames *[]string) {
	ProcessChildElements(decoder, "Complemento", func(se xml.StartElement, childDecoder *xml.Decoder) error {
		// Record complement name
		*complementNames = append(*complementNames, se.Name.Local)

		// Handle TFD11
		if se.Name.Local == "TimbreFiscalDigital" && se.Name.Space == "http://www.sat.gob.mx/TimbreFiscalDigital" {
			tfd := models.TFD11{
				Version:          h.builder.ExtractString(se, "Version"),
				NoCertificadoSAT: h.builder.ExtractString(se, "NoCertificadoSAT"),
				UUID:             strings.ToUpper(h.builder.ExtractString(se, "UUID")),
				FechaTimbrado:    h.builder.ExtractString(se, "FechaTimbrado"),
				RfcProvCert:      h.builder.ExtractString(se, "RfcProvCertif"),
				SelloCFD:         h.builder.ExtractCompact(se, "SelloCFD"),
				SelloSAT:         h.builder.ExtractCompact(se, "SelloSAT"),
			}
			data.TFD11 = append(data.TFD11, tfd)
		}

		// Handle Nomina 1.2
		if h.config.ParseNomina12 && se.Name.Local == "Nomina" && se.Name.Space == "http://www.sat.gob.mx/nomina12" {
			nomina12Handler := NewNomina12Handler(h.config)
			nomina12Data, err := nomina12Handler.ProcessNomina12Element(se, childDecoder)
			if err == nil && nomina12Data != nil {
				data.Nomina12 = append(data.Nomina12, *nomina12Data)
			}
		}

		// Handle Pagos 2.0
		if h.config.ParsePagos20 && se.Name.Local == "Pagos" && se.Name.Space == "http://www.sat.gob.mx/Pagos20" {
			pagosHandler := NewPagos20Handler(h.config)
			pagosData, err := pagosHandler.ProcessPagosElement(se, childDecoder)
			if err == nil && pagosData != nil {
				data.Pagos20 = append(data.Pagos20, *pagosData)
			}
		}

		// Handle VentaVehiculos 1.1
		if h.config.ParseVentaVehiculos11 && se.Name.Local == "VentaVehiculos" && se.Name.Space == "http://www.sat.gob.mx/ventavehiculos" {
			ventaVehiculos11Handler := NewVentaVehiculos11Handler(h.config)
			ventaVehiculos11Data, err := ventaVehiculos11Handler.ProcessVentaVehiculosElement(se, childDecoder)
			if err == nil && ventaVehiculos11Data != nil {
				data.VentaVehiculos11 = append(data.VentaVehiculos11, *ventaVehiculos11Data)
			}
		}

		// Handle ImpuestosLocales 1.0
		if h.config.ParseImpuestosLocales && se.Name.Local == "ImpuestosLocales" && se.Name.Space == "http://www.sat.gob.mx/implocal" {
			impLocHandler := NewImpuestosLocales10Handler(h.config)
			impLocData, err := impLocHandler.ProcessImpuestosLocalesElement(se, childDecoder)
			if err == nil && impLocData != nil {
				data.ImpuestosLocales = append(data.ImpuestosLocales, *impLocData)
			}
		}
		return nil
	})
}

func (h *CFDI40Handler) transformAddenda(decoder *xml.Decoder, data *models.CFDI40Data) {
	var addendaNames []string

	ProcessChildElements(decoder, "Addenda", func(se xml.StartElement, childDecoder *xml.Decoder) error {
		addendaNames = append(addendaNames, se.Name.Local)
		return nil
	})

	if len(addendaNames) > 0 {
		data.CFDI40.Addendas = strings.Join(addendaNames, " ")
	}
}
