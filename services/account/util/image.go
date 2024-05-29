package util

import "strings"

func ExtractImageType(contentType string) (string, error) {
	allowedTypes := map[string]struct{}{
		"image/png":  {},
		"image/jpg":  {},
		"image/jpeg": {},
	}
	if _, exists := allowedTypes[contentType]; !exists {
		return "", NewBadRequest("Unallowed content type")
	}

	splitedType := strings.Split(contentType, "/")
	if len(splitedType) < 2 {
		return "", NewBadRequest("Could not parse the content type")
	}

	return splitedType[1], nil
}
