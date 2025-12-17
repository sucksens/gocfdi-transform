package models

type Nomina12Data struct {
	Version           string               `json:"version"`
	TipoNomina        string               `json:"tipo_nomina"`
	FechaPago         string               `json:"fecha_pago"`
	FechaInicialPago  string               `json:"fecha_inicial_pago"`
	FechaFinalPago    string               `json:"fecha_final_pago"`
	NumDiasPagados    string               `json:"num_dias_pagados"`
	TotalPercepciones string               `json:"total_percepciones,omitempty"`
	TotalDeducciones  string               `json:"total_deducciones,omitempty"`
	TotalOtrosPagos   string               `json:"total_otros_pagos,omitempty"`
	Emisor            Nomina12Emisor       `json:"emisor,omitempty"`
	Receptor          Nomina12Receptor     `json:"receptor"`
	Percepciones      Nomina12Percepciones `json:"percepciones,omitempty"`
	Deducciones       Deducciones          `json:"deducciones,omitempty"`
	OtrosPagos        OtrosPagos           `json:"otros_pagos,omitempty"`
	Incapacidades     Incapacidades        `json:"incapacidades,omitempty"`
}

type Nomina12Emisor struct {
	Curp             string      `json:"curp,omitempty"`
	RegistroPatronal string      `json:"registro_patronal,omitempty"`
	RfcPatronOrigen  string      `json:"rfc_patron_origen,omitempty"`
	EntidadSNCF      EntidadSNCF `json:"entidad_sncf,omitempty"`
}

type EntidadSNCF struct {
	OrigenRecurso      string `json:"origen_recurso"`
	MontoRecursoPropio string `json:"monto_recurso_propio,omitempty"`
}

type Nomina12Receptor struct {
	Curp                   string            `json:"curp"`
	NumSeguridadSocial     string            `json:"num_seguridad_social,omitempty"`
	FechaIniciaRelLaboral  string            `json:"fecha_inicia_rel_laboral,omitempty"`
	Antiguedad             string            `json:"antigüedad,omitempty"`
	TipoContrato           string            `json:"tipo_contrato,omitempty"`
	Sindicalizado          string            `json:"sindicalizado,omitempty"`
	TipoJornada            string            `json:"tipo_jornada,omitempty"`
	TipoRegimen            string            `json:"tipo_regimen"`
	NumEmpleado            string            `json:"num_empleado"`
	Departamento           string            `json:"departamento,omitempty"`
	Puesto                 string            `json:"puesto,omitempty"`
	RiesgoPuesto           string            `json:"riesgo_puesto,omitempty"`
	PeriodicidadPago       string            `json:"periodicidad_pago"`
	Banco                  string            `json:"banco,omitempty"`
	CuentaBancaria         string            `json:"cuenta_bancaria,omitempty"`
	SalarioBaseCotApor     string            `json:"salario_base_cot_apor,omitempty"`
	SalarioDiarioIntegrado string            `json:"salario_diario_integrado,omitempty"`
	ClaveEntFed            string            `json:"clave_ent_fed"`
	Subcontrataciones      []Subcontratacion `json:"subcontratacion,omitempty"`
}

type Subcontratacion struct {
	RfcLabora        string `json:"rfc_labora"`
	PorcentajeTiempo string `json:"porcentaje_tiempo"`
}

type Nomina12Percepciones struct {
	TotalSueldos                 string                  `json:"total_sueldos,omitempty"`
	TotalSeparacionIndemnizacion string                  `json:"total_separacion_indemnizacion,omitempty"`
	TotalJubilacionPensionRetiro string                  `json:"total_jubilacion_pension_retiro,omitempty"`
	TotalGravado                 string                  `json:"total_gravado"`
	TotalExento                  string                  `json:"total_exento"`
	Percepcion                   []Nomina12Percepcion    `json:"percepcion"`
	JubilacionPensionRetiro      JubilacionPensionRetiro `json:"jubilacion_pension_retiro,omitempty"`
	SeparacionIndemnizacion      SeparacionIndemnizacion `json:"separacion_indemnizacion,omitempty"`
}

type Nomina12Percepcion struct {
	TipoPercepcion   string           `json:"tipo_percepcion"`
	Clave            string           `json:"clave"`
	Concepto         string           `json:"concepto"`
	ImporteGravado   string           `json:"importe_gravado"`
	ImporteExento    string           `json:"importe_exento"`
	AccionesOTitulos AccionesOTitulos `json:"acciones_o_titulos,omitempty"`
	HorasExtra       []HorasExtra     `json:"horas_extra,omitempty"`
}

// AccionesOTitulos representa ingresos por acciones o títulos.
type AccionesOTitulos struct {
	ValorMercado      string `json:"valor_mercado"`
	PrecioAlOtorgarse string `json:"precio_al_otorgarse"`
}

// HorasExtra representa horas extra trabajadas.
type HorasExtra struct {
	Dias          string `json:"dias"`
	TipoHoras     string `json:"tipo_horas"`
	HorasExtra    string `json:"horas_extra"`
	ImportePagado string `json:"importe_pagado"`
}

// JubilacionPensionRetiro representa pagos por jubilación o pensión.
type JubilacionPensionRetiro struct {
	TotalUnaExhibicion  string `json:"total_una_exhibicion,omitempty"`
	TotalParcialidad    string `json:"total_parcialidad,omitempty"`
	MontoDiario         string `json:"monto_diario,omitempty"`
	IngresoAcumulable   string `json:"ingreso_acumulable"`
	IngresoNoAcumulable string `json:"ingreso_no_acumulable"`
}

// SeparacionIndemnizacion representa pagos por separación.
type SeparacionIndemnizacion struct {
	TotalPagado         string `json:"total_pagado"`
	NumAnosServicio     string `json:"num_años_servicio"`
	UltimoSueldoMensOrd string `json:"ultimo_sueldo_mens_ord"`
	IngresoAcumulable   string `json:"ingreso_acumulable"`
	IngresoNoAcumulable string `json:"ingreso_no_acumulable"`
}

// Deducciones representa las deducciones de nómina.
type Deducciones struct {
	TotalOtrasDeducciones   string      `json:"total_otras_deducciones,omitempty"`
	TotalImpuestosRetenidos string      `json:"total_impuestos_retenidos,omitempty"`
	Deduccion               []Deduccion `json:"deduccion"`
}

// Deduccion representa una deducción individual.
type Deduccion struct {
	TipoDeduccion string `json:"tipo_deduccion"`
	Clave         string `json:"clave"`
	Concepto      string `json:"concepto"`
	Importe       string `json:"importe"`
}

// OtrosPagos representa otros pagos de nómina.
type OtrosPagos struct {
	OtroPago []OtroPago `json:"otro_pago"`
}

// OtroPago representa un pago adicional.
type OtroPago struct {
	TipoOtroPago             string                   `json:"tipo_otro_pago"`
	Clave                    string                   `json:"clave"`
	Concepto                 string                   `json:"concepto"`
	Importe                  string                   `json:"importe"`
	SubsidioAlEmpleo         SubsidioAlEmpleo         `json:"subsidio_al_empleo,omitempty"`
	CompensacionSaldosAFavor CompensacionSaldosAFavor `json:"compensacion_saldos_a_favor,omitempty"`
}

// SubsidioAlEmpleo representa el subsidio al empleo.
type SubsidioAlEmpleo struct {
	SubsidioCausado string `json:"subsidio_causado"`
}

// CompensacionSaldosAFavor representa compensación de saldos a favor.
type CompensacionSaldosAFavor struct {
	SaldoAFavor     string `json:"saldo_a_favor"`
	Ano             string `json:"año"`
	RemanenteSalFav string `json:"remanente_sal_fav"`
}

// Incapacidades representa las incapacidades del empleado.
type Incapacidades struct {
	Incapacidad []Incapacidad `json:"incapacidad"`
}

// Incapacidad representa una incapacidad individual.
type Incapacidad struct {
	DiasIncapacidad  string `json:"dias_incapacidad"`
	TipoIncapacidad  string `json:"tipo_incapacidad"`
	ImporteMonetario string `json:"importe_monetario,omitempty"`
}
