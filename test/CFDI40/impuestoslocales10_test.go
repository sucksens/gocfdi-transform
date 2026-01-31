// Package cfdi40_test contiene tests unitarios para el handler de CFDI 4.0
// y todos sus complementos soportados.
package cfdi40_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sucksens/gocfdi-transform/sax"
)

// TestImpuestosLocales10 es la suite de tests para el complemento de Impuestos Locales 1.0.
//
// Estos tests verifican el correcto funcionamiento del parsing del complemento
// ImpuestosLocales cuando se encuentra dentro de un CFDI 4.0.
//
// Características evaluadas:
//   - Parseo básico de retenciones y traslados
//   - Uso del XML de ejemplo real del SAT
//   - Comportamiento cuando el flag está deshabilitado
//   - Manejo de SafeNumerics con valores vacíos
func TestImpuestosLocales10(t *testing.T) {
	// xmlStr contiene un CFDI 4.0 completo con el complemento de Impuestos Locales 1.0.
	// Este XML incluye:
	//   - Datos básicos del CFDI (emisor, receptor, conceptos)
	//   - Complemento de Impuestos Locales con 2 retenciones y 2 traslados
	xmlStr := `
	<cfdi:Comprobante xmlns:cfdi="http://www.sat.gob.mx/cfd/4" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:implocal="http://www.sat.gob.mx/implocal" xsi:schemaLocation="http://www.sat.gob.mx/cfd/4 cfdv40.xsd http://www.sat.gob.mx/implocal http://www.sat.gob.mx/sitio_internet/cfd/implocal/implocal.xsd" Version="4.0" Fecha="2021-12-08T23:59:59" Moneda="XXX" SubTotal="0" Total="0" TipoDeComprobante="P" FormaPago="02" Exportacion="03" LugarExpedicion="99999" NoCertificado="30001000000300023708" Sello="" Certificado="">
		<cfdi:CfdiRelacionados TipoRelacion="09">
			<cfdi:CfdiRelacionado UUID="F4F09AEF-57F2-4BE0-A828-87D1A80ED61C" />
		</cfdi:CfdiRelacionados>
		<cfdi:Emisor Rfc="AAA010101AAA" RegimenFiscal="622" Nombre="Esta es una demostración" />
		<cfdi:Receptor Rfc="BASJ600902KL9" UsoCFDI="P01" Nombre="Juanito Bananas De la Sierra" DomicilioFiscalReceptor="99999" RegimenFiscalReceptor="630" />
		<cfdi:Conceptos>
			<cfdi:Concepto ClaveProdServ="84111506" Cantidad="1" ClaveUnidad="ACT" Descripcion="Descripcion" ValorUnitario="0" Importe="0" ObjetoImp="01" />
		</cfdi:Conceptos>
		<cfdi:Complemento>
			<implocal:ImpuestosLocales version="1.0" TotaldeRetenciones="1357.80" TotaldeTraslados="1578.24">
				<implocal:RetencionesLocales ImpLocRetenido="Impuesto Estatal sobre Nómina" TasadeRetencion="2.00" Importe="678.90"/>
				<implocal:RetencionesLocales ImpLocRetenido="Impuesto Local de Traslación" TasadeRetencion="1.50" Importe="678.90"/>
				<implocal:TrasladosLocales ImpLocTrasladado="Impuesto Estatal" TasadeTraslado="3.00" Importe="789.12" />
				<implocal:TrasladosLocales ImpLocTrasladado="Impuesto Local por Servicios" TasadeTraslado="2.50" Importe="789.12" />
			</implocal:ImpuestosLocales>
		</cfdi:Complemento>
	</cfdi:Comprobante>
	`

	// Test Case 1: Parseo básico del complemento ImpuestosLocales 1.0
	//
	// Este test verifica que:
	//   - El complemento se detecta y parsea correctamente cuando el flag está habilitado
	//   - Los atributos del elemento principal (version, totales) se extraen correctamente
	//   - Las retenciones se parsean con todos sus atributos (nombre, tasa, importe)
	//   - Los traslados se parsean con todos sus atributos (nombre, tasa, importe)
	//   - Los arrays contienen la cantidad correcta de elementos
	t.Run("Parse ImpuestosLocales 1.0 complement", func(t *testing.T) {
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UseImpuestosLocales()
		data, err := handler.TransformFromString(xmlStr)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verificar que se haya parseado el complemento
		if len(data.ImpuestosLocales) == 0 {
			t.Fatal("Expected ImpuestosLocales complement, got none")
		}

		impLoc := data.ImpuestosLocales[0]

		// Validar atributos del elemento principal
		assert.Equal(t, "1.0", impLoc.Version)
		assert.Equal(t, "1357.80", impLoc.TotaldeRetenciones)
		assert.Equal(t, "1578.24", impLoc.TotaldeTraslados)

		// Validar retenciones locales (deben haber 2)
		assert.Len(t, impLoc.RetencionesLocales, 2)

		// Primera retención
		r1 := impLoc.RetencionesLocales[0]
		assert.Equal(t, "Impuesto Estatal sobre Nómina", r1.ImpLocRetenido)
		assert.Equal(t, "2.00", r1.TasadeRetencion)
		assert.Equal(t, "678.90", r1.Importe)

		// Segunda retención
		r2 := impLoc.RetencionesLocales[1]
		assert.Equal(t, "Impuesto Local de Traslación", r2.ImpLocRetenido)
		assert.Equal(t, "1.50", r2.TasadeRetencion)
		assert.Equal(t, "678.90", r2.Importe)

		// Validar traslados locales (deben haber 2)
		assert.Len(t, impLoc.TrasladosLocales, 2)

		// Primer traslado
		t1 := impLoc.TrasladosLocales[0]
		assert.Equal(t, "Impuesto Estatal", t1.ImpLocTrasladado)
		assert.Equal(t, "3.00", t1.TasadeTraslado)
		assert.Equal(t, "789.12", t1.Importe)

		// Segundo traslado
		t2 := impLoc.TrasladosLocales[1]
		assert.Equal(t, "Impuesto Local por Servicios", t2.ImpLocTrasladado)
		assert.Equal(t, "2.50", t2.TasadeTraslado)
		assert.Equal(t, "789.12", t2.Importe)
	})

	// Test Case 2: Uso del XML de ejemplo real del SAT
	//
	// Este test verifica que el XML de ejemplo proporcionado por el SAT
	// se puede parsear correctamente usando el archivo cfdi40_implocal.xml.
	//
	// Objetivos:
	//   - Validar integración con archivos XML reales
	//   - Verificar que el handler funciona con TransformFromFile
	//   - Comprobar que los valores del XML de ejemplo se extraen correctamente
	t.Run("Use existing CFDI40 file with ImpuestosLocales", func(t *testing.T) {
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UseImpuestosLocales()
		data, err := handler.TransformFromFile("../recursos/cfdi40_implocal.xml")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verificar que se haya parseado el complemento
		if len(data.ImpuestosLocales) == 0 {
			t.Fatal("Expected ImpuestosLocales complement, got none")
		}

		impLoc := data.ImpuestosLocales[0]

		// Validar valores esperados del XML de ejemplo
		assert.Equal(t, "1.0", impLoc.Version)
		assert.Equal(t, "5678.90", impLoc.TotaldeRetenciones)
		assert.Equal(t, "1234.56", impLoc.TotaldeTraslados)
		assert.Len(t, impLoc.RetencionesLocales, 2)
		assert.Len(t, impLoc.TrasladosLocales, 2)
	})

	// Test Case 3: Parseo sin flag de configuración habilitado
	//
	// Este test verifica que cuando el flag ParseImpuestosLocales está deshabilitado,
	// el complemento de Impuestos Locales NO se parsea.
	//
	// Objetivos:
	//   - Validar que el flag de configuración funciona correctamente
	//   - Verificar que el complemento se ignora cuando el flag es false
	//   - Asegurar que el CFDI se parsea correctamente sin el complemento
	t.Run("Parse without ImpuestosLocales flag should not parse", func(t *testing.T) {
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig())
		data, err := handler.TransformFromString(xmlStr)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verificar que NO se haya parseado el complemento
		assert.Len(t, data.ImpuestosLocales, 0, "Expected no ImpuestosLocales when flag is false")
	})

	// Test Case 4: Manejo de SafeNumerics con valores vacíos
	//
	// Este test verifica el comportamiento del handler cuando SafeNumerics está habilitado
	// y los atributos numéricos están vacíos.
	//
	// Características de SafeNumerics:
	//   - Los campos numéricos vacíos se reemplazan con "0.00"
	//   - Los campos de string vacíos se reemplazan con "" (emptyChar)
	//
	// Objetivos:
	//   - Validar que los valores vacíos usan "0.00" cuando SafeNumerics=true
	//   - Verificar que el parsing no falla con valores vacíos
	//   - Comprobar que la configuración SafeNumerics se aplica correctamente
	t.Run("SafeNumerics should use default values", func(t *testing.T) {
		config := sax.NewDefaultConfig()
		config.SafeNumerics = true
		handler := sax.NewCFDI40Handler(config).UseImpuestosLocales()

		// XML con valores numéricos vacíos para probar SafeNumerics
		xmlEmpty := `
		<cfdi:Comprobante xmlns:cfdi="http://www.sat.gob.mx/cfd/4" xmlns:implocal="http://www.sat.gob.mx/implocal" Version="4.0" Fecha="2021-12-08T23:59:59" Moneda="XXX" SubTotal="0" Total="0" TipoDeComprobante="P" LugarExpedicion="99999">
			<cfdi:Emisor Rfc="AAA010101AAA" RegimenFiscal="622" />
			<cfdi:Receptor Rfc="BASJ600902KL9" UsoCFDI="P01" DomicilioFiscalReceptor="99999" />
			<cfdi:Conceptos>
				<cfdi:Concepto ClaveProdServ="84111506" Cantidad="1" ClaveUnidad="ACT" Descripcion="Descripcion" ValorUnitario="0" Importe="0" ObjetoImp="01" />
			</cfdi:Conceptos>
			<cfdi:Complemento>
				<implocal:ImpuestosLocales version="1.0" TotaldeRetenciones="" TotaldeTraslados="">
					<implocal:RetencionesLocales ImpLocRetenido="" TasadeRetencion="" Importe=""/>
				</implocal:ImpuestosLocales>
			</cfdi:Complemento>
		</cfdi:Comprobante>
		`

		data, err := handler.TransformFromString(xmlEmpty)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		impLoc := data.ImpuestosLocales[0]

		// Validar que los campos numéricos vacíos usen "0.00" cuando SafeNumerics=true
		assert.Equal(t, "0.00", impLoc.TotaldeRetenciones)
		assert.Equal(t, "0.00", impLoc.TotaldeTraslados)

		// Validar las retenciones
		if len(impLoc.RetencionesLocales) > 0 {
			assert.Equal(t, "0.00", impLoc.RetencionesLocales[0].TasadeRetencion)
			assert.Equal(t, "0.00", impLoc.RetencionesLocales[0].Importe)
		}
	})
}
