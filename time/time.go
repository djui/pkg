package time

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

// RFC3339Time allows RFC 3339 compliant un/marshaling.
type RFC3339Time struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.RFC3339)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.RFC3339)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *RFC3339Time) UnmarshalJSON(json []byte) error {
	parse, err := time.Parse(`"`+time.RFC3339+`"`, string(json))
	if err != nil {
		return err
	}
	t.Time = parse
	return nil
}

// MarshalXML implements the xml.Marshaler interface.
func (t RFC3339Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(t.Time.Format(time.RFC3339), start)
}

// MarshalXMLAttr implements the xml.MarshalerAttr interface.
func (t RFC3339Time) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	attr := xml.Attr{
		Name:  name,
		Value: t.Time.Format(time.RFC3339),
	}
	return attr, nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (t *RFC3339Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	parse, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}
	t.Time = parse
	return nil
}

// UnmarshalXMLAttr implements the xml.UnmarshalerAttr interface.
func (t *RFC3339Time) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, err := time.Parse(time.RFC3339, attr.Value)
	if err != nil {
		return err
	}
	t.Time = parse
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (t RFC3339Time) MarshalBinary() ([]byte, error) {
	return []byte(t.Time.Format(time.RFC3339)), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *RFC3339Time) UnmarshalBinary(data []byte) error {
	parse, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return err
	}
	t.Time = parse
	return nil

}

// MarshalText implements the encoding.TextMarshaler interface.
func (t RFC3339Time) MarshalText() ([]byte, error) {
	return []byte(t.Time.Format(time.RFC3339)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *RFC3339Time) UnmarshalText(text []byte) error {
	parse, err := time.Parse(time.RFC3339, string(text))
	if err != nil {
		return err
	}
	t.Time = parse
	return nil
}

// GobEncode implements the gob.GobEncoder interface.
func (t RFC3339Time) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface.
func (t *RFC3339Time) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}
