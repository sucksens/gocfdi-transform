# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2025-02-17

### Added
- `ModelBuilder` para extracción eficiente de atributos XML
- `BaseHandler` con funcionalidades comunes compartidas entre handlers
- Utilidades de procesamiento XML en `xml_utils.go`
- Validación de versión integrada en `ValidateVersion`

### Changed
- Refactorización completa de handlers para integrar con `BaseHandler` y `ModelBuilder`
- `CFDI40Handler` - Integración con arquitectura unificada
- `Nomina12Handler` - Integración con arquitectura unificada
- `Pagos20Handler` - Integración con arquitectura unificada
- `VentaVehiculos11Handler` - Integración con arquitectura unificada
- `TFD11Handler` - Integración con arquitectura unificada
- `ImpuestosLocales10Handler` - Integración con arquitectura unificada
- Consistencia en la lógica de `TransformFromString` en VentaVehículos11
- Integración de `ValidateVersion` en TFD handler
- Eliminación de código muerto de transformación de atributos en `BaseHandler`

### Fixed
- Error en `ElementProcessor` para consistencia con `ProcessElementUntil`

## [0.5.0] - Previous Release

### Added
- Soporte inicial para CFDI 4.0
- Extracción de datos de Timbre Fiscal Digital (TFD) 1.1
- Complemento Nómina 1.2
- Complemento Pagos 2.0
- Complemento Venta de Vehículos 1.1
- Complemento Impuestos Locales 1.0
- Funciones `Use` para habilitar/deshabilitar complementos específicos

[0.6.0]: https://github.com/sucksens/gocfdi-transform/releases/tag/v0.6.0
[0.5.0]: https://github.com/sucksens/gocfdi-transform/releases/tag/v0.5.0
