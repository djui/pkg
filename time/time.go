package time

import (
	"encoding/xml"
	"strings"
	"time"
)

var formatReplacer *strings.Replacer
var formatPatterns = map[string]string{
	// Day
	"%d": "02",     // Two-digit day of the month (with leading zeros)
	"%e": "2",      // Day of the month, with a single digits
	"%a": "Mon",    // An abbreviated textual representation of the day
	"%A": "Monday", // A full textual representation of the day
	// Month
	"%b": "Jan",     // Abbreviated month name
	"%h": "Jan",     // Abbreviated month name (an alias of %b)
	"%B": "January", // Full month name
	"%m": "01",      // Two digit representation of the month
	// Year
	"%y": "06",   // Two digit representation of the year
	"%Y": "2006", // Four digit representation for the year
	// Time
	"%H": "15",          // Two digit representation of the hour in 24-hour format
	"%k": "15",          // Two digit representation of the hour in 24-hour format, with a space preceding single digits
	"%I": "03",          // Two digit representation of the hour in 12-hour format
	"%l": "3",           // Hour in 12-hour format, with a space preceding single digits
	"%M": "04",          // Two digit representation of the minute
	"%P": "pm",          // UPPER-CASE 'AM' or 'PM' based on the given time
	"%p": "PM",          // lower-case 'am' or 'pm' based on the given time
	"%r": "03:04:05 PM", // Same as "%I:%M:%S %p"
	"%R": "15:04",       // Same as "%H:%M"
	"%S": "05",          // Two digit representation of the second
	"%T": "15:04:05",    // Same as "%H:%M:%S"
	"%z": "-0700",       // The time zone offset. Not implemented as described on Windows. See below for more information.
	"%Z": "MST",         // The time zone abbreviation.
	// Time and Date Stamps
	"%D": "01/02/2006", // Same as "%m/%d/%y"
	"%F": "2006-01-02", // Same as "%Y-%m-%d"
	// Miscellaneous
	"%n": "\n", // A newline character ("\n")
	"%t": "\t", // A Tab character ("\t")
	"%%": "%",  // A literal percentage character ("%")
}

func init() {
	var fp []string
	for k, v := range formatPatterns {
		fp = append(fp, k, v)
	}
	formatReplacer = strings.NewReplacer(fp...)
}

// ParseUnix parses a formatted string and returns the time value it represents.
// The layout defines the format declared in
// http://pubs.opengroup.org/onlinepubs/007908799/xsh/strftime.html . The
// implementation differs from the original as no local information is taken
// into account.
func ParseUnix(layout, value string) (time.Time, error) {
	layout = formatReplacer.Replace(layout)
	return time.ParseInLocation(layout, value, time.Local)
}

// FormatUnix returns a textual representation of the time value formatted
// according to layout, which defines the format declared in
// http://pubs.opengroup.org/onlinepubs/007908799/xsh/strftime.html . The
// implementation differs from the original as no local information is taken
// into account.
func FormatUnix(t time.Time, layout string) string {
	layout = formatReplacer.Replace(layout)
	return t.Format(layout)
}

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
