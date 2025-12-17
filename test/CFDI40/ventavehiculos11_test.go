package test

import (
	"testing"

	"github.com/sucksens/gocfdi-transform/sax"
)

func TestVentaVehiculos11Handler(t *testing.T) {
	xmlStr := `
	<cfdi:Comprobante xmlns:cfdi="http://www.sat.gob.mx/cfd/4" xmlns:ventavehiculos="http://www.sat.gob.mx/ventavehiculos" Version="4.0">
		<cfdi:Complemento>
			<ventavehiculos:VentaVehiculos Version="1.1" ClaveVehicular="123456" Niv="ABC1234567890">
				<ventavehiculos:InformacionAduanera Numero="123" Fecha="2023-01-01" Aduana="Aduana1"/>
				<ventavehiculos:Parte Cantidad="1" Descripcion="Parte1">
					<ventavehiculos:InformacionAduanera Numero="456" Fecha="2023-01-02" Aduana="Aduana2"/>
				</ventavehiculos:Parte>
			</ventavehiculos:VentaVehiculos>
		</cfdi:Complemento>
	</cfdi:Comprobante>
	`

	t.Run("Parse VentaVehiculos11", func(t *testing.T) {
		handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UseVentaVehiculos11()
		data, err := handler.TransformFromString(xmlStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(data.VentaVehiculos11) == 0 {
			t.Fatal("Expected VentaVehiculos11 data, got none")
		}

		vv := data.VentaVehiculos11[0]

		// Verify Root Attributes
		if vv.Version != "1.1" {
			t.Errorf("Expected Version 1.1, got %s", vv.Version)
		}
		if vv.ClaveVehicular != "123456" {
			t.Errorf("Expected ClaveVehicular 123456, got %s", vv.ClaveVehicular)
		}
		if vv.Niv != "ABC1234567890" {
			t.Errorf("Expected Niv ABC1234567890, got %s", vv.Niv)
		}

		// Verify InformacionAduanera (Root)
		if len(vv.InformacionAduanera) != 1 {
			t.Fatalf("Expected 1 InformacionAduanera at root, got %d", len(vv.InformacionAduanera))
		}
		ia := vv.InformacionAduanera[0]
		if ia.Numero != "123" {
			t.Errorf("Expected InformacionAduanera Numero 123, got %s", ia.Numero)
		}
		if ia.Aduana != "Aduana1" {
			t.Errorf("Expected InformacionAduanera Aduana Aduana1, got %s", ia.Aduana)
		}

		// Verify Parte
		if len(vv.Partes) != 1 {
			t.Fatalf("Expected 1 Parte, got %d", len(vv.Partes))
		}
		p := vv.Partes[0]
		if p.Cantidad != "1" {
			t.Errorf("Expected Parte Cantidad 1, got %s", p.Cantidad)
		}
		if p.Descripcion != "Parte1" {
			t.Errorf("Expected Parte Descripcion Parte1, got %s", p.Descripcion)
		}

		// Verify InformacionAduanera (Parte)
		if len(p.InformacionAduanera) != 1 {
			t.Fatalf("Expected 1 InformacionAduanera in Parte, got %d", len(p.InformacionAduanera))
		}
		pia := p.InformacionAduanera[0]
		if pia.Numero != "456" {
			t.Errorf("Expected Parte InformacionAduanera Numero 456, got %s", pia.Numero)
		}
	})
}
