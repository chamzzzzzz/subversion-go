package subversion

import (
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
)

type LogEntry struct {
	Revision string `xml:"revision,attr"`
	Author   string `xml:"author"`
	Date     string `xml:"date"`
	Message  string `xml:"msg"`
}

type Log struct {
	Entries []LogEntry `xml:"logentry"`
}

type Client struct {
	Remote   string
	Local    string
	Username string
	Password string
}

func (c *Client) Log(limit int, revision string, search []string, searchand []string) (*Log, error) {
	args := []string{"log", "--xml"}
	if c.Username != "" && c.Password != "" {
		args = append(args, "--username", c.Username, "--password", c.Password)
	}
	if limit > 0 {
		args = append(args, "-l", fmt.Sprintf("%d", limit))
	}
	if revision != "" {
		args = append(args, "-r", revision)
	}
	for _, s := range search {
		args = append(args, "--search", s)
	}
	for _, s := range searchand {
		args = append(args, "--search-and", s)
	}
	if c.Local != "" {
		args = append(args, c.Local)
	} else if c.Remote != "" {
		args = append(args, c.Remote)
	}
	cmd := exec.Command("svn", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(output)
		idx := strings.Index(msg, "svn: E")
		if idx >= 0 {
			msg := strings.TrimSpace(msg[idx:])
			return nil, fmt.Errorf(msg)
		}
		return nil, fmt.Errorf(msg)
	}

	log := &Log{}
	err = xml.Unmarshal(output, log)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return log, nil
}

func (c *Client) Diff(change, revision string) (patch []byte, err error) {
	args := []string{"diff"}
	if c.Username != "" && c.Password != "" {
		args = append(args, "--username", c.Username, "--password", c.Password)
	}
	if revision != "" {
		args = append(args, "-r", revision)
	}
	if change != "" {
		args = append(args, "-c", change)
	}
	if c.Local != "" {
		args = append(args, c.Local)
	} else if c.Remote != "" {
		args = append(args, c.Remote)
	}
	cmd := exec.Command("svn", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(output)
		idx := strings.Index(msg, "svn: E")
		if idx >= 0 {
			msg := strings.TrimSpace(msg[idx:])
			return nil, fmt.Errorf(msg)
		}
		return nil, fmt.Errorf(msg)
	}
	return output, nil
}
