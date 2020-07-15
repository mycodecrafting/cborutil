package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	cbor "github.com/brianolson/cbor_go"
)

type Container struct {
	object interface{}
}

func DotPathToSlice(path string) []string {
	hierarchy := strings.Split(path, ".")
	return hierarchy
}

func (c *Container) Search(hierarchy ...string) (*Container, error) {
	object := c.object
	for i := 0; i < len(hierarchy); i++ {
		pathSeg := hierarchy[i]
		if imap, ok := object.(map[interface{}]interface{}); ok {
			object, ok = imap[pathSeg]
			if !ok {
				return nil, fmt.Errorf("failed to resolve path segment '%v': key '%v' was not found", i, pathSeg)
			}
		} else if mmap, ok := object.(map[string]interface{}); ok {
			object, ok = mmap[pathSeg]
			if !ok {
				return nil, fmt.Errorf("failed to resolve path segment '%v': key '%v' was not found", i, pathSeg)
			}
		} else if marray, ok := object.([]interface{}); ok {
			index, err := strconv.Atoi(pathSeg)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve path segment '%v': found array but segment value '%v' could not be parsed into array index: %v", i, pathSeg, err)
			}
			if len(marray) <= index {
				return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' exceeded target array size of '%v'", i, pathSeg, len(marray))
			}
			object = marray[index]
		} else {
			return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was not found", i, pathSeg)
		}
	}
	return &Container{object}, nil
}

func (c *Container) Set(value interface{}, hierarchy ...string) (*Container, error) {
	if len(hierarchy) == 0 {
		c.object = value
		return c, nil
	}

	if c.object == nil {
		c.object = map[string]interface{}{}
	}
	object := c.object

	for i := 0; i < len(hierarchy); i++ {
		pathSeg := hierarchy[i]
		if imap, ok := object.(map[interface{}]interface{}); ok {
			if i == len(hierarchy)-1 {
				object = value
				imap[pathSeg] = object
			} else if object = imap[pathSeg]; object == nil {
				imap[pathSeg] = map[string]interface{}{}
				object = imap[pathSeg]
			}
		} else if mmap, ok := object.(map[string]interface{}); ok {
			if i == len(hierarchy)-1 {
				object = value
				mmap[pathSeg] = object
			} else if object = mmap[pathSeg]; object == nil {
				mmap[pathSeg] = map[string]interface{}{}
				object = mmap[pathSeg]
			}
		} else if marray, ok := object.([]interface{}); ok {
			index, err := strconv.Atoi(pathSeg)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve path segment '%v': found array but segment value '%v' could not be parsed into array index: %v", i, pathSeg, err)
			}
			if len(marray) <= index {
				return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' exceeded target array size of '%v'", i, pathSeg, len(marray))
			}
			if i == len(hierarchy)-1 {
				object = value
				marray[index] = object
			} else if object = marray[index]; object == nil {
				return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was not found", i, pathSeg)
			}
		}
	}
	return &Container{object}, nil
}

func ParseStr(s string, b64 bool) (*Container, error) {
	if b64 {
		return ParseBase64(s)
	} else {
		return ParseHex(s)
	}
}

func ParseHex(s string) (*Container, error) {
	b, _ := hex.DecodeString(s)
	return LoadBytes(b)
}

func ParseBase64(s string) (*Container, error) {
	b, _ := base64.StdEncoding.DecodeString(s)
	return LoadBytes(b)
}

func LoadBytes(b []byte) (*Container, error) {
	var c Container
	if err := cbor.Loads(b, &c.object); err != nil {
		return nil, err
	}
	return &c, nil
}

func Encode(c *Container, b64 bool) (string, error) {
	cborBytes, err := cbor.Dumps(c.object)
	if err != nil {
		return "", err
	}
	return EncodeToStr(cborBytes, b64), nil
}

func EncodeToStr(b []byte, b64 bool) string {
	if b64 {
		return base64.StdEncoding.EncodeToString(b)
	} else {
		return hex.EncodeToString(b)
	}
}

func DecodePath(cborStr string, path string, isBase64 bool) (interface{}, error) {
	c, err := ParseStr(cborStr, isBase64)
	if err != nil {
		return "", err
	}
	if path == "" {
		return c.object, nil
	}
	c2, err := c.Search(DotPathToSlice(path)...)
	return c2.object, err
}

func UpdatePath(cborStr string, path string, data interface{}, isBase64 bool) (string, error) {
	c, err := ParseStr(cborStr, isBase64)
	if err != nil {
		return "", err
	}
	c.Set(data, DotPathToSlice(path)...)
	return Encode(c, isBase64)
}
