package cfdi40_test

import (
	"testing"

	"github.com/sucksens/gocfdi-transform/sax"
)

func TestCFDI40Handler(t *testing.T) {
	xmlStr := `
	<cfdi:Comprobante xmlns:cfdi="http://www.sat.gob.mx/cfd/4" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sat.gob.mx/cfd/4 http://www.sat.gob.mx/sitio_internet/cfd/4/cfdv40.xsd" Version="4.0" Fecha="2025-01-15T10:30:00" Sello="SELLO_DE_EJEMPLO_1234567890_ABCDEFGHIJKLMNOPQRSTUVWXYZ_==" FormaPago="03" NoCertificado="30001000000300023788" Certificado="CERTIFICADO_DE_EJEMPLO_ABCDEF1234567890_==" SubTotal="1000.00" Moneda="MXN" Total="1160.00" TipoDeComprobante="I" Exportacion="01" MetodoPago="PUE" LugarExpedicion="01000" Serie="AAA" Folio="12345">
		<cfdi:CfdiRelacionados TipoRelacion="04">
			<cfdi:CfdiRelacionado UUID="5FB2822E-396D-4725-8521-500FAB000222"/>
		</cfdi:CfdiRelacionados>
		<cfdi:Emisor Rfc="AAA010101AAA" Nombre="EMISOR DE PRUEBA SA DE CV" RegimenFiscal="601"/>
		<cfdi:Receptor Rfc="XAXX010101000" Nombre="PUBLICO EN GENERAL" DomicilioFiscalReceptor="01000" RegimenFiscalReceptor="616" UsoCFDI="G03"/>
		<cfdi:Conceptos>
			<cfdi:Concepto ClaveProdServ="84111506" Cantidad="1" ClaveUnidad="ACT" Descripcion="SERVICIO DE EJEMPLO" ValorUnitario="1000" Importe="1000" ObjetoImp="02">
				<cfdi:Impuestos>
					<cfdi:Traslados>
						<cfdi:Traslado Base="1000.00" Impuesto="002" TipoFactor="Tasa" TasaOCuota="0.160000" Importe="160.00"/>
					</cfdi:Traslados>
				</cfdi:Impuestos>
			</cfdi:Concepto>
		</cfdi:Conceptos>
		<cfdi:Impuestos TotalImpuestosTrasladados="160.00">
			<cfdi:Traslados>
				<cfdi:Traslado Base="1000.00" Impuesto="002" TipoFactor="Tasa" TasaOCuota="0.160000" Importe="160.00"/>
			</cfdi:Traslados>
		</cfdi:Impuestos>
		<cfdi:Complemento>
			<tfd:TimbreFiscalDigital xmlns:tfd="http://www.sat.gob.mx/TimbreFiscalDigital" xsi:schemaLocation="http://www.sat.gob.mx/TimbreFiscalDigital http://www.sat.gob.mx/sitio_internet/cfd/TimbreFiscalDigital/TimbreFiscalDigitalv11.xsd" Version="1.1" UUID="a3c6a0d7-8f4b-4e2a-9b5c-1d8e9f7a6b2c" FechaTimbrado="2025-01-15T10:30:01" RfcProvCertif="AAA010101AAA" SelloCFD="SELLO_CFD_DE_EJEMPLO_1234567890_==" NoCertificadoSAT="30001000000300023789" SelloSAT="SELLO_SAT_DE_EJEMPLO_1234567890_=="/>
		</cfdi:Complemento>
	</cfdi:Comprobante>
	`
	t.Run("Parse CFDI40 with default config", func(t *testing.T) {
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig())
		data, err := handler.TransformFromString(xmlStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Validacion de que sea version 4.0
		if data.CFDI40.Version != "4.0" {
			t.Errorf("Expected CFDI40 version 4.0, got %s", data.CFDI40.Version)
		}
		//validacion de la serie
		if data.CFDI40.Serie != "AAA" {
			t.Errorf("Expected CFDI40 serie AAA, got %s", data.CFDI40.Serie)
		}
		//validacion del folio
		if data.CFDI40.Folio != "12345" {
			t.Errorf("Expected CFDI40 folio 12345, got %s", data.CFDI40.Folio)
		}

		//validacion del emisor
		if data.CFDI40.Emisor.RFC != "AAA010101AAA" {
			t.Errorf("Expected Emisor40 rfc AAA010101AAA, got %s", data.CFDI40.Emisor.RFC)
		}
		//validacion del receptor
		if data.CFDI40.Receptor.RFC != "XAXX010101000" {
			t.Errorf("Expected Receptor40 rfc XAXX010101000, got %s", data.CFDI40.Receptor.RFC)
		}

		// Validacion de que tenga un TFD11
		if len(data.TFD11) == 0 {
			t.Fatal("Expected TFD11 data, got none")
		}

		// Validacion de atributos del Comprobante
		if data.CFDI40.Fecha != "2025-01-15T10:30:00" {
			t.Errorf("Expected CFDI40 Fecha 2025-01-15T10:30:00, got %s", data.CFDI40.Fecha)
		}
		if data.CFDI40.Sello != "SELLO_DE_EJEMPLO_1234567890_ABCDEFGHIJKLMNOPQRSTUVWXYZ_==" {
			t.Errorf("Expected CFDI40 Sello mismatch, got %s", data.CFDI40.Sello)
		}
		if data.CFDI40.FormaPago != "03" {
			t.Errorf("Expected CFDI40 FormaPago 03, got %s", data.CFDI40.FormaPago)
		}
		if data.CFDI40.NoCertificado != "30001000000300023788" {
			t.Errorf("Expected CFDI40 NoCertificado 30001000000300023788, got %s", data.CFDI40.NoCertificado)
		}
		if data.CFDI40.Certificado != "CERTIFICADO_DE_EJEMPLO_ABCDEF1234567890_==" {
			t.Errorf("Expected CFDI40 Certificado mismatch, got %s", data.CFDI40.Certificado)
		}
		if data.CFDI40.SubTotal != "1000.00" {
			t.Errorf("Expected CFDI40 SubTotal 1000.00, got %s", data.CFDI40.SubTotal)
		}
		if data.CFDI40.Moneda != "MXN" {
			t.Errorf("Expected CFDI40 Moneda MXN, got %s", data.CFDI40.Moneda)
		}
		if data.CFDI40.Total != "1160.00" {
			t.Errorf("Expected CFDI40 Total 1160.00, got %s", data.CFDI40.Total)
		}
		if data.CFDI40.TipoComprobante != "I" {
			t.Errorf("Expected CFDI40 TipoComprobante I, got %s", data.CFDI40.TipoComprobante)
		}
		if data.CFDI40.Exportacion != "01" {
			t.Errorf("Expected CFDI40 Exportacion 01, got %s", data.CFDI40.Exportacion)
		}
		if data.CFDI40.MetodoPago != "PUE" {
			t.Errorf("Expected CFDI40 MetodoPago PUE, got %s", data.CFDI40.MetodoPago)
		}
		if data.CFDI40.LugarExpedicion != "01000" {
			t.Errorf("Expected CFDI40 LugarExpedicion 01000, got %s", data.CFDI40.LugarExpedicion)
		}

		// Validacion de atributos del Emisor
		if data.CFDI40.Emisor.Nombre != "EMISOR DE PRUEBA SA DE CV" {
			t.Errorf("Expected Emisor Nombre EMISOR DE PRUEBA SA DE CV, got %s", data.CFDI40.Emisor.Nombre)
		}
		if data.CFDI40.Emisor.RegimenFiscal != "601" {
			t.Errorf("Expected Emisor RegimenFiscal 601, got %s", data.CFDI40.Emisor.RegimenFiscal)
		}

		// Validacion de atributos del Receptor
		if data.CFDI40.Receptor.Nombre != "PUBLICO EN GENERAL" {
			t.Errorf("Expected Receptor Nombre PUBLICO EN GENERAL, got %s", data.CFDI40.Receptor.Nombre)
		}
		if data.CFDI40.Receptor.DomicilioFiscalReceptor != "01000" {
			t.Errorf("Expected Receptor DomicilioFiscalReceptor 01000, got %s", data.CFDI40.Receptor.DomicilioFiscalReceptor)
		}
		if data.CFDI40.Receptor.RegimenFiscalReceptor != "616" {
			t.Errorf("Expected Receptor RegimenFiscalReceptor 616, got %s", data.CFDI40.Receptor.RegimenFiscalReceptor)
		}
		if data.CFDI40.Receptor.UsoCFDI != "G03" {
			t.Errorf("Expected Receptor UsoCFDI G03, got %s", data.CFDI40.Receptor.UsoCFDI)
		}

		// Validacion de Impuestos Globales
		if data.CFDI40.Impuestos.TotalImpuestosTrasladados != "160.00" {
			t.Errorf("Expected Impuestos TotalImpuestosTrasladados 160.00, got %s", data.CFDI40.Impuestos.TotalImpuestosTrasladados)
		}
		if len(data.CFDI40.Impuestos.Traslados) != 1 {
			t.Errorf("Expected 1 global Traslado, got %d", len(data.CFDI40.Impuestos.Traslados))
		} else {
			tr := data.CFDI40.Impuestos.Traslados[0]
			if tr.Base != "1000.00" {
				t.Errorf("Expected Traslado Base 1000.00, got %s", tr.Base)
			}
			if tr.Impuesto != "002" {
				t.Errorf("Expected Traslado Impuesto 002, got %s", tr.Impuesto)
			}
			if tr.TipoFactor != "Tasa" {
				t.Errorf("Expected Traslado TipoFactor Tasa, got %s", tr.TipoFactor)
			}
			if tr.TasaOCuota != "0.160000" {
				t.Errorf("Expected Traslado TasaOCuota 0.160000, got %s", tr.TasaOCuota)
			}
			if tr.Importe != "160.00" {
				t.Errorf("Expected Traslado Importe 160.00, got %s", tr.Importe)
			}
		}

		// Validacion de TFD11
		tfd := data.TFD11[0]
		if tfd.Version != "1.1" {
			t.Errorf("Expected TFD Version 1.1, got %s", tfd.Version)
		}
		if tfd.UUID != "A3C6A0D7-8F4B-4E2A-9B5C-1D8E9F7A6B2C" {
			t.Errorf("Expected TFD UUID A3C6A0D7-8F4B-4E2A-9B5C-1D8E9F7A6B2C, got %s", tfd.UUID)
		}
		if tfd.FechaTimbrado != "2025-01-15T10:30:01" {
			t.Errorf("Expected TFD FechaTimbrado 2025-01-15T10:30:01, got %s", tfd.FechaTimbrado)
		}
		if tfd.RfcProvCert != "AAA010101AAA" {
			t.Errorf("Expected TFD RfcProvCert AAA010101AAA, got %s", tfd.RfcProvCert)
		}
		if tfd.SelloCFD != "SELLO_CFD_DE_EJEMPLO_1234567890_==" {
			t.Errorf("Expected TFD SelloCFD mismatch, got %s", tfd.SelloCFD)
		}
		if tfd.NoCertificadoSAT != "30001000000300023789" {
			t.Errorf("Expected TFD NoCertificadoSAT 30001000000300023789, got %s", tfd.NoCertificadoSAT)
		}
		if tfd.SelloSAT != "SELLO_SAT_DE_EJEMPLO_1234567890_==" {
			t.Errorf("Expected TFD SelloSAT mismatch, got %s", tfd.SelloSAT)
		}
	})

	t.Run("UseConcepts - Parsea conceptos sin impuestos", func(t *testing.T) {
		// Habilitamos el parseo de conceptos pero NO de sus impuestos
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UseConcepts()
		data, err := handler.TransformFromString(xmlStr)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		// Validamos que se hayan parseado los conceptos
		if len(data.CFDI40.Conceptos) == 0 {
			t.Fatal("Se esperaban conceptos, no se obtuvo ninguno")
		}

		c := data.CFDI40.Conceptos[0]
		if c.ClaveProdServ != "84111506" {
			t.Errorf("Se esperaba ClaveProdServ 84111506, se obtuvo %s", c.ClaveProdServ)
		}
		if c.Descripcion != "SERVICIO DE EJEMPLO" {
			t.Errorf("Se esperaba Descripcion SERVICIO DE EJEMPLO, se obtuvo %s", c.Descripcion)
		}

		// Validamos que NO se hayan parseado los impuestos del concepto (por defecto)
		if len(c.Traslados) > 0 {
			t.Errorf("No se esperaban traslados en el concepto, se obtuvieron %d", len(c.Traslados))
		}
	})

	t.Run("UseConceptsWithTaxes - Parsea conceptos con impuestos", func(t *testing.T) {
		// Habilitamos el parseo de conceptos Y de sus impuestos
		// Nota: UseConceptsWithTaxes debería usarse junto con UseConcepts
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UseConcepts().UseConceptsWithTaxes()
		data, err := handler.TransformFromString(xmlStr)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		if len(data.CFDI40.Conceptos) == 0 {
			t.Fatal("Se esperaban conceptos, no se obtuvo ninguno")
		}

		c := data.CFDI40.Conceptos[0]
		// Validamos que se hayan parseado los impuestos del concepto
		if len(c.Traslados) == 0 {
			t.Fatal("Se esperaban traslados en el concepto, no se obtuvo ninguno")
		}

		tr := c.Traslados[0]
		if tr.Base != "1000.00" {
			t.Errorf("Se esperaba Base de traslado 1000.00, se obtuvo %s", tr.Base)
		}
		if tr.Impuesto != "002" {
			t.Errorf("Se esperaba Impuesto de traslado 002, se obtuvo %s", tr.Impuesto)
		}
		if tr.TasaOCuota != "0.160000" {
			t.Errorf("Se esperaba TasaOCuota de traslado 0.160000, se obtuvo %s", tr.TasaOCuota)
		}
	})

	t.Run("UseRelatedCFDIs - Parsea CFDIs relacionados", func(t *testing.T) {
		// Habilitamos el parseo de CFDIs relacionados
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UseRelatedCFDIs()
		data, err := handler.TransformFromString(xmlStr)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		// Validamos que se hayan parseado los CFDIs relacionados
		if len(data.CFDI40.CFDIsRelacionados) == 0 {
			t.Fatal("Se esperaban CFDIs relacionados, no se obtuvo ninguno")
		}

		rel := data.CFDI40.CFDIsRelacionados[0]
		if rel.TipoRelacion != "04" {
			t.Errorf("Se esperaba TipoRelacion 04, se obtuvo %s", rel.TipoRelacion)
		}
		if rel.UUID != "5FB2822E-396D-4725-8521-500FAB000222" {
			t.Errorf("Se esperaba UUID relacionado 5FB2822E-396D-4725-8521-500FAB000222, se obtuvo %s", rel.UUID)
		}
	})
}
