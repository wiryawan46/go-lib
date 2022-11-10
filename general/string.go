package general

import "strings"

// GetMimeFromImage
// Split image extension using explode
func GetMimeFromImage(base64Image string) string {
	explodeList := []string{",", ";", ":", "/"}
	var result []string

	for i := 0; i < len(explodeList); i++ {
		step := explodeList[i]
		if i == 0 {
			result = Explode(step, base64Image)
		}
		if step == "/" {
			result = Explode(step, result[1])
		} else {
			result = Explode(step, result[0])
		}
	}

	return "." + result[1]
}

// Explode string
func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	}
	return strings.Split(text, delimiter)
}
