package time

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

// RFC3339Time allows RFC 3339 compliant XML un/marshaling.
type RFC3339Time struct {
	time.Time
}

// MarshalJSON satisfies the json.Marshaler interface.
func (c *RFC3339Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Time.Format(time.RFC3339))
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (c *RFC3339Time) UnmarshalJSON(json []byte) error {
	parse, err := time.Parse(time.RFC3339, string(json))
	if err != nil {
		return err
	}
	c.Time = parse
	return nil
}

// MarshalXML satisfies the xml.Marshaler interface.
func (c *RFC3339Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(c.Time.Format(time.RFC3339), start)
}

// MarshalXMLAttr satsfies the xml.MarshalerAttr interface.
func (c *RFC3339Time) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	attr := xml.Attr{
		Name:  name,
		Value: c.Time.Format(time.RFC3339),
	}
	return attr, nil
}

// UnmarshalXML satisfies the xml.Unmarshaler interface.
func (c *RFC3339Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	parse, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}
	c.Time = parse
	return nil
}

// UnmarshalXMLAttr satisfies the xml.UnmarshalerAttr interface.
func (c *RFC3339Time) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, err := time.Parse(time.RFC3339, attr.Value)
	if err != nil {
		return err
	}
	c.Time = parse
	return nil
}

// MarshalBinary satisfies the encoding.BinaryMarshaler interface.
func (c *RFC3339Time) MarshalBinary() ([]byte, error) {
	return []byte(c.Time.Format(time.RFC3339)), nil
}

// UnmarshalBinary satisfies the encoding.BinaryUnmarshaler interface.
func (c *RFC3339Time) UnmarshalBinary(data []byte) error {
	parse, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return err
	}
	c.Time = parse
	return nil

}

// MarshalText satisfies the encoding.TextMarshaler interface.
func (c *RFC3339Time) MarshalText() ([]byte, error) {
	return []byte(c.Time.Format(time.RFC3339)), nil
}

// UnmarshalText satisfies the encoding.TextUnmarshaler interface.
func (c *RFC3339Time) UnmarshalText(text []byte) error {
	parse, err := time.Parse(time.RFC3339, string(text))
	if err != nil {
		return err
	}
	c.Time = parse
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
