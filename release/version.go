package release

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Base of minor, major, and patch version numbers
	numberBase = 10
	// Number of bits to represent version numbers
	uintBits = 8
)

var major, minor, patch uint8

func init() {
	var err error
	major, minor, patch, err = ParseVersion(Version())

	if err != nil {
		panic(fmt.Errorf("Could not parse version: '%s'", Version()))
	}
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
			fmt.Errorf("Version string must have three '.' separated parts "+
				"but '%s' does not.", ver)
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
