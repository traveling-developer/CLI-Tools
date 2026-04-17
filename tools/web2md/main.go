package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	url := flag.String("url", "", "URL of the web page to convert")
	flag.StringVar(url, "u", "", "shorthand for --url")
	output := flag.String("output", "", "write output to file instead of stdout")
	flag.StringVar(output, "o", "", "shorthand for --output")
	help := flag.Bool("help", false, "show this help message")
	flag.BoolVar(help, "h", false, "shorthand for --help")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: web2md --url <url> [options]

Convert a web page to Markdown.

Options:
  -u, --url <url>       URL of the web page to convert (required)
  -o, --output <file>   Write output to file (default: stdout)
  -h, --help            Show this help message

Examples:
  web2md --url https://example.com
  web2md --url https://example.com --output article.md
  web2md --url https://example.com -o article.md
`)
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *url == "" {
		fmt.Fprintln(os.Stderr, "error: --url is required")
		flag.Usage()
		os.Exit(1)
	}

	binaryPath, cleanup, err := extractDefuddle()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer cleanup()

	args := []string{"parse", "--markdown", *url}
	if *output != "" {
		args = append(args, "--output", *output)
	}

	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if ps := exitErr.ProcessState; ps != nil && !ps.Exited() {
				fmt.Fprintf(os.Stderr, "error: defuddle terminated abnormally: %s\n", ps.String())
				os.Exit(1)
			}
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func extractDefuddle() (string, func(), error) {
	tmp, err := os.CreateTemp("", "defuddle-*")
	if err != nil {
		return "", func() {}, err
	}

	if _, err := tmp.Write(defuddleBinary); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", func() {}, err
	}
	tmp.Close()

	if err := os.Chmod(tmp.Name(), 0755); err != nil {
		os.Remove(tmp.Name())
		return "", func() {}, err
	}

	if runtime.GOOS == "darwin" {
		_ = exec.Command("codesign", "--remove-signature", tmp.Name()).Run()
		if out, err := exec.Command("codesign", "--sign", "-", "--force", tmp.Name()).CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: codesign failed (continuing, bun may have pre-signed): %v: %s\n", err, out)
		}
	}

	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}
