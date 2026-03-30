package digest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func generateDiscordTableOfContents(digestPath, version string) error {
	var builder = strings.Builder{}

	builder.WriteString(fmt.Sprintf("# Release %s\n", version))

	err := filepath.Walk(digestPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}
		if info.Name() == fileDiscord {
			return nil
		}
		data, err := os.ReadFile(path)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSuffix(scanner.Text(), version)
			if strings.HasPrefix(line, "# ") {
				header := strings.TrimPrefix(line, "# ")
				builder.WriteString("## " + header + "\n")
				continue
			}
			if strings.Contains(line, "Table of Contents") {
				continue
			}
			if strings.Contains(line, "Notes") {
				continue
			}
			if strings.HasPrefix(line, "## ") {
				header := strings.TrimPrefix(line, "## ")
				link := strings.ReplaceAll(strings.ToLower(header), " ", "-")
				builder.WriteString(
					fmt.Sprintf(
						" * [%s](%s#%s)\n",
						header,
						info.Name(),
						link,
					),
				)
				continue
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not parse documentation for discord toc: %v", err)
	}

	err = os.WriteFile(filepath.Join(digestPath, fileDiscord), []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("could not generate discord toc: %v", err)
	}

	return nil
}
