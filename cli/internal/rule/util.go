package rule

func getString(key string, raw map[string]interface{}, def string) string {
	value, ok := raw[key]
	if !ok {
		return def
	}
	str, ok := value.(string)
	if !ok {
		return def
	}
	return str
}

func getBool(key string, raw map[string]interface{}, def bool) bool {
	value, ok := raw[key]
	if !ok {
		return def
	}
	b, ok := value.(bool)
	if !ok {
		return def
	}
	return b
}
