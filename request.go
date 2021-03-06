package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Request struct {
	Filename string
	Content  string
	Command  string
	Image    string
	Format   string
}

var FilenameRegexp = regexp.MustCompile(`\A([a-z\d\-\_]+)\.[a-z]{1,6}\z`)

func normalizeString(val string) string {
	return strings.ToLower(strings.TrimSpace(val))
}

func ParseRequest(r *http.Request) (*Request, error) {
	req := Request{
		Filename: normalizeString(r.FormValue("filename")),
		Command:  normalizeString(r.FormValue("command")),
		Content:  r.FormValue("content"),
	}

	if req.Filename == "" {
		return nil, fmt.Errorf("Filename is required")
	}

	if !FilenameRegexp.Match([]byte(req.Filename)) {
		return nil, fmt.Errorf("Invalid filename")
	}

	if req.Content == "" {
		return nil, fmt.Errorf("Content is required")
	}

	lang, err := GetLanguageConfig(req.Filename)
	if err != nil {
		return nil, err
	}

	req.Image = lang.Image
	req.Format = lang.Format

	if req.Format == "" {
		req.Format = `text/plain; charset="UTF-8"`
	}

	if req.Command == "" {
		req.Command = fmt.Sprintf(lang.Command, req.Filename)
	}

	return &req, nil
}
