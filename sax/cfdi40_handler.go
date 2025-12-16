package sax

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sucksens/gocfdi-transform/helpers"
	"github.com/sucksens/gocfdi-transform/models"
)

// CFDI40Handler handles parsing of CFDI 4.0 XML documents.
type CFDI40Handler struct {
	config      HandlerConfig
	complements ComplementRegistry
}

// NewCFDI40Handler creates a new CFDI40Handler with the given configuration.
func NewCFDI40Handler(cfg HandlerConfig) *CFDI40Handler {
	return &CFDI40Handler{
		config:      cfg,
		complements: DefaultCFDI40Complements(),
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

// UsePagos20 enables parsing of Pagos 2.0 complement.
func (h *CFDI40Handler) UsePagos20() *CFDI40Handler {
	h.config.ParsePagos20 = true
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
	version := getAttrValue(se, "Version")
	if version != "4.0" {
		return errors.New("incorrect type of CFDI, this handler only supports CFDI version 4.0")
	}

	data.CFDI40.Version = version
	data.CFDI40.Serie = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Serie", h.config.EmptyChar))
	data.CFDI40.Folio = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Folio", h.config.EmptyChar))
	data.CFDI40.Fecha = getAttrValue(se, "Fecha")
	data.CFDI40.NoCertificado = getAttrValue(se, "NoCertificado")
	data.CFDI40.SubTotal = getAttrValue(se, "SubTotal")
	data.CFDI40.Descuento = helpers.GetOrDefault(getAttrValue(se, "Descuento"), h.config.EmptyChar, h.config.SafeNumerics)
	data.CFDI40.Total = getAttrValue(se, "Total")
	data.CFDI40.Moneda = getAttrValue(se, "Moneda")
	data.CFDI40.TipoCambio = helpers.GetOrDefaultOne(getAttrValue(se, "TipoCambio"), h.config.EmptyChar, h.config.SafeNumerics)
	data.CFDI40.TipoComprobante = getAttrValue(se, "TipoDeComprobante")
	data.CFDI40.MetodoPago = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "MetodoPago", h.config.EmptyChar))
	data.CFDI40.FormaPago = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "FormaPago", h.config.EmptyChar))
	data.CFDI40.CondicionesPago = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "CondicionesDePago", h.config.EmptyChar))
	data.CFDI40.LugarExpedicion = getAttrValue(se, "LugarExpedicion")
	data.CFDI40.Exportacion = getAttrValue(se, "Exportacion")
	data.CFDI40.Sello = helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "Sello"))
	data.CFDI40.Certificado = helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "Certificado"))
	data.CFDI40.Confirmacion = getAttrValueOrDefault(se, "Confirmacion", h.config.EmptyChar)

	return nil
}

func (h *CFDI40Handler) transformEmisor(se xml.StartElement, data *models.CFDI40Data) {
	data.CFDI40.Emisor.RFC = getAttrValue(se, "Rfc")
	data.CFDI40.Emisor.Nombre = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Nombre", h.config.EmptyChar))
	data.CFDI40.Emisor.RegimenFiscal = getAttrValue(se, "RegimenFiscal")
	data.CFDI40.Emisor.FacAtrAdquirente = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "FacAtrAdquirente", h.config.EmptyChar))
}

func (h *CFDI40Handler) transformReceptor(se xml.StartElement, data *models.CFDI40Data) {
	data.CFDI40.Receptor.RFC = getAttrValue(se, "Rfc")
	data.CFDI40.Receptor.Nombre = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Nombre", h.config.EmptyChar))
	data.CFDI40.Receptor.DomicilioFiscalReceptor = getAttrValue(se, "DomicilioFiscalReceptor")
	data.CFDI40.Receptor.ResidenciaFiscal = getAttrValueOrDefault(se, "ResidenciaFiscal", h.config.EmptyChar)
	data.CFDI40.Receptor.NumRegIdTrib = helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "NumRegIdTrib", h.config.EmptyChar))
	data.CFDI40.Receptor.RegimenFiscalReceptor = getAttrValue(se, "RegimenFiscalReceptor")
	data.CFDI40.Receptor.UsoCFDI = getAttrValue(se, "UsoCFDI")
}

func (h *CFDI40Handler) transformConcepto(se xml.StartElement) *models.Concepto40 {
	return &models.Concepto40{
		ClaveProdServ:    getAttrValue(se, "ClaveProdServ"),
		NoIdentificacion: helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "NoIdentificacion", h.config.EmptyChar)),
		Cantidad:         getAttrValue(se, "Cantidad"),
		ClaveUnidad:      getAttrValue(se, "ClaveUnidad"),
		Unidad:           helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Unidad", h.config.EmptyChar)),
		Descripcion:      helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "Descripcion")),
		ValorUnitario:    getAttrValue(se, "ValorUnitario"),
		Importe:          getAttrValue(se, "Importe"),
		Descuento:        helpers.GetOrDefault(getAttrValue(se, "Descuento"), h.config.EmptyChar, h.config.SafeNumerics),
		ObjetoImp:        getAttrValue(se, "ObjetoImp"),
	}
}

func (h *CFDI40Handler) transformImpuestos(se xml.StartElement, decoder *xml.Decoder, data *models.CFDI40Data) {
	data.CFDI40.Impuestos.TotalImpuestosTrasladados = helpers.GetOrDefault(getAttrValue(se, "TotalImpuestosTrasladados"), h.config.EmptyChar, h.config.SafeNumerics)
	data.CFDI40.Impuestos.TotalImpuestosRetenidos = helpers.GetOrDefault(getAttrValue(se, "TotalImpuestosRetenidos"), h.config.EmptyChar, h.config.SafeNumerics)

	// Parse child elements
	for {
		token, err := decoder.Token()
		if err != nil {
			return
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "Traslado":
				traslado := models.Traslado{
					Base:       helpers.GetOrDefault(getAttrValue(t, "Base"), h.config.EmptyChar, h.config.SafeNumerics),
					Impuesto:   helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(t, "Impuesto", h.config.EmptyChar)),
					TipoFactor: helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(t, "TipoFactor", h.config.EmptyChar)),
					TasaOCuota: helpers.GetOrDefault(getAttrValue(t, "TasaOCuota"), h.config.EmptyChar, h.config.SafeNumerics),
					Importe:    helpers.GetOrDefault(getAttrValue(t, "Importe"), h.config.EmptyChar, h.config.SafeNumerics),
				}
				data.CFDI40.Impuestos.Traslados = append(data.CFDI40.Impuestos.Traslados, traslado)

			case "Retencion":
				retencion := models.Retencion{
					Impuesto: helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(t, "Impuesto", h.config.EmptyChar)),
					Importe:  helpers.GetOrDefault(getAttrValue(t, "Importe"), h.config.EmptyChar, h.config.SafeNumerics),
				}
				data.CFDI40.Impuestos.Retenciones = append(data.CFDI40.Impuestos.Retenciones, retencion)
			}

		case xml.EndElement:
			if t.Name.Local == "Impuestos" {
				return
			}
		}
	}
}

func (h *CFDI40Handler) transformImpuestosConcepto(se xml.StartElement, decoder *xml.Decoder, concept *models.Concepto40) {
	for {
		token, err := decoder.Token()
		if err != nil {
			return
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "Traslado":
				traslado := models.TrasladoConcepto{
					Base:       helpers.GetOrDefault(getAttrValue(t, "Base"), h.config.EmptyChar, h.config.SafeNumerics),
					Impuesto:   helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(t, "Impuesto", h.config.EmptyChar)),
					TipoFactor: helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(t, "TipoFactor", h.config.EmptyChar)),
					TasaOCuota: helpers.GetOrDefault(getAttrValue(t, "TasaOCuota"), h.config.EmptyChar, h.config.SafeNumerics),
					Importe:    helpers.GetOrDefault(getAttrValue(t, "Importe"), h.config.EmptyChar, h.config.SafeNumerics),
				}
				concept.Traslados = append(concept.Traslados, traslado)

			case "Retencion":
				retencion := models.RetencionConcepto{
					Impuesto: helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(t, "Impuesto", h.config.EmptyChar)),
					Importe:  helpers.GetOrDefault(getAttrValue(t, "Importe"), h.config.EmptyChar, h.config.SafeNumerics),
				}
				concept.Retenciones = append(concept.Retenciones, retencion)
			}

		case xml.EndElement:
			if t.Name.Local == "Impuestos" {
				return
			}
		}
	}
}

func (h *CFDI40Handler) transformCFDIsRelacionados(se xml.StartElement, decoder *xml.Decoder, data *models.CFDI40Data) {
	tipoRelacion := getAttrValue(se, "TipoRelacion")

	for {
		token, err := decoder.Token()
		if err != nil {
			return
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "CfdiRelacionado" {
				cfdiRel := models.CFDIRelacionado{
					UUID:         strings.ToUpper(getAttrValue(t, "UUID")),
					TipoRelacion: tipoRelacion,
				}
				data.CFDI40.CFDIsRelacionados = append(data.CFDI40.CFDIsRelacionados, cfdiRel)
			}

		case xml.EndElement:
			if t.Name.Local == "CfdiRelacionados" {
				return
			}
		}
	}
}

func (h *CFDI40Handler) transformComplemento(decoder *xml.Decoder, data *models.CFDI40Data, complementNames *[]string) {
	for {
		token, err := decoder.Token()
		if err != nil {
			return
		}

		switch t := token.(type) {
		case xml.StartElement:
			// Record complement name
			*complementNames = append(*complementNames, t.Name.Local)

			// Handle TFD11
			if t.Name.Local == "TimbreFiscalDigital" && t.Name.Space == "http://www.sat.gob.mx/TimbreFiscalDigital" {
				tfd := models.TFD11{
					Version:          getAttrValue(t, "Version"),
					NoCertificadoSAT: getAttrValue(t, "NoCertificadoSAT"),
					UUID:             strings.ToUpper(getAttrValue(t, "UUID")),
					FechaTimbrado:    getAttrValue(t, "FechaTimbrado"),
					RfcProvCert:      getAttrValue(t, "RfcProvCertif"),
					SelloCFD:         helpers.CompactString(h.config.EscDelimiters, getAttrValue(t, "SelloCFD")),
					SelloSAT:         helpers.CompactString(h.config.EscDelimiters, getAttrValue(t, "SelloSAT")),
				}
				data.TFD11 = append(data.TFD11, tfd)
			}

			// Handle Pagos 2.0
			if h.config.ParsePagos20 && t.Name.Local == "Pagos" && t.Name.Space == "http://www.sat.gob.mx/Pagos20" {
				pagosHandler := NewPagos20Handler(h.config)
				pagosData, err := pagosHandler.ProcessPagosElement(t, decoder)
				if err == nil && pagosData != nil {
					data.Pagos20 = append(data.Pagos20, *pagosData)
				}
			}

		case xml.EndElement:
			if t.Name.Local == "Complemento" {
				return
			}
		}
	}
}

func (h *CFDI40Handler) transformAddenda(decoder *xml.Decoder, data *models.CFDI40Data) {
	var addendaNames []string

	for {
		token, err := decoder.Token()
		if err != nil {
			return
		}

		switch t := token.(type) {
		case xml.StartElement:
			addendaNames = append(addendaNames, t.Name.Local)

		case xml.EndElement:
			if t.Name.Local == "Addenda" {
				if len(addendaNames) > 0 {
					data.CFDI40.Addendas = strings.Join(addendaNames, " ")
				}
				return
			}
		}
	}
}
