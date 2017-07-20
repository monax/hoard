package release

import (
	"bytes"
	"html/template"
)

var changeLogTemplate = template.Must(template.New("changelog_template").
	Parse(`# Hoard Changelog{{ range . }}
## Version {{ .Version }}
{{ .Changes }}
{{ end }}`))

func Changes() string {
	return hoardReleases[0].Changes
}

func Changelog() string {
	completeChangeLog, err := changelogForReleases(hoardReleases)
	if err != nil {
		panic(err)
	}
	return completeChangeLog
}

func changelogForReleases(rels []Release) (string, error) {
	buf := new(bytes.Buffer)
	err := changeLogTemplate.Execute(buf, rels)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
