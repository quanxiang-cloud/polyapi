package apipath

import (
	"fmt"
	"path"
	"strings"
)

// GetAPIType get api type from api name separated by '.'
// eg: "/system/api.r" => "r"
func GetAPIType(fullPath string) string {
	for i := len(fullPath) - 1; i >= 0; i-- {
		switch c := fullPath[i]; c {
		case '/':
			return "" //do not find '.'
		case '.':
			return fullPath[i+1:]
		}
	}
	return ""
}

// GenerateAPIName generate api name by apiType
// eg: ("api","r") => "api.r"
func GenerateAPIName(name string, apiType string) string {
	if !strings.ContainsRune(name, '.') {
		return fmt.Sprintf("%s.%s", name, apiType)
	}
	return name
}

// BaseName get api base name without type suffix.
// eg: /system/foo.r => foo
func BaseName(apiPath string) string {
	name := Name(apiPath)
	if index := strings.IndexRune(name, '.'); index >= 0 {
		return name[:index]
	}
	return name
}

// Split parse full path with path and name
func Split(full string) (string, string) {
	if !strings.HasPrefix(full, "/") {
		full = "/" + full
	}
	path, name := "", full
	if index := strings.LastIndex(full, "/"); index >= 0 {
		if path = full[:index]; path == "" {
			path = "/"
		}
		name = full[index+1:]
	}
	return path, name
}

// Name get the last name of a full path
func Name(fullPath string) string {
	_, name := Split(fullPath)
	return name
}

// Parent get the parent path of a full path
func Parent(fullPath string) string {
	parent, _ := Split(fullPath)
	return parent
}

// Join join the namespace and name as full path
func Join(namespace, name string) string {
	if !strings.HasPrefix(namespace, "/") {
		namespace = "/" + namespace
	}
	return path.Join(namespace, name)
}

// Format convert full path as standard format
func Format(full string) string {
	path, name := Split(full)
	return Join(path, name)
}

// FormatPrefix format path with sub path
func FormatPrefix(fullPath string) string {
	fullPath = Format(fullPath)
	return fullPath + "/*"
}

// MakeRequestURL generate request url by parts
func MakeRequestURL(schema, host, apiPath string) string {
	if len(apiPath) == 0 || apiPath[0] != '/' {
		apiPath = "/" + apiPath
	}
	protocol := strings.ToLower(schema) + "://"
	return protocol + host + apiPath
}
