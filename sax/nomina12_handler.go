package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/helpers"
	"github.com/sucksens/gocfdi-transform/models"
)

type Nomina12Handler struct {
	config HandlerConfig
}

func NewNomina12Handler(config HandlerConfig) *Nomina12Handler {
	return &Nomina12Handler{config: config}
}

func (h *Nomina12Handler) ProcessNomina12Element(se xml.StartElement, decoder *xml.Decoder) (*models.Nomina12Data, error) {
	version := strings.TrimSpace(getAttrValue(se, "Version"))
	if version != "1.2" {
		return nil, errors.New("incorrect type of Nomina 12, this handler only supports Nomina 12 version 1.2")
	}
	data := &models.Nomina12Data{
		Version:           version,
		TipoNomina:        getAttrValue(se, "TipoNomina"), // Enum
		FechaPago:         getAttrValue(se, "FechaPago"),
		FechaInicialPago:  getAttrValue(se, "FechaInicialPago"),
		FechaFinalPago:    getAttrValue(se, "FechaFinalPago"),
		NumDiasPagados:    helpers.GetOrDefault(getAttrValue(se, "NumDiasPagados"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalPercepciones: helpers.GetOrDefault(getAttrValue(se, "TotalPercepciones"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalDeducciones:  helpers.GetOrDefault(getAttrValue(se, "TotalDeducciones"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalOtrosPagos:   helpers.GetOrDefault(getAttrValue(se, "TotalOtrosPagos"), h.config.EmptyChar, h.config.SafeNumerics),
		Emisor:            models.Nomina12Emisor{},
		Receptor:          models.Nomina12Receptor{},
		Percepciones:      models.Nomina12Percepciones{},
		Deducciones:       models.Nomina12Deducciones{},
		OtrosPagos:        models.Nomina12OtrosPagos{},
		Incapacidades:     models.Nomina12Incapacidades{},
	}

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
			case "Emisor":
				data.Emisor = h.transformNomina12EmisorElement(t, decoder)
			case "Receptor":
				data.Receptor = h.transformNomina12ReceptorElement(t, decoder)
			case "Percepciones":
				data.Percepciones = h.transformNomina12PercepcionesElement(t, decoder)
			case "Deducciones":
				data.Deducciones = h.transformNomina12DeduccionesElement(t, decoder)
			case "OtrosPagos":
				data.OtrosPagos = h.transformNomina12OtrosPagosElement(t, decoder)
			case "Incapacidades":
				data.Incapacidades = h.transformNomina12IncapacidadesElement(t, decoder)
			}
		case xml.EndElement:
			if t.Name.Local == "Nomina12" || t.Name.Local == "Nomina" {
				return data, nil
			}
		}
	}
	return data, nil
}

func (h *Nomina12Handler) transformBytes(xmlBytes []byte) (*models.Nomina12Data, error) {
	return h.transformString(string(xmlBytes))
}

func (h *Nomina12Handler) transformString(xmlString string) (*models.Nomina12Data, error) {
	data := &models.Nomina12Data{
		Emisor:        models.Nomina12Emisor{},
		Receptor:      models.Nomina12Receptor{},
		Percepciones:  models.Nomina12Percepciones{},
		Deducciones:   models.Nomina12Deducciones{},
		OtrosPagos:    models.Nomina12OtrosPagos{},
		Incapacidades: models.Nomina12Incapacidades{},
	}

	decoder := xml.NewDecoder(strings.NewReader(xmlString))

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
			case "Emisor":
				data.Emisor = h.transformNomina12EmisorElement(t, decoder)
			case "Receptor":
				data.Receptor = h.transformNomina12ReceptorElement(t, decoder)
			case "Percepciones":
				data.Percepciones = h.transformNomina12PercepcionesElement(t, decoder)
			case "Deducciones":
				data.Deducciones = h.transformNomina12DeduccionesElement(t, decoder)
			case "OtrosPagos":
				data.OtrosPagos = h.transformNomina12OtrosPagosElement(t, decoder)
			case "Incapacidades":
				data.Incapacidades = h.transformNomina12IncapacidadesElement(t, decoder)
			}
		case xml.EndElement:
			if t.Name.Local == "Nomina12" || t.Name.Local == "Nomina" {
				return data, nil
			}
		}
	}
	return data, nil
}

func (h *Nomina12Handler) transformNomina12EmisorElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Emisor {
	emisor := models.Nomina12Emisor{
		Curp:             getAttrValue(se, "Curp"),
		RegistroPatronal: helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "RegistroPatronal", h.config.EmptyChar)),
		RfcPatronOrigen:  getAttrValueOrDefault(se, "RfcPatronOrigen", h.config.EmptyChar),
		EntidadSNCF:      models.EntidadSNCF{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return emisor
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "EntidadSNCF":
				emisor.EntidadSNCF = h.transformEntidadSNCFElement(t, decoder)
			}
		case xml.EndElement:
			if t.Name.Local == "Emisor" {
				return emisor
			}
		}
	}
	return emisor
}

func (h *Nomina12Handler) transformEntidadSNCFElement(se xml.StartElement, decoder *xml.Decoder) models.EntidadSNCF {
	return models.EntidadSNCF{
		OrigenRecurso:      getAttrValue(se, "OrigenRecurso"),
		MontoRecursoPropio: helpers.GetOrDefault(getAttrValue(se, "MontoRecursoPropio"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformNomina12ReceptorElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Receptor {
	receptor := models.Nomina12Receptor{
		Curp:                   getAttrValue(se, "Curp"),
		NumSeguridadSocial:     getAttrValueOrDefault(se, "NumSeguridadSocial", h.config.EmptyChar),
		FechaInicioRelLaboral:  getAttrValueOrDefault(se, "FechaInicioRelLaboral", h.config.EmptyChar),
		Antiguedad:             getAttrValueOrDefault(se, "Antigüedad", h.config.EmptyChar),
		TipoContrato:           getAttrValue(se, "TipoContrato"),
		Sindicalizado:          getAttrValueOrDefault(se, "Sindicalizado", h.config.EmptyChar),
		TipoJornada:            getAttrValueOrDefault(se, "TipoJornada", h.config.EmptyChar),
		TipoRegimen:            getAttrValue(se, "TipoRegimen"),
		NumEmpleado:            helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "NumEmpleado")),
		Departamento:           helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Departamento", h.config.EmptyChar)),
		Puesto:                 helpers.CompactString(h.config.EscDelimiters, getAttrValueOrDefault(se, "Puesto", h.config.EmptyChar)),
		RiesgoPuesto:           getAttrValueOrDefault(se, "RiesgoPuesto", h.config.EmptyChar),
		PeriodicidadPago:       getAttrValue(se, "PeriodicidadPago"),
		Banco:                  getAttrValueOrDefault(se, "Banco", h.config.EmptyChar),
		CuentaBancaria:         getAttrValueOrDefault(se, "CuentaBancaria", h.config.EmptyChar),
		SalarioBaseCotApor:     helpers.GetOrDefault(getAttrValue(se, "SalarioBaseCotApor"), h.config.EmptyChar, h.config.SafeNumerics),
		SalarioDiarioIntegrado: helpers.GetOrDefault(getAttrValue(se, "SalarioDiarioIntegrado"), h.config.EmptyChar, h.config.SafeNumerics),
		ClaveEntFed:            getAttrValue(se, "ClaveEntFed"),
		Subcontrataciones:      []models.Subcontratacion{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return receptor
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "SubContratacion":
				receptor.Subcontrataciones = append(receptor.Subcontrataciones, h.transformSubcontratacionElement(t, decoder))
			}
		case xml.EndElement:
			if t.Name.Local == "Receptor" {
				return receptor
			}
		}
	}
	return receptor
}

func (h *Nomina12Handler) transformSubcontratacionElement(se xml.StartElement, decoder *xml.Decoder) models.Subcontratacion {
	return models.Subcontratacion{
		RfcLabora:        getAttrValue(se, "RfcLabora"),
		PorcentajeTiempo: helpers.GetOrDefault(getAttrValue(se, "PorcentajeTiempo"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformNomina12PercepcionesElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Percepciones {
	Percepciones := models.Nomina12Percepciones{
		TotalSueldos:                 helpers.GetOrDefault(getAttrValue(se, "TotalSueldos"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalSeparacionIndemnizacion: helpers.GetOrDefault(getAttrValue(se, "TotalSeparacionIndemnizacion"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalJubilacionPensionRetiro: helpers.GetOrDefault(getAttrValue(se, "TotalJubilacionPensionRetiro"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalGravado:                 helpers.GetOrDefault(getAttrValue(se, "TotalGravado"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalExento:                  helpers.GetOrDefault(getAttrValue(se, "TotalExento"), h.config.EmptyChar, h.config.SafeNumerics),
		Percepcion:                   []models.Nomina12Percepcion{},
		JubilacionPensionRetiro:      models.JubilacionPensionRetiro{},
		SeparacionIndemnizacion:      models.SeparacionIndemnizacion{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Percepciones
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "Percepcion":
				Percepciones.Percepcion = append(Percepciones.Percepcion, h.transformNomina12PercepcionElement(t, decoder))
			case "JubilacionPensionRetiro":
				Percepciones.JubilacionPensionRetiro = h.transformJubilacionPensionRetiroElement(t, decoder)
			case "SeparacionIndemnizacion":
				Percepciones.SeparacionIndemnizacion = h.transformSeparacionIndemnizacionElement(t, decoder)
			}
		case xml.EndElement:
			if t.Name.Local == "Percepciones" {
				return Percepciones
			}
		}
	}
	return Percepciones
}

func (h *Nomina12Handler) transformNomina12PercepcionElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Percepcion {
	percepcion := models.Nomina12Percepcion{
		TipoPercepcion:   getAttrValue(se, "TipoPercepcion"),
		Clave:            getAttrValue(se, "Clave"),
		Concepto:         helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "Concepto")),
		ImporteGravado:   helpers.GetOrDefault(getAttrValue(se, "ImporteGravado"), h.config.EmptyChar, h.config.SafeNumerics),
		ImporteExento:    helpers.GetOrDefault(getAttrValue(se, "ImporteExento"), h.config.EmptyChar, h.config.SafeNumerics),
		AccionesOTitulos: models.AccionesOTitulos{},
		HorasExtra:       []models.HorasExtra{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return percepcion
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "AccionesOTitulos":
				percepcion.AccionesOTitulos = h.transformAccionesOTitulosElement(t, decoder)
			case "HorasExtra":
				percepcion.HorasExtra = append(percepcion.HorasExtra, h.transformHorasExtraElement(t, decoder))
			}
		case xml.EndElement:
			if t.Name.Local == "Percepcion" {
				return percepcion
			}
		}
	}
	return percepcion
}

func (h *Nomina12Handler) transformAccionesOTitulosElement(se xml.StartElement, decoder *xml.Decoder) models.AccionesOTitulos {
	return models.AccionesOTitulos{
		ValorMercado:      helpers.GetOrDefault(getAttrValue(se, "ValorMercado"), h.config.EmptyChar, h.config.SafeNumerics),
		PrecioAlOtorgarse: helpers.GetOrDefault(getAttrValue(se, "PrecioAlOtorgarse"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformHorasExtraElement(se xml.StartElement, decoder *xml.Decoder) models.HorasExtra {
	return models.HorasExtra{
		Dias:          helpers.GetOrDefault(getAttrValue(se, "Dias"), h.config.EmptyChar, h.config.SafeNumerics),
		TipoHoras:     getAttrValue(se, "TipoHoras"),
		HorasExtra:    helpers.GetOrDefault(getAttrValue(se, "HorasExtra"), h.config.EmptyChar, h.config.SafeNumerics),
		ImportePagado: helpers.GetOrDefault(getAttrValue(se, "ImportePagado"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformJubilacionPensionRetiroElement(se xml.StartElement, decoder *xml.Decoder) models.JubilacionPensionRetiro {
	return models.JubilacionPensionRetiro{
		TotalUnaExhibicion:  helpers.GetOrDefault(getAttrValue(se, "TotalUnaExhibicion"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalParcialidad:    helpers.GetOrDefault(getAttrValue(se, "TotalParcialidad"), h.config.EmptyChar, h.config.SafeNumerics),
		MontoDiario:         helpers.GetOrDefault(getAttrValue(se, "MontoDiario"), h.config.EmptyChar, h.config.SafeNumerics),
		IngresoAcumulable:   helpers.GetOrDefault(getAttrValue(se, "IngresoAcumulable"), h.config.EmptyChar, h.config.SafeNumerics),
		IngresoNoAcumulable: helpers.GetOrDefault(getAttrValue(se, "IngresoNoAcumulable"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformSeparacionIndemnizacionElement(se xml.StartElement, decoder *xml.Decoder) models.SeparacionIndemnizacion {
	return models.SeparacionIndemnizacion{
		TotalPagado:         helpers.GetOrDefault(getAttrValue(se, "TotalPagado"), h.config.EmptyChar, h.config.SafeNumerics),
		NumAnosServicio:     helpers.GetOrDefault(getAttrValue(se, "NumAñosServicio"), h.config.EmptyChar, h.config.SafeNumerics),
		UltimoSueldoMensOrd: helpers.GetOrDefault(getAttrValue(se, "UltimoSueldoMensOrd"), h.config.EmptyChar, h.config.SafeNumerics),
		IngresoAcumulable:   helpers.GetOrDefault(getAttrValue(se, "IngresoAcumulable"), h.config.EmptyChar, h.config.SafeNumerics),
		IngresoNoAcumulable: helpers.GetOrDefault(getAttrValue(se, "IngresoNoAcumulable"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformNomina12DeduccionesElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Deducciones {
	deducciones := models.Nomina12Deducciones{
		TotalOtrasDeducciones:   helpers.GetOrDefault(getAttrValue(se, "TotalOtrasDeducciones"), h.config.EmptyChar, h.config.SafeNumerics),
		TotalImpuestosRetenidos: helpers.GetOrDefault(getAttrValue(se, "TotalImpuestosRetenidos"), h.config.EmptyChar, h.config.SafeNumerics),
		Deduccion:               []models.Nomina12Deduccion{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return deducciones
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "Deduccion":
				deducciones.Deduccion = append(deducciones.Deduccion, h.transformNomina12DeduccionElement(t, decoder))
			}
		case xml.EndElement:
			if t.Name.Local == "Deducciones" {
				return deducciones
			}
		}
	}
	return deducciones
}

func (h *Nomina12Handler) transformNomina12DeduccionElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Deduccion {
	return models.Nomina12Deduccion{
		TipoDeduccion: getAttrValue(se, "TipoDeduccion"),
		Clave:         getAttrValue(se, "Clave"),
		Concepto:      helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "Concepto")),
		Importe:       helpers.GetOrDefault(getAttrValue(se, "Importe"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformNomina12OtrosPagosElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12OtrosPagos {
	otrosPagos := models.Nomina12OtrosPagos{
		OtroPago: []models.Nomina12OtroPago{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return otrosPagos
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "OtroPago":
				otrosPagos.OtroPago = append(otrosPagos.OtroPago, h.transformNomina12OtroPagoElement(t, decoder))
			}
		case xml.EndElement:
			if t.Name.Local == "OtrosPagos" {
				return otrosPagos
			}
		}
	}
	return otrosPagos
}

func (h *Nomina12Handler) transformNomina12OtroPagoElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12OtroPago {
	otroPago := models.Nomina12OtroPago{
		TipoOtroPago:             getAttrValue(se, "TipoOtroPago"),
		Clave:                    getAttrValue(se, "Clave"),
		Concepto:                 helpers.CompactString(h.config.EscDelimiters, getAttrValue(se, "Concepto")),
		Importe:                  helpers.GetOrDefault(getAttrValue(se, "Importe"), h.config.EmptyChar, h.config.SafeNumerics),
		SubsidioAlEmpleo:         models.SubsidioAlEmpleo{},
		CompensacionSaldosAFavor: models.CompensacionSaldosAFavor{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return otroPago
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "SubsidioAlEmpleo":
				otroPago.SubsidioAlEmpleo = h.transformSubsidioAlEmpleoElement(t, decoder)
			case "CompensacionSaldosAFavor":
				otroPago.CompensacionSaldosAFavor = h.transformCompensacionSaldosAFavorElement(t, decoder)
			}
		case xml.EndElement:
			if t.Name.Local == "OtroPago" {
				return otroPago
			}
		}
	}
	return otroPago
}

func (h *Nomina12Handler) transformSubsidioAlEmpleoElement(se xml.StartElement, decoder *xml.Decoder) models.SubsidioAlEmpleo {
	return models.SubsidioAlEmpleo{
		SubsidioCausado: helpers.GetOrDefault(getAttrValue(se, "SubsidioCausado"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformCompensacionSaldosAFavorElement(se xml.StartElement, decoder *xml.Decoder) models.CompensacionSaldosAFavor {
	return models.CompensacionSaldosAFavor{
		SaldoAFavor:     helpers.GetOrDefault(getAttrValue(se, "SaldoAFavor"), h.config.EmptyChar, h.config.SafeNumerics),
		Ano:             getAttrValue(se, "Año"),
		RemanenteSalFav: helpers.GetOrDefault(getAttrValue(se, "RemanenteSalFav"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}

func (h *Nomina12Handler) transformNomina12IncapacidadesElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Incapacidades {
	incapacidades := models.Nomina12Incapacidades{
		Incapacidad: []models.Nomina12Incapacidad{},
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return incapacidades
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "Incapacidad":
				incapacidades.Incapacidad = append(incapacidades.Incapacidad, h.transformNomina12IncapacidadElement(t, decoder))
			}
		case xml.EndElement:
			if t.Name.Local == "Incapacidades" {
				return incapacidades
			}
		}
	}
	return incapacidades
}

func (h *Nomina12Handler) transformNomina12IncapacidadElement(se xml.StartElement, decoder *xml.Decoder) models.Nomina12Incapacidad {
	return models.Nomina12Incapacidad{
		DiasIncapacidad:  helpers.GetOrDefault(getAttrValue(se, "DiasIncapacidad"), h.config.EmptyChar, h.config.SafeNumerics),
		TipoIncapacidad:  getAttrValue(se, "TipoIncapacidad"),
		ImporteMonetario: helpers.GetOrDefault(getAttrValue(se, "ImporteMonetario"), h.config.EmptyChar, h.config.SafeNumerics),
	}
}
