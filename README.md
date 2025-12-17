# gocfdi-transform

Package `gocfdi_transform` proporciona handlers para transformar un CFDI (Comprobante Fiscal Digital por Internet) en una estructura de datos fácil de manejar en Go.

Este proyecto fue inspirado por [pycfdi-transform](https://github.com/swsapien/pycfdi-transform), con el objetivo de portar su funcionalidad al ecosistema de Go.

## Características

- Soporte para **CFDI 4.0**.
- Extracción de datos de **Timbre Fiscal Digital (TFD) 1.1**.
- Soporte para complementos:
    - **Nómina 1.2**
    - **Pagos 2.0**
    - **Venta de Vehículos**

## Instalación

```bash
go get github.com/sucksens/gocfdi-transform
```

## Uso

### Ejemplo Básico (CFDI 4.0)

```go
package main

import (
	"fmt"
	"log"

	"github.com/sucksens/gocfdi-transform/sax"
)

func main() {
	// Crear un nuevo handler con configuración por defecto
	handler := sax.NewCFDI40Handler(sax.NewDefaultConfig())
    
	// Transformar desde un archivo XML
	data, err := handler.TransformFromFile("factura.xml")
	if err != nil {
		log.Fatal(err)
	}

	// Acceder a los datos del CFDI
	fmt.Printf("UUID: %s\n", data.TFD11[0].UUID)
	fmt.Printf("Total: %s\n", data.CFDI40.Total)
	fmt.Printf("Emisor: %s\n", data.CFDI40.Emisor.Nombre)
}
```

### Habilitar Complementos (Ej. Nómina 1.2)

Para extraer información de complementos específicos, debes habilitarlos en el handler:

```go
package main

import (
	"fmt"
	"log"

	"github.com/sucksens/gocfdi-transform/sax"
)

func main() {
	// Habilitar el complemento de Nómina 1.2
	handler := sax.NewCFDI40Handler(sax.NewDefaultConfig()).UseNomina12()
    
	data, err := handler.TransformFromFile("recibo_nomina.xml")
	if err != nil {
		log.Fatal(err)
	}

	if len(data.Nomina12) > 0 {
		nomina := data.Nomina12[0]
		fmt.Printf("Tipo Nómina: %s\n", nomina.TipoNomina)
		fmt.Printf("Fecha Pago: %s\n", nomina.FechaPago)
		fmt.Printf("Total Percepciones: %s\n", nomina.TotalPercepciones)
	}
}
```

## Estructura de Datos (Referencia 4.0)

A continuación se muestra una representación JSON de cómo se ve una estructura `CFDI40Data` completa (habilitando todos los complementos soportados):

```json
{
  "cfdi40": {
    "version": "4.0",
    "serie": "A",
    "folio": "12345",
    "fecha": "2023-01-01T12:00:00",
    "subtotal": "1000.00",
    "total": "1160.00",
    "moneda": "MXN",
    "emisor": {
      "rfc": "AAA010101AAA",
      "nombre": "EMISOR DE PRUEBA",
      "regimen_fiscal": "601"
    },
    "receptor": {
      "rfc": "XAXX010101000",
      "nombre": "PUBLICO EN GENERAL",
      "uso_cfdi": "G03"
    },
    "conceptos": [
      {
        "clave_prod_serv": "84111506",
        "descripcion": "SERVICIO",
        "valor_unitario": "1000.00",
        "importe": "1000.00",
        "traslados": [
            {
                "impuesto": "002",
                "tipo_factor": "Tasa",
                "tasa_o_cuota": "0.160000",
                "importe": "160.00"
            }
        ]
      }
    ]
  },
  "tfd11": [
      {
          "uuid": "5FB2822E-396D-4725-8521-500FAB000222",
          "fecha_timbrado": "2023-01-01T12:05:00"
      }
  ],
  "nomina_12": [
      {
          "version": "1.2",
          "tipo_nomina": "O",
          "fecha_pago": "2023-01-15",
          "emisor": { "registro_patronal": "Y54545454" },
          "receptor": { "num_empleado": "123", "puesto": "DESAROLLADOR" }
      }
  ],
  "pagos20": [
      {
          "totales": { "monto_total_pagos": "500.00" },
          "pagos": [
              {
                  "fecha_pago": "2023-01-10T10:00:00",
                  "forma_de_pago_p": "03",
                  "monto": "500.00"
              }
          ]
      }
  ]
}
```

## Estatus de Integraciones

Tabla de soporte de versiones de CFDI y complementos:

| Integración | Versión | Estatus |
| :--- | :--- | :--- |
| **CFDI** | 4.0 | ✅ Implementado |
| **Timbre Fiscal Digital (TFD)** | 1.1 | ✅ Implementado |
| **Nómina** | 1.2 | ✅ Implementado |
| **Pagos** | 2.0 | ✅ Implementado |
| **Venta de Vehículos** | 1.1 | ✅ Implementado |
| **Carta Porte** | 3.0 / 3.1 | ❌ Pendiente |
| **Impuestos Locales** | 1.1 | ❌ Pendiente |
| **Comercio Exterior** | 2.0 | ❌ Pendiente |
| **Retenciones** | 2.0 | ❌ Pendiente |


## Licencia

Este proyecto está bajo la Licencia **GPLv3**. Consulta el archivo [LICENSE](LICENSE) para más detalles.
