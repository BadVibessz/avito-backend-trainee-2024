package _map

type AnyMap = map[string]any

func GetIntFromAnyMap(anyMap AnyMap, key string) (int, error) {
	valAny, exists := anyMap[key]
	if !exists {
		return -1, ErrNoSuchKey
	}

	val, ok := valAny.(float64)
	if !ok {
		return -1, ErrCannotConver
	}

	return int(val), nil
}

func GetStringFromAnyMap(anyMap AnyMap, key string) (string, error) {
	valAny, exists := anyMap[key]
	if !exists {
		return "", ErrNoSuchKey
	}

	val, ok := valAny.(string)
	if !ok {
		return "", ErrCannotConver
	}

	return val, nil
}

func GetBoolFromAnyMap(anyMap AnyMap, key string) (bool, error) {
	valAny, exists := anyMap[key]
	if !exists {
		return false, ErrNoSuchKey
	}

	val, ok := valAny.(bool)
	if !ok {
		return false, ErrCannotConver
	}

	return val, nil
}
