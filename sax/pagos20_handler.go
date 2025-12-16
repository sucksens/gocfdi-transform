package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/helpers"
	"github.com/sucksens/gocfdi-transform/models"
)

// Pagos20Handler handles parsing of Pagos 2.0 complement.
type Pagos20Handler struct {
	config HandlerConfig
}

// NewPagos20Handler creates a new Pagos20Handler.
func NewPagos20Handler(cfg HandlerConfig) *Pagos20Handler {
	return &Pagos20Handler{config: cfg}
}

// ProcessPagosElement processes the Pagos element from an existing decoder stream.
func (h *Pagos20Handler) ProcessPagosElement(se xml.StartElement, decoder *xml.Decoder) (*models.Pagos20Data, error) {
	version := strings.TrimSpace(getAttrValue(se, "Version"))
	if version != "2.0" {
		return nil, errors.New("incorrect type of Pagos, this handler only supports Pagos version 2.0")
	}

	data := &models.Pagos20Data{
		Version: version,
		Pagos:   []models.Pago20{},
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "Totales":
				h.transformTotales(t, data)

			case "Pago":
				pago := h.transformPago(t, decoder)
				data.Pagos = append(data.Pagos, pago)
			}

		case xml.EndElement:
			if t.Name.Local == "Pagos" {
				return data, nil
			}
		}
	}
}

// TransformFromBytes parses a Pagos 2.0 XML byte slice.
func (h *Pagos20Handler) TransformFromBytes(xmlBytes []byte) (*models.Pagos20Data, error) {
	return h.TransformFromString(string(xmlBytes))
}

// TransformFromString parses a Pagos 2.0 XML string.
func (h *Pagos20Handler) TransformFromString(xmlStr string) (*models.Pagos20Data, error) {
	data := &models.Pagos20Data{
		Pagos: []models.Pago20{},
	}

	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

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
			case "Pagos":
				version := strings.TrimSpace(getAttrValue(se, "Version"))
				if version != "2.0" {
					return nil, errors.New("incorrect type of Pagos, this handler only supports Pagos version 2.0")
				}
				data.Version = version

			case "Totales":
				h.transformTotales(se, data)

			case "Pago":
				pago := h.transformPago(se, decoder)
				data.Pagos = append(data.Pagos, pago)
			}
		}
	}

	return data, nil
}

func (h *Pagos20Handler) transformTotales(se xml.StartElement, data *models.Pagos20Data) {
	data.Totales = models.Totales20{
		TotalRetencionesIVA:         helpers.GetOrDefault(getAttrValue(se, "TotalRetencionesIVA"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalRetencionesISR:         helpers.GetOrDefault(getAttrValue(se, "TotalRetencionesISR"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalRetencionesIEPS:        helpers.GetOrDefault(getAttrValue(se, "TotalRetencionesIEPS"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalTrasladosBaseIVA16:     helpers.GetOrDefault(getAttrValue(se, "TotalTrasladosBaseIVA16"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalTrasladosImpuestoIVA16: helpers.GetOrDefault(getAttrValue(se, "TotalTrasladosImpuestoIVA16"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalTrasladosBaseIVA8:      helpers.GetOrDefault(getAttrValue(se, "TotalTrasladosBaseIVA8"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalTrasladosImpuestoIVA8:  helpers.GetOrDefault(getAttrValue(se, "TotalTrasladosImpuestoIVA8"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalTrasladosBaseIVA0:      helpers.GetOrDefault(getAttrValue(se, "TotalTrasladosBaseIVA0"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalTrasladosImpuestoIVA0:  helpers.GetOrDefault(getAttrValue(se, "TotalTrasladosImpuestoIVA0"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalTrasladosBaseIVAExento: helpers.GetOrDefault(getAttrValue(se, "TotalTrasladosBaseIVAExento"), h.config.EmptyChar, h.config.SafeNumerics),
		MontoTotalPagos:             helpers.GetOrDefault(getAttrValue(se, "MontoTotalPagos"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Pagos20Handler) transformPago(se xml.StartElement, decoder *xml.Decoder) models.Pago20 {
	pago := models.Pago20{
		FechaPago:        getAttrValue(se, "FechaPago"),
		FormaDePagoP:     getAttrValue(se, "FormaDePagoP"),
		MonedaP:          getAttrValue(se, "MonedaP"),
		TipoCambioP:      helpers.GetOrDefaultOne(getAttrValue(se, "TipoCambioP"), h.config.EmptyChar, h.config.SafeNumerics),
		Monto:            getAttrValue(se, "Monto"),
		NumOperacion:     helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "NumOperacion", h.config.EmptyChar)),
		RfcEmisorCtaOrd:  getAttrValueOrDefault(se, "RfcEmisorCtaOrd", h.config.EmptyChar),
		NomBancoOrdExt:   helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "NomBancoOrdExt", h.config.EmptyChar)),
		CtaOrdenante:     getAttrValueOrDefault(se, "CtaOrdenante", h.config.EmptyChar),
		RfcEmisorCtaBen:  getAttrValueOrDefault(se, "RfcEmisorCtaBen", h.config.EmptyChar),
		CtaBeneficiario:  getAttrValueOrDefault(se, "CtaBeneficiario", h.config.EmptyChar),
		TipoCadPago:      getAttrValueOrDefault(se, "TipoCadPago", h.config.EmptyChar),
		CertPago:         helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "CertPago", h.config.EmptyChar)),
		CadPago:          helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "CadPago", h.config.EmptyChar)),
		SelloPago:        helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "SelloPago", h.config.EmptyChar)),
		DoctoRelacionado: []models.DoctoRelacionado20{},
		ImpuestosP:       []models.ImpuestosP{},
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return pago
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "DoctoRelacionado":
				docto := h.transformDoctoRelacionado(t, decoder)
				pago.DoctoRelacionado = append(pago.DoctoRelacionado, docto)

			case "ImpuestosP":
				impuestos := h.transformImpuestosP(t, decoder)
				pago.ImpuestosP = append(pago.ImpuestosP, impuestos)
			}

		case xml.EndElement:
			if t.Name.Local == "Pago" {
				return pago
			}
		}
	}
}

func (h *Pagos20Handler) transformDoctoRelacionado(se xml.StartElement, decoder *xml.Decoder) models.DoctoRelacionado20 {
	docto := models.DoctoRelacionado20{
		IdDocumento:      getAttrValue(se, "IdDocumento"),
		Serie:            helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Serie", h.config.EmptyChar)),
		Folio:            helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Folio", h.config.EmptyChar)),
		MonedaDR:         getAttrValue(se, "MonedaDR"),
		EquivalenciaDR:   helpers.GetOrDefault(getAttrValue(se, "EquivalenciaDR"), h.config.EmptyChar, h.config.SafeNumerics),
		NumParcialidad:   getAttrValueOrDefault(se, "NumParcialidad", h.config.EmptyChar),
		ImpSaldoAnt:      helpers.GetOrDefault(getAttrValue(se, "ImpSaldoAnt"), h.config.EmptyChar, h.config.SafeNumerics),
		ImpPagado:        helpers.GetOrDefault(getAttrValue(se, "ImpPagado"), h.config.EmptyChar, h.config.SafeNumerics),
		ImpSaldoInsoluto: helpers.GetOrDefault(getAttrValue(se, "ImpSaldoInsoluto"), h.config.EmptyChar, h.config.SafeNumerics),
		ObjetoImpDR:      getAttrValue(se, "ObjetoImpDR"),
		ImpuestosDR:      []models.ImpuestosDR{},
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return docto
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "ImpuestosDR" {
				impuestosDR := h.transformImpuestosDR(t, decoder)
				docto.ImpuestosDR = append(docto.ImpuestosDR, impuestosDR)
			}

		case xml.EndElement:
			if t.Name.Local == "DoctoRelacionado" {
				return docto
			}
		}
	}
}

func (h *Pagos20Handler) transformImpuestosDR(se xml.StartElement, decoder *xml.Decoder) models.ImpuestosDR {
	impuestos := models.ImpuestosDR{
		RetencionesDR: []models.ImpuestoDRItem{},
		TrasladosDR:   []models.ImpuestoDRItem{},
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return impuestos
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "RetencionDR":
				item := models.ImpuestoDRItem{
					BaseDR:       getAttrValue(t, "BaseDR"),
					ImpuestoDR:   getAttrValue(t, "ImpuestoDR"),
					TipoFactorDR: getAttrValue(t, "TipoFactorDR"),
					TasaOCuotaDR: helpers.GetOrDefault(getAttrValue(t, "TasaOCuotaDR"), h.config.EmptyChar, h.config.SafeNumerics),
					ImporteDR:    helpers.GetOrDefault(getAttrValue(t, "ImporteDR"), h.config.EmptyChar, h.config.SafeNumerics),
				}
				impuestos.RetencionesDR = append(impuestos.RetencionesDR, item)

			case "TrasladoDR":
				item := models.ImpuestoDRItem{
					BaseDR:       getAttrValue(t, "BaseDR"),
					ImpuestoDR:   getAttrValue(t, "ImpuestoDR"),
					TipoFactorDR: getAttrValue(t, "TipoFactorDR"),
					TasaOCuotaDR: helpers.GetOrDefault(getAttrValue(t, "TasaOCuotaDR"), h.config.EmptyChar, h.config.SafeNumerics),
					ImporteDR:    helpers.GetOrDefault(getAttrValue(t, "ImporteDR"), h.config.EmptyChar, h.config.SafeNumerics),
				}
				impuestos.TrasladosDR = append(impuestos.TrasladosDR, item)
			}

		case xml.EndElement:
			if t.Name.Local == "ImpuestosDR" {
				return impuestos
			}
		}
	}
}

func (h *Pagos20Handler) transformImpuestosP(se xml.StartElement, decoder *xml.Decoder) models.ImpuestosP {
	impuestos := models.ImpuestosP{
		RetencionesP: []models.RetencionP{},
		TrasladosP:   []models.TrasladoP{},
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return impuestos
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "RetencionP":
				retencion := models.RetencionP{
					ImpuestoP: getAttrValue(t, "ImpuestoP"),
					ImporteP:  getAttrValue(t, "ImporteP"),
				}
				impuestos.RetencionesP = append(impuestos.RetencionesP, retencion)

			case "TrasladoP":
				traslado := models.TrasladoP{
					BaseP:       getAttrValue(t, "BaseP"),
					ImpuestoP:   getAttrValue(t, "ImpuestoP"),
					TipoFactorP: getAttrValue(t, "TipoFactorP"),
					TasaOCuotaP: helpers.GetOrDefault(getAttrValue(t, "TasaOCuotaP"), h.config.EmptyChar, h.config.SafeNumerics),
					ImporteP:    helpers.GetOrDefault(getAttrValue(t, "ImporteP"), h.config.EmptyChar, h.config.SafeNumerics),
				}
				impuestos.TrasladosP = append(impuestos.TrasladosP, traslado)
			}

		case xml.EndElement:
			if t.Name.Local == "ImpuestosP" {
				return impuestos
			}
		}
	}
}
