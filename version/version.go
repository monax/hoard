package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// This is the single authoritative version. Should be in semantic version form.
const version = "0.0.1"

const (
	// Base of minor, major, and patch version numbers
	numberBase = 10
	// Number of bits to represent version numbers
	uintBits = 8
)

var major, minor, patch uint8

func init() {
	var err error
	major, minor, patch, err = ParseVersion(version)

	if err != nil {
		panic(fmt.Errorf("Could not parse version: '%s'", version))
	}
}

func String() string {
	return version
}

func Major() uint8 {
	return major
}

func Minor() uint8 {
	return minor
}

func Patch() uint8 {
	return patch
}

func ParseVersion(ver string) (uint8, uint8, uint8, error) {
	parts := strings.Split(ver, ".")
	if len(parts) != 3 {
		return 0, 0, 0,
			errors.New("Version string must have three '.' separated parts")
	}
	maj, err := strconv.ParseUint(parts[0], numberBase, uintBits)
	if err != nil {
		return 0, 0, 0, err
	}
	min, err := strconv.ParseUint(parts[1], numberBase, uintBits)
	if err != nil {
		return 0, 0, 0, err
	}
	pat, err := strconv.ParseUint(parts[2], numberBase, uintBits)
	if err != nil {
		return 0, 0, 0, err
	}
	return uint8(maj), uint8(min), uint8(pat), err
}
