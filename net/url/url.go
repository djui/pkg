package url

import "net/url"

// AppendQuery appends additional query values to a given URL's query.
func AppendQuery(u url.URL, q url.Values) url.URL {
	uq := u.Query()
	for k, vv := range q {
		for _, v := range vv {
			uq.Add(k, v)
		}
	}
	u.RawQuery = uq.Encode()
	return u
}
