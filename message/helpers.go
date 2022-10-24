package message

import "strconv"

func parseBool(key string, values map[string]string) (bool, error) {
	value, exists := values[key]
	if !exists {
		return false, nil
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return result, nil
}

func parseInt(key string, values map[string]string) (int, error) {
	value, exists := values[key]
	if !exists {
		return 0, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func firstValues(form map[string][]string) map[string]string {
	result := map[string]string{}
	for key, values := range form {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}
