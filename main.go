package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

type File struct {
	Name      string
	Extension string
}

type Files map[string]File

func main() {
	path, key := getArgs()
	files := getFiles(path)

	fmt.Printf("%d files found.\n-----\n", len(files))

	for _, file := range files {
		convertFile(file, path, key)
	}

	fmt.Println("Complete!")
}

func getArgs() (string, string) {
	if len(os.Args) < 3 {
		fmt.Println("Argument for path and activation bytes required, e.g. convert ./some-dir abcd1234")
		fmt.Println("Activation bytes can be retrieved from https://audible-converter.ml")
		os.Exit(1)
	}

	if strings.Contains(os.Args[1], ".aax") {
		fmt.Println("Please provide a directory path, file passed.")
		os.Exit(1)
	}

	if len(os.Args[2]) != 8 {
		fmt.Println("Activation bytes too long, value needs to be 4 bytes.")
		os.Exit(1)
	}

	return os.Args[1], os.Args[2]
}

func getFiles(path string) Files {
	files, err := locateFiles(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return files
}

func locateFiles(path string) (Files, error) {
	fileSystem := os.DirFS(path)
	files := make(map[string]File)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsRegular() || !strings.Contains(d.Name(), ".aax") {
			return nil
		}

		nameParts := strings.Split(d.Name(), ".")
		name := strings.Join(nameParts[:1], "")
		extension := strings.Join(nameParts[len(nameParts)-1:], "")

		files[path] = File{
			Name:      name,
			Extension: extension,
		}

		return nil
	})

	return files, nil
}

func convertFile(file File, path, key string) {
	fmt.Printf("Converting %s...\n", file.Name)

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-activation_bytes",
		key,
		"-i",
		fmt.Sprintf("%s/%s.%s", path, file.Name, file.Extension),
		"-codec",
		"copy",
		fmt.Sprintf("%s/%s.%s", path, file.Name, "m4b"),
	)

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("-")
}
