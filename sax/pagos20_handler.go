package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/models"
)

// Pagos20Handler handles parsing of Pagos 2.0 complement.
type Pagos20Handler struct {
	*BaseHandler
	builder *ModelBuilder
}

// NewPagos20Handler creates a new Pagos20Handler.
func NewPagos20Handler(config HandlerConfig) *Pagos20Handler {
	return &Pagos20Handler{
		BaseHandler: NewBaseHandler(config),
		builder:     NewModelBuilder(config),
	}
}

// ProcessPagosElement processes the Pagos element from an existing decoder stream.
func (h *Pagos20Handler) ProcessPagosElement(se xml.StartElement, decoder *xml.Decoder) (*models.Pagos20Data, error) {
	if err := h.ValidateVersion(se, "2.0"); err != nil {
		return nil, errors.New("incorrect type of Pagos, this handler only supports Pagos version 2.0")
	}
	// Inicializar la estructura de datos
	data := &models.Pagos20Data{
		Version: "2.0",
		Pagos:   []models.Pago20{},
	}
	// Iterar sobre los elementos hijos de pagos
	err := ProcessChildElements(decoder, "Pagos", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "Totales":
			h.transformTotales(childSE, data)

		case "Pago":
			pago := h.transformPago(childSE, childDecoder)
			data.Pagos = append(data.Pagos, pago)
		}
		return nil
	})
	return data, err
}

// TransformFromString parses a Pagos 2.0 XML string.
func (h *Pagos20Handler) TransformFromString(xmlStr string) (*models.Pagos20Data, error) {
	// Crear el decoder XML para  parsear la cadena
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))
	// Buscar el elemento Pagos
	se, err := FindElement(decoder, "Pagos")
	if err != nil {
		if err == io.EOF {
			return &models.Pagos20Data{
				Pagos: []models.Pago20{},
			}, nil
		}
		return nil, err
	}
	return h.ProcessPagosElement(*se, decoder)
}

func (h *Pagos20Handler) transformTotales(se xml.StartElement, data *models.Pagos20Data) {
	data.Totales = models.Totales20{
		TotalRetencionesIVA:         h.builder.ExtractNumeric(se, "TotalRetencionesIVA"),
		TotalRetencionesISR:         h.builder.ExtractNumeric(se, "TotalRetencionesISR"),
		TotalRetencionesIEPS:        h.builder.ExtractNumeric(se, "TotalRetencionesIEPS"),
		TotalTrasladosBaseIVA16:     h.builder.ExtractNumeric(se, "TotalTrasladosBaseIVA16"),
		TotalTrasladosImpuestoIVA16: h.builder.ExtractNumeric(se, "TotalTrasladosImpuestoIVA16"),
		TotalTrasladosBaseIVA8:      h.builder.ExtractNumeric(se, "TotalTrasladosBaseIVA8"),
		TotalTrasladosImpuestoIVA8:  h.builder.ExtractNumeric(se, "TotalTrasladosImpuestoIVA8"),
		TotalTrasladosBaseIVA0:      h.builder.ExtractNumeric(se, "TotalTrasladosBaseIVA0"),
		TotalTrasladosImpuestoIVA0:  h.builder.ExtractNumeric(se, "TotalTrasladosImpuestoIVA0"),
		TotalTrasladosBaseIVAExento: h.builder.ExtractNumeric(se, "TotalTrasladosBaseIVAExento"),
		MontoTotalPagos:             h.builder.ExtractNumeric(se, "MontoTotalPagos"),
	}
}

func (h *Pagos20Handler) transformPago(se xml.StartElement, decoder *xml.Decoder) models.Pago20 {
	pago := models.Pago20{
		FechaPago:        h.builder.ExtractString(se, "FechaPago"),
		FormaDePagoP:     h.builder.ExtractString(se, "FormaDePagoP"),
		MonedaP:          h.builder.ExtractString(se, "MonedaP"),
		TipoCambioP:      h.builder.ExtractNumericOne(se, "TipoCambioP"),
		Monto:            h.builder.ExtractString(se, "Monto"),
		NumOperacion:     h.builder.ExtractCompact(se, "NumOperacion"),
		RfcEmisorCtaOrd:  h.builder.ExtractStringOrDefault(se, "RfcEmisorCtaOrd"),
		NomBancoOrdExt:   h.builder.ExtractCompact(se, "NomBancoOrdExt"),
		CtaOrdenante:     h.builder.ExtractStringOrDefault(se, "CtaOrdenante"),
		RfcEmisorCtaBen:  h.builder.ExtractStringOrDefault(se, "RfcEmisorCtaBen"),
		CtaBeneficiario:  h.builder.ExtractStringOrDefault(se, "CtaBeneficiario"),
		TipoCadPago:      h.builder.ExtractStringOrDefault(se, "TipoCadPago"),
		CertPago:         h.builder.ExtractCompact(se, "CertPago"),
		CadPago:          h.builder.ExtractCompact(se, "CadPago"),
		SelloPago:        h.builder.ExtractCompact(se, "SelloPago"),
		DoctoRelacionado: []models.DoctoRelacionado20{},
		ImpuestosP:       []models.ImpuestosP{},
	}

	ProcessChildElements(decoder, "Pago", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "DoctoRelacionado":
			docto := h.transformDoctoRelacionado(childSE, childDecoder)
			pago.DoctoRelacionado = append(pago.DoctoRelacionado, docto)

		case "ImpuestosP":
			impuestos := h.transformImpuestosP(childSE, childDecoder)
			pago.ImpuestosP = append(pago.ImpuestosP, impuestos)
		}
		return nil
	})
	return pago
}

func (h *Pagos20Handler) transformDoctoRelacionado(se xml.StartElement, decoder *xml.Decoder) models.DoctoRelacionado20 {
	docto := models.DoctoRelacionado20{
		IdDocumento:      h.builder.ExtractString(se, "IdDocumento"),
		Serie:            h.builder.ExtractCompact(se, "Serie"),
		Folio:            h.builder.ExtractCompact(se, "Folio"),
		MonedaDR:         h.builder.ExtractString(se, "MonedaDR"),
		EquivalenciaDR:   h.builder.ExtractNumeric(se, "EquivalenciaDR"),
		NumParcialidad:   h.builder.ExtractStringOrDefault(se, "NumParcialidad"),
		ImpSaldoAnt:      h.builder.ExtractNumeric(se, "ImpSaldoAnt"),
		ImpPagado:        h.builder.ExtractNumeric(se, "ImpPagado"),
		ImpSaldoInsoluto: h.builder.ExtractNumeric(se, "ImpSaldoInsoluto"),
		ObjetoImpDR:      h.builder.ExtractString(se, "ObjetoImpDR"),
		ImpuestosDR:      []models.ImpuestosDR{},
	}

	ProcessChildElements(decoder, "DoctoRelacionado", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "ImpuestosDR":
			impuestosDR := h.transformImpuestosDR(childSE, childDecoder)
			docto.ImpuestosDR = append(docto.ImpuestosDR, impuestosDR)
		}
		return nil
	})
	return docto
}

func (h *Pagos20Handler) transformImpuestosDR(se xml.StartElement, decoder *xml.Decoder) models.ImpuestosDR {
	impuestos := models.ImpuestosDR{
		RetencionesDR: []models.ImpuestoDRItem{},
		TrasladosDR:   []models.ImpuestoDRItem{},
	}

	ProcessChildElements(decoder, "ImpuestosDR", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "RetencionDR":
			item := models.ImpuestoDRItem{
				BaseDR:       h.builder.ExtractString(childSE, "BaseDR"),
				ImpuestoDR:   h.builder.ExtractString(childSE, "ImpuestoDR"),
				TipoFactorDR: h.builder.ExtractString(childSE, "TipoFactorDR"),
				TasaOCuotaDR: h.builder.ExtractNumeric(childSE, "TasaOCuotaDR"),
				ImporteDR:    h.builder.ExtractNumeric(childSE, "ImporteDR"),
			}
			impuestos.RetencionesDR = append(impuestos.RetencionesDR, item)

		case "TrasladoDR":
			item := models.ImpuestoDRItem{
				BaseDR:       h.builder.ExtractString(childSE, "BaseDR"),
				ImpuestoDR:   h.builder.ExtractString(childSE, "ImpuestoDR"),
				TipoFactorDR: h.builder.ExtractString(childSE, "TipoFactorDR"),
				TasaOCuotaDR: h.builder.ExtractNumeric(childSE, "TasaOCuotaDR"),
				ImporteDR:    h.builder.ExtractNumeric(childSE, "ImporteDR"),
			}
			impuestos.TrasladosDR = append(impuestos.TrasladosDR, item)
		}
		return nil
	})
	return impuestos
}

func (h *Pagos20Handler) transformImpuestosP(se xml.StartElement, decoder *xml.Decoder) models.ImpuestosP {
	impuestos := models.ImpuestosP{
		RetencionesP: []models.RetencionP{},
		TrasladosP:   []models.TrasladoP{},
	}
	ProcessChildElements(decoder, "ImpuestosP", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "RetencionP":
			retencion := models.RetencionP{
				ImpuestoP: h.builder.ExtractString(childSE, "ImpuestoP"),
				ImporteP:  h.builder.ExtractString(childSE, "ImporteP"),
			}
			impuestos.RetencionesP = append(impuestos.RetencionesP, retencion)

		case "TrasladoP":
			traslado := models.TrasladoP{
				BaseP:       h.builder.ExtractString(childSE, "BaseP"),
				ImpuestoP:   h.builder.ExtractString(childSE, "ImpuestoP"),
				TipoFactorP: h.builder.ExtractString(childSE, "TipoFactorP"),
				TasaOCuotaP: h.builder.ExtractNumeric(childSE, "TasaOCuotaP"),
				ImporteP:    h.builder.ExtractNumeric(childSE, "ImporteP"),
			}
			impuestos.TrasladosP = append(impuestos.TrasladosP, traslado)
		}
		return nil
	})
	return impuestos
}
