package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/models"
)

type Nomina12Handler struct {
	*BaseHandler
	builder *ModelBuilder
}

func NewNomina12Handler(config HandlerConfig) *Nomina12Handler {
	return &Nomina12Handler{
		BaseHandler: NewBaseHandler(config),
		builder:     NewModelBuilder(config),
	}
}

func (h *Nomina12Handler) ProcessNomina12Element(se xml.StartElement, decoder *xml.Decoder) (*models.Nomina12Data, error) {
	if err := h.ValidateVersion(se, "1.2"); err != nil {
		return nil, errors.New("incorrect type of Nomina, this handler only supports Nomina version 1.2")
	}
	data := &models.Nomina12Data{
		Version:           "1.2",
		TipoNomina:        h.builder.ExtractString(se, "TipoNomina"), // Enum
		FechaPago:         h.builder.ExtractString(se, "FechaPago"),
		FechaInicialPago:  h.builder.ExtractString(se, "FechaInicialPago"),
		FechaFinalPago:    h.builder.ExtractString(se, "FechaFinalPago"),
		NumDiasPagados:    h.builder.ExtractNumeric(se, "NumDiasPagados"),
		TotalPercepciones: h.builder.ExtractNumeric(se, "TotalPercepciones"),
		TotalDeducciones:  h.builder.ExtractNumeric(se, "TotalDeducciones"),
		TotalOtrosPagos:   h.builder.ExtractNumeric(se, "TotalOtrosPagos"),
		Emisor:            models.Nomina12Emisor{},
		Receptor:          models.Nomina12Receptor{},
		Percepciones:      models.Nomina12Percepciones{},
		Deducciones:       models.Nomina12Deducciones{},
		OtrosPagos:        models.Nomina12OtrosPagos{},
		Incapacidades:     models.Nomina12Incapacidades{},
	}

	err := ProcessChildElements(decoder, "Nomina", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "Emisor":
			data.Emisor = h.transformNomina12EmisorElement(childSE, childDecoder)
		case "Receptor":
			data.Receptor = h.transformNomina12ReceptorElement(childSE, childDecoder)
		case "Percepciones":
			data.Percepciones = h.transformNomina12PercepcionesElement(childSE, childDecoder)
		case "Deducciones":
			data.Deducciones = h.transformNomina12DeduccionesElement(childSE, childDecoder)
		case "OtrosPagos":
			data.OtrosPagos = h.transformNomina12OtrosPagosElement(childSE, childDecoder)
		case "Incapacidades":
			data.Incapacidades = h.transformNomina12IncapacidadesElement(childSE, childDecoder)
		}
		return nil
	})
	return data, err
}

func (h *Nomina12Handler) TransformFromString(xmlString string) (*models.Nomina12Data, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	se, err := FindElement(decoder, "Nomina")
	if err != nil {
		if err == io.EOF {
			return &models.Nomina12Data{
				Emisor:        models.Nomina12Emisor{},
				Receptor:      models.Nomina12Receptor{},
				Percepciones:  models.Nomina12Percepciones{},
				Deducciones:   models.Nomina12Deducciones{},
				OtrosPagos:    models.Nomina12OtrosPagos{},
				Incapacidades: models.Nomina12Incapacidades{},
			}, nil
		}
		return nil, err
	}

	return h.ProcessNomina12Element(*se, decoder)
}

func (h *Nomina12Handler) transformNomina12EmisorElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Emisor {
	emisor := models.Nomina12Emisor{
		Curp:             h.builder.ExtractString(se, "Curp"),
		RegistroPatronal: h.builder.ExtractCompact(se, "RegistroPatronal"),
		RfcPatronOrigen:  h.builder.ExtractStringOrDefault(se, "RfcPatronOrigen"),
		EntidadSNCF:      models.EntidadSNCF{},
	}
	ProcessChildElements(decoder, "Emisor", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "EntidadSNCF":
			emisor.EntidadSNCF = h.transformEntidadSNCFElement(childSE)
		}
		return nil
	})
	return emisor
}

func (h *Nomina12Handler) transformEntidadSNCFElement(se xml.StartElement) models.EntidadSNCF {
	return models.EntidadSNCF{
		OrigenRecurso:      h.builder.ExtractString(se, "OrigenRecurso"),
		MontoRecursoPropio: h.builder.ExtractNumeric(se, "MontoRecursoPropio"),
	}
}

func (h *Nomina12Handler) transformNomina12ReceptorElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Receptor {
	receptor := models.Nomina12Receptor{
		Curp:                   h.builder.ExtractString(se, "Curp"),
		NumSeguridadSocial:     h.builder.ExtractStringOrDefault(se, "NumSeguridadSocial"),
		FechaInicioRelLaboral:  h.builder.ExtractStringOrDefault(se, "FechaInicioRelLaboral"),
		Antiguedad:             h.builder.ExtractStringOrDefault(se, "Antigüedad"),
		TipoContrato:           h.builder.ExtractString(se, "TipoContrato"),
		Sindicalizado:          h.builder.ExtractStringOrDefault(se, "Sindicalizado"),
		TipoJornada:            h.builder.ExtractStringOrDefault(se, "TipoJornada"),
		TipoRegimen:            h.builder.ExtractString(se, "TipoRegimen"),
		NumEmpleado:            h.builder.ExtractCompact(se, "NumEmpleado"),
		Departamento:           h.builder.ExtractCompact(se, "Departamento"),
		Puesto:                 h.builder.ExtractCompact(se, "Puesto"),
		RiesgoPuesto:           h.builder.ExtractStringOrDefault(se, "RiesgoPuesto"),
		PeriodicidadPago:       h.builder.ExtractString(se, "PeriodicidadPago"),
		Banco:                  h.builder.ExtractStringOrDefault(se, "Banco"),
		CuentaBancaria:         h.builder.ExtractStringOrDefault(se, "CuentaBancaria"),
		SalarioBaseCotApor:     h.builder.ExtractNumeric(se, "SalarioBaseCotApor"),
		SalarioDiarioIntegrado: h.builder.ExtractNumeric(se, "SalarioDiarioIntegrado"),
		ClaveEntFed:            h.builder.ExtractString(se, "ClaveEntFed"),
		Subcontrataciones:      []models.Subcontratacion{},
	}

	ProcessChildElements(decoder, "Receptor", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "SubContratacion":
			receptor.Subcontrataciones = append(receptor.Subcontrataciones, h.transformSubcontratacionElement(childSE))
		}
		return nil
	})
	return receptor
}

func (h *Nomina12Handler) transformSubcontratacionElement(se xml.StartElement) models.Subcontratacion {
	return models.Subcontratacion{
		RfcLabora:        h.builder.ExtractString(se, "RfcLabora"),
		PorcentajeTiempo: h.builder.ExtractNumeric(se, "PorcentajeTiempo"),
	}
}

func (h *Nomina12Handler) transformNomina12PercepcionesElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Percepciones {
	Percepciones := models.Nomina12Percepciones{
		TotalSueldos:                 h.builder.ExtractNumeric(se, "TotalSueldos"),
		TotalSeparacionIndemnizacion: h.builder.ExtractNumeric(se, "TotalSeparacionIndemnizacion"),
		TotalJubilacionPensionRetiro: h.builder.ExtractNumeric(se, "TotalJubilacionPensionRetiro"),
		TotalGravado:                 h.builder.ExtractNumeric(se, "TotalGravado"),
		TotalExento:                  h.builder.ExtractNumeric(se, "TotalExento"),
		Percepcion:                   []models.Nomina12Percepcion{},
		JubilacionPensionRetiro:      models.JubilacionPensionRetiro{},
		SeparacionIndemnizacion:      models.SeparacionIndemnizacion{},
	}
	ProcessChildElements(decoder, "Percepciones", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "Percepcion":
			Percepciones.Percepcion = append(Percepciones.Percepcion, h.transformNomina12PercepcionElement(childSE, childDecoder))
		case "JubilacionPensionRetiro":
			Percepciones.JubilacionPensionRetiro = h.transformJubilacionPensionRetiroElement(childSE)
		case "SeparacionIndemnizacion":
			Percepciones.SeparacionIndemnizacion = h.transformSeparacionIndemnizacionElement(childSE)
		}
		return nil
	})
	return Percepciones
}

func (h *Nomina12Handler) transformNomina12PercepcionElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Percepcion {
	percepcion := models.Nomina12Percepcion{
		TipoPercepcion:   h.builder.ExtractString(se, "TipoPercepcion"),
		Clave:            h.builder.ExtractString(se, "Clave"),
		Concepto:         h.builder.ExtractCompact(se, "Concepto"),
		ImporteGravado:   h.builder.ExtractNumeric(se, "ImporteGravado"),
		ImporteExento:    h.builder.ExtractNumeric(se, "ImporteExento"),
		AccionesOTitulos: models.AccionesOTitulos{},
		HorasExtra:       []models.HorasExtra{},
	}
	ProcessChildElements(decoder, "Percepcion", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "AccionesOTitulos":
			percepcion.AccionesOTitulos = h.transformAccionesOTitulosElement(childSE)
		case "HorasExtra":
			percepcion.HorasExtra = append(percepcion.HorasExtra, h.transformHorasExtraElement(childSE))
		}
		return nil
	})
	return percepcion
}

func (h *Nomina12Handler) transformAccionesOTitulosElement(se xml.StartElement) models.AccionesOTitulos {
	return models.AccionesOTitulos{
		ValorMercado:      h.builder.ExtractNumeric(se, "ValorMercado"),
		PrecioAlOtorgarse: h.builder.ExtractNumeric(se, "PrecioAlOtorgarse"),
	}
}

func (h *Nomina12Handler) transformHorasExtraElement(se xml.StartElement) models.HorasExtra {
	return models.HorasExtra{
		Dias:          h.builder.ExtractNumeric(se, "Dias"),
		TipoHoras:     h.builder.ExtractString(se, "TipoHoras"),
		HorasExtra:    h.builder.ExtractNumeric(se, "HorasExtra"),
		ImportePagado: h.builder.ExtractNumeric(se, "ImportePagado"),
	}
}

func (h *Nomina12Handler) transformJubilacionPensionRetiroElement(se xml.StartElement) models.JubilacionPensionRetiro {
	return models.JubilacionPensionRetiro{
		TotalUnaExhibicion:  h.builder.ExtractNumeric(se, "TotalUnaExhibicion"),
		TotalParcialidad:    h.builder.ExtractNumeric(se, "TotalParcialidad"),
		MontoDiario:         h.builder.ExtractNumeric(se, "MontoDiario"),
		IngresoAcumulable:   h.builder.ExtractNumeric(se, "IngresoAcumulable"),
		IngresoNoAcumulable: h.builder.ExtractNumeric(se, "IngresoNoAcumulable"),
	}
}

func (h *Nomina12Handler) transformSeparacionIndemnizacionElement(se xml.StartElement) models.SeparacionIndemnizacion {
	return models.SeparacionIndemnizacion{
		TotalPagado:         h.builder.ExtractNumeric(se, "TotalPagado"),
		NumAnosServicio:     h.builder.ExtractNumeric(se, "NumAñosServicio"),
		UltimoSueldoMensOrd: h.builder.ExtractNumeric(se, "UltimoSueldoMensOrd"),
		IngresoAcumulable:   h.builder.ExtractNumeric(se, "IngresoAcumulable"),
		IngresoNoAcumulable: h.builder.ExtractNumeric(se, "IngresoNoAcumulable"),
	}
}

func (h *Nomina12Handler) transformNomina12DeduccionesElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Deducciones {
	deducciones := models.Nomina12Deducciones{
		TotalOtrasDeducciones:   h.builder.ExtractNumeric(se, "TotalOtrasDeducciones"),
		TotalImpuestosRetenidos: h.builder.ExtractNumeric(se, "TotalImpuestosRetenidos"),
		Deduccion:               []models.Nomina12Deduccion{},
	}

	ProcessChildElements(decoder, "Deducciones", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "Deduccion":
			deducciones.Deduccion = append(deducciones.Deduccion, h.transformNomina12DeduccionElement(childSE, childDecoder))
		}
		return nil
	})
	return deducciones
}

func (h *Nomina12Handler) transformNomina12DeduccionElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Deduccion {
	return models.Nomina12Deduccion{
		TipoDeduccion: h.builder.ExtractString(se, "TipoDeduccion"),
		Clave:         h.builder.ExtractString(se, "Clave"),
		Concepto:      h.builder.ExtractCompact(se, "Concepto"),
		Importe:       h.builder.ExtractNumeric(se, "Importe"),
	}
}

func (h *Nomina12Handler) transformNomina12OtrosPagosElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12OtrosPagos {
	otrosPagos := models.Nomina12OtrosPagos{
		OtroPago: []models.Nomina12OtroPago{},
	}
	ProcessChildElements(decoder, "OtrosPagos", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "OtroPago":
			otrosPagos.OtroPago = append(otrosPagos.OtroPago, h.transformNomina12OtroPagoElement(childSE, childDecoder))
		}
		return nil
	})
	return otrosPagos
}

func (h *Nomina12Handler) transformNomina12OtroPagoElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12OtroPago {
	otroPago := models.Nomina12OtroPago{
		TipoOtroPago:             h.builder.ExtractString(se, "TipoOtroPago"),
		Clave:                    h.builder.ExtractString(se, "Clave"),
		Concepto:                 h.builder.ExtractCompact(se, "Concepto"),
		Importe:                  h.builder.ExtractNumeric(se, "Importe"),
		SubsidioAlEmpleo:         models.SubsidioAlEmpleo{},
		CompensacionSaldosAFavor: models.CompensacionSaldosAFavor{},
	}
	ProcessChildElements(decoder, "OtroPago", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "SubsidioAlEmpleo":
			otroPago.SubsidioAlEmpleo = h.transformSubsidioAlEmpleoElement(childSE)
		case "CompensacionSaldosAFavor":
			otroPago.CompensacionSaldosAFavor = h.transformCompensacionSaldosAFavorElement(childSE)
		}
		return nil
	})
	return otroPago
}

func (h *Nomina12Handler) transformSubsidioAlEmpleoElement(se xml.StartElement) models.SubsidioAlEmpleo {
	return models.SubsidioAlEmpleo{
		SubsidioCausado: h.builder.ExtractNumeric(se, "SubsidioCausado"),
	}
}

func (h *Nomina12Handler) transformCompensacionSaldosAFavorElement(se xml.StartElement) models.CompensacionSaldosAFavor {
	return models.CompensacionSaldosAFavor{
		SaldoAFavor:     h.builder.ExtractNumeric(se, "SaldoAFavor"),
		Ano:             h.builder.ExtractString(se, "Año"),
		RemanenteSalFav: h.builder.ExtractNumeric(se, "RemanenteSalFav"),
	}
}

func (h *Nomina12Handler) transformNomina12IncapacidadesElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Incapacidades {
	incapacidades := models.Nomina12Incapacidades{
		Incapacidad: []models.Nomina12Incapacidad{},
	}
	ProcessChildElements(decoder, "Incapacidades", func(childSE xml.StartElement, childDecoder *xml.Decoder) error {
		switch childSE.Name.Local {
		case "Incapacidad":
			incapacidades.Incapacidad = append(incapacidades.Incapacidad, h.transformNomina12IncapacidadElement(childSE))
		}
		return nil
	})
	return incapacidades
}

func (h *Nomina12Handler) transformNomina12IncapacidadElement(se xml.StartElement) models.Nomina12Incapacidad {
	return models.Nomina12Incapacidad{
		DiasIncapacidad:  h.builder.ExtractNumeric(se, "DiasIncapacidad"),
		TipoIncapacidad:  h.builder.ExtractString(se, "TipoIncapacidad"),
		ImporteMonetario: h.builder.ExtractNumeric(se, "ImporteMonetario"),
	}
}
