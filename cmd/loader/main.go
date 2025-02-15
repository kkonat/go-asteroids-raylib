package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/yeka/zip" // added for unzipping password-protected files
)

//go:embed splash.txt
var splash string

const program = "bangbang.exe"
const args = "-v"
const programPath = "./" + program

const githubRepo = "kkonat/go-asteroids-raylib"

type GitHubReleaseAssets struct {
	DownloadUrl string `json:"browser_download_url"`
}

type GitHubRelease struct {
	TagName       string `json:"name"`
	PublishedDate string `json:"published_at"`
	Assets        []GitHubReleaseAssets
}

func parseVersion(version string) (int, int, int) {
	var major, minor, build int
	// search from the name end for the 'v' character
	start := strings.LastIndex(version, "v")
	if start == -1 {
		fmt.Println("Error: version does not contain 'v' character")
		return -1, 0, 0
	}
	version = version[start:]
	n, err := fmt.Sscanf(version, "v%d.%d.%d", &major, &minor, &build)
	if n != 3 || err != nil {
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Error parsing version from string: %s\n", version)
		return -1, 0, 0
	}
	return major, minor, build
}

func checkLatestGithubRelease() (int, int, int, string) {
	// /releases/latest - does not work for some reason
	url := "https://api.github.com/repos/" + githubRepo + "/releases"
	header := map[string]string{
		"Accept": "application/vnd.github+json",
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		return 0, 0, 0, ""
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching latest release: %s\n", err)
		return 0, 0, 0, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received non-200 response code: %d\n", resp.StatusCode)
		return 0, 0, 0, ""
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		fmt.Printf("Error decoding JSON response: %s\n", err)
		return 0, 0, 0, ""
	}
	if len(releases) == 0 {
		fmt.Println("No releases found")
		return 0, 0, 0, ""
	}

	// sort releases by published date
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].PublishedDate > releases[j].PublishedDate
	})
	if len(releases[0].Assets) == 0 {
		fmt.Println("No assets found")
	}
	if len(releases[0].Assets) == 0 {
		fmt.Println("No assets found")
	}
	downloadUrl := releases[0].Assets[0].DownloadUrl
	var major, minor, build int
	n, err := fmt.Sscanf(releases[0].TagName, "v%d.%d.%d", &major, &minor, &build)
	if n != 3 || err != nil {
		fmt.Printf("Error parsing version from tag: %s\n", releases[0].TagName)
		if err != nil {
			fmt.Println(err)
		}
	}
	return major, minor, build, downloadUrl
}
func downloadFile(url string) (error, string) {
	// Split URL to get the filename
	filename := url[strings.LastIndex(url, "/")+1:]

	// Create the file
	outFile, err := os.Create(filename)
	if err != nil {
		return err, ""
	}

	// Connect to the server
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return err, ""
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status), ""
	}

	// Get the size of the file
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return fmt.Errorf("unable to get content length: %v", err), ""
	}

	// Create a progress reader
	progressReader := &ProgressReader{
		Reader: resp.Body,
		Total:  int64(size),
	}

	// Write the body to file
	_, err = io.Copy(outFile, progressReader)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err), ""
	}

	fmt.Println()
	outFile.Close()
	return nil, filename
}

// ProgressReader is a custom reader to track progress
type ProgressReader struct {
	Reader    io.Reader
	Total     int64
	ReadBytes int64
}

// Read reads data and updates the progress
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if n > 0 {
		pr.ReadBytes += int64(n)
		fmt.Printf("\r  %d%% complete", int(float64(pr.ReadBytes)/float64(pr.Total)*100))
	}
	return n, err
}

func main() {
	fmt.Println(splash)
	// execute ../bangbang.exe

	cmd := exec.Command(programPath, args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Can't execute the main program (%s %s):%s", programPath, args, err)
		return
	}
	lmajor, lminor, lbuild := parseVersion(string(output))
	if lmajor == -1 {
		fmt.Printf("Error parsing local version: %s\n", output)
		return
	}
	gmajor, gminor, gbuild, downloadUrl := checkLatestGithubRelease()
	fmt.Printf("v%d.%d.%d\n", lmajor, lminor, lbuild)
	var fileName = ""
	if lmajor < gmajor || lminor < gminor || lbuild < gbuild {
		fmt.Printf("New version available! (%d.%d.%d)\n", gmajor, gminor, gbuild)
		// download the latest version
		fmt.Println("  Downloading the latest version...")
		err, fileName = downloadFile(downloadUrl)
		if err != nil {
			fmt.Println("Error downloading the latest version:", err)
			return
		}

		// unzip the file
		archive, err := zip.OpenReader(fileName)
		if err != nil {
			fmt.Println("Error opening zip file:", err)
			return
		}
		var outFile *os.File
		extracted := false
		for _, file := range archive.File {
			if file.Name == program {
				zipFile, err := file.Open()
				if err != nil {
					fmt.Println("Error opening file in zip:", err)
					return
				}

				outFile, err = os.Create(program + ".new")
				if err != nil {
					fmt.Println("Error creating game.exe:", err)
					return
				}

				_, err = io.Copy(outFile, zipFile)
				if err != nil {
					fmt.Println("Error extracting game.exe:", err)
					return
				}

				extracted = true
				zipFile.Close()
				break
			}

		}
		if !extracted {
			fmt.Println("game.exe not found inside zip file.")
			return
		}

		// Explicitly close the archive to release the file handle.
		if err := archive.Close(); err != nil {
			fmt.Println("Error closing zip archive:", err)
			return
		}
		// Explicitly close the output file.
		if err := outFile.Close(); err != nil {
			fmt.Println("Error closing game.exe:", err)
			return
		}
		// Clean up: delete the downloaded zip file.
		if err := os.Remove(fileName); err != nil {
			fmt.Println("Error deleting zip file:", err)
			return

		}
		// move program + ".new" to program
		if err := os.Rename(program+".new", program); err != nil {
			fmt.Println("Error renaming game.exe:", err)
			return
		}
		fmt.Println("Update applied.")
	} else {
		fmt.Printf("You have the latest version")
	}

	cmd = exec.Command(programPath, "-from-loader")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error executing the main program (%s %s): %s\n", programPath, args, err)
		return
	}
}
