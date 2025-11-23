package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(text string) string {
	caser := cases.Title(language.Und)
	return caser.String(text)
}

func NormalizePath(path string) string {
	if strings.HasPrefix(path, "~") {
		path, _ = homedir.Expand(path)
	}
	return path
}

func BinaryExists(binaryname string) bool {
	_, err := exec.LookPath(binaryname)
	return err == nil
}

func AppendToContent(filename string, data []byte) (string, error) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	if _, err := f.Write(data); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return filename, nil
}

func FileExists(filename string) bool {
	filename = NormalizePath(filename)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func FolderExists(foldername string) bool {
	foldername = NormalizePath(foldername)
	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileWithoutExtension(fileNamePath string) string {
	fileName := filepath.Base(fileNamePath)
	return strings.TrimSuffix(fileName, path.Ext(fileName))
}

func CleanPath(raw string) string {
	raw = NormalizePath(raw)
	base := raw
	if FileExists(base) {
		base = filepath.Base(raw)
	}
	out := strings.ReplaceAll(base, "/", "_")
	out = strings.ReplaceAll(out, ":", "_")
	return out
}

func ReadFile(filename string) []string {
	var result []string
	if strings.HasPrefix(filename, "~") {
		filename = NormalizePath(filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		Logger.Error().Str("file", filename).Msg("Unable to open file")
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		if val == "" {
			continue
		}
		result = append(result, val)
	}

	if err := scanner.Err(); err != nil {
		return result
	}
	return result
}

func GetFileContent(filename string) []byte {
	var result []byte
	if strings.Contains(filename, "~") {
		filename, _ = homedir.Expand(filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return result
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return result
	}
	return b
}

func WriteFile(filename, data string) error {
	if strings.HasPrefix(filename, "~") {
		filename = NormalizePath(filename)
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data+"\n")
	if err != nil {
		return err
	}
	return file.Sync()
}

func ExecCmd(cmd, cwd, binPath string, show_output bool) (output string, err error) {
	Logger.Debug().Str("script", cmd).Msg("Execute command")
	envPath := fmt.Sprintf("%s:%s", os.Getenv("PATH"), binPath)
	os.Setenv("PATH", envPath)

	commands := []string{
		"bash",
		"-c",
		cmd,
	}
	realCmd := exec.Command(commands[0], commands[1:]...)

	// define current working directory on the workspace
	realCmd.Dir = cwd

	// retrieve output
	stdout, err := realCmd.StdoutPipe()
	if err != nil {
		return output, err
	}
	realCmd.Stderr = realCmd.Stdout

	if err = realCmd.Start(); err != nil {
		return output, err
	}

	out, _ := io.ReadAll(stdout)

	// Only print log if there is an output
	if len(out) > 1 && show_output {
		Logger.Info().Msg(string(out))
	}

	if err = realCmd.Wait(); err != nil {
		return output, err
	}

	return output, err
}
