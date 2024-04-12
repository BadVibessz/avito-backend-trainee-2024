package url

import "net/url"

func RemoveQueryParamByKey(url url.URL, key string) *url.URL {
	query := url.Query()
	query.Del(key)

	url.RawQuery = query.Encode()

	return &url
}
