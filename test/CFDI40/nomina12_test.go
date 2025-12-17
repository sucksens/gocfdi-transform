package cfdi40_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sucksens/gocfdi-transform/sax"
)

func TestNomina12(t *testing.T) {
	// Parse using the main CFDI40Handler with Nomina12 enabled
	handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UseNomina12()
	data, err := handler.TransformFromFile("../recursos/nomina12.xml")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Validate Nomina12 data was parsed
	if len(data.Nomina12) == 0 {
		t.Fatal("Expected Nomina12 complement, got none")
	}

	nomina := data.Nomina12[0]

	// Validaciones Generales
	assert.Equal(t, "1.2", nomina.Version)
	assert.Equal(t, "O", nomina.TipoNomina)
	assert.Equal(t, "2016-10-15", nomina.FechaPago)
	assert.Equal(t, "123.45", nomina.TotalPercepciones)

	// Emisor
	assert.Equal(t, "OAAJ840102HJCVRN00", nomina.Emisor.Curp)
	assert.Equal(t, "E23-12345-12-1", nomina.Emisor.RegistroPatronal)
	assert.Equal(t, "IP", nomina.Emisor.EntidadSNCF.OrigenRecurso)

	// Receptor
	assert.Equal(t, "123456789012345", nomina.Receptor.NumSeguridadSocial)
	assert.Equal(t, "P3Y2M23D", nomina.Receptor.Antiguedad)
	assert.Len(t, nomina.Receptor.Subcontrataciones, 2)
	assert.Equal(t, "AAA010101AAA", nomina.Receptor.Subcontrataciones[0].RfcLabora)

	// Percepciones
	assert.Equal(t, "123.45", nomina.Percepciones.TotalSueldos)
	assert.Len(t, nomina.Percepciones.Percepcion, 2)
	p1 := nomina.Percepciones.Percepcion[0]
	assert.Equal(t, "001", p1.TipoPercepcion)
	assert.Equal(t, "12345.67", p1.AccionesOTitulos.ValorMercado)
	assert.Len(t, p1.HorasExtra, 2)

	// Jubilacion
	assert.Equal(t, "223.45", nomina.Percepciones.JubilacionPensionRetiro.TotalUnaExhibicion)

	// Deducciones
	assert.Len(t, nomina.Deducciones.Deduccion, 2)
	assert.Equal(t, "001", nomina.Deducciones.Deduccion[0].TipoDeduccion)

	// OtrosPagos
	assert.Len(t, nomina.OtrosPagos.OtroPago, 2)
	op1 := nomina.OtrosPagos.OtroPago[0]
	assert.Equal(t, "003", op1.Clave)
	assert.Equal(t, "1234.56", op1.SubsidioAlEmpleo.SubsidioCausado)
	assert.Equal(t, "12345.67", op1.CompensacionSaldosAFavor.SaldoAFavor)

	// Incapacidades
	assert.Len(t, nomina.Incapacidades.Incapacidad, 2)
	assert.Equal(t, "1", nomina.Incapacidades.Incapacidad[0].DiasIncapacidad)
}
