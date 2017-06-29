package shared

import (
	"fmt"
	"os"
	"strings"
)

func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func SplitListenURL(listenOn string) (string, string, error) {
	// net.Listen does not want a parsed url.URL so it seems to make more sense
	// just to do a dumb split here to support the various networks
	listenParts := strings.Split(listenOn, "://")
	if len(listenParts) != 2 {
		return "", "", fmt.Errorf("Expected a Go net.Listen URL of the form "+
			"'<net>://<laddr>', but got: '%s'", listenOn)
	}
	if listenParts[0] == "" {
		return "", "", fmt.Errorf("Expected the URL scheme to be present, "+
			"but got '%s'", listenOn)
	}
	if listenParts[1] == "" {
		return "", "", fmt.Errorf("Expected the URL host to be present, "+
			"but got '%s'", listenOn)
	}
	return listenParts[0], listenParts[1], nil
}
