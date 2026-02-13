package sax

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/sucksens/gocfdi-transform/models"
)

// TFD11Handler handles parsing of Timbre Fiscal Digital 1.1 complement.
type TFD11Handler struct {
	*BaseHandler
	builder *ModelBuilder
}

// NewTFD11Handler creates a new TFD11Handler.
func NewTFD11Handler(config HandlerConfig) *TFD11Handler {
	return &TFD11Handler{
		BaseHandler: NewBaseHandler(config),
		builder:     NewModelBuilder(config),
	}
}

// TransformFromString parses a TFD 1.1 XML string.
func (h *TFD11Handler) TransformFromString(xmlStr string) (*models.TFD11, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "TimbreFiscalDigital" {
				return h.transformTFD(se)
			}
		}
	}

	return nil, errors.New("TimbreFiscalDigital element not found")
}

func (h *TFD11Handler) transformTFD(se xml.StartElement) (*models.TFD11, error) {
	version := getAttrValue(se, "Version")
	if version != "1.1" {
		return nil, errors.New("incorrect type of TFD, this handler only supports TFD version 1.1")
	}

	return &models.TFD11{
		Version:          version,
		NoCertificadoSAT: h.builder.ExtractString(se, "NoCertificadoSAT"),
		UUID:             h.builder.ExtractUpper(se, "UUID"),
		FechaTimbrado:    h.builder.ExtractString(se, "FechaTimbrado"),
		RfcProvCert:      h.builder.ExtractString(se, "RfcProvCertif"),
		SelloCFD:         h.builder.ExtractCompact(se, "SelloCFD"),
		SelloSAT:         h.builder.ExtractCompact(se, "SelloSAT"),
	}, nil
}
