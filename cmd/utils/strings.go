package utils

import (
    "regexp"
    "strings"
)

func ToSnakeCase(str string) string {
    re := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
    str = re.ReplaceAllString(str, "${1}_${2}")

    re = regexp.MustCompile(`([a-z\d])([A-Z])`)
    str = re.ReplaceAllString(str, "${1}_${2}")

    // Convert the entire string to lowercase
    return strings.ToLower(str)
}

func Contains[T comparable](slice []T, value T) bool {
    for _, v := range slice {
        if v == value {
            return true
        }
    }

    return false
}