package aggregdepscore

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var gopkginRegexp = regexp.MustCompile(`^gopkg\.in\/(?<package>[a-zA-Z0-9_\-./]+)\.v[0-9]+$`)

func getGopkginRepository(packageName string) (string, error) {
	// note: a more "correct" way would be to send a request to gopkg.in
	// and look for a "<meta name='go-import'>" tag in the response
	// (see https://pkg.go.dev/cmd/go#hdr-Remote_import_paths)
	// but using a regexp is simpler and saves us a network request

	// we follow the logic in https://labix.org/gopkg.in

	matches := gopkginRegexp.FindStringSubmatch(packageName)
	if len(matches) == 0 {
		return "", errors.New("not a gopkg.in package")
	}

	p := matches[gopkginRegexp.SubexpIndex("package")]

	parts := strings.Split(p, "/")

	if len(parts) == 2 {
		return fmt.Sprintf("github.com/%s/%s", parts[0], parts[1]), nil
	}

	if len(parts) == 1 {
		return fmt.Sprintf("github.com/go-%s/%s", parts[0], parts[0]), nil
	}

	return "", fmt.Errorf("unexpected number of parts in gopkg.in package name: %s", p)
}
