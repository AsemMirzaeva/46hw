package main

import (
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type GoRelease struct {
    Version string `json:"version"`
    Files   []struct {
        Filename string `json:"filename"`
        Sha256   string `json:"sha256"`
    } `json:"files"`
}

const goURL = "https://go.dev/dl/go1.22.3.src.tar.gz"
const jsonURL = "https://go.dev/dl/?mode=json"
const fileName = "go1.22.3.src.tar.gz"

func downloadGoArchive(url, filePath string) error {
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("error downloading the file: %w", err)
    }
    defer resp.Body.Close()

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("error creating the file: %w", err)
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    if err != nil {
        return fmt.Errorf("error saving the file: %w", err)
    }
    return nil
}

func getExpectedHash(url, version, fileName string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("error fetching the JSON: %w", err)
    }
    defer resp.Body.Close()

    var releases []GoRelease
    err = json.NewDecoder(resp.Body).Decode(&releases)
    if err != nil {
        return "", fmt.Errorf("error decoding the JSON: %w", err)
    }

    for _, release := range releases {
        if release.Version == version {
            for _, file := range release.Files {
                if file.Filename == fileName {
                    return file.Sha256, nil
                }
            }
        }
    }
    return "", fmt.Errorf("hash not found in JSON")
}

func calculateFileHash(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("error opening the file: %w", err)
    }
    defer file.Close()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", fmt.Errorf("error calculating the hash: %w", err)
    }
    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func main() {
	
    if err := downloadGoArchive(goURL, fileName); err != nil {
        fmt.Println(err)
        return
    }

    expectedHash, err := getExpectedHash(jsonURL, "go1.22.3", fileName)
    if err != nil {
        fmt.Println(err)
        return
    }

    calculatedHash, err := calculateFileHash(fileName)
    if err != nil {
        fmt.Println(err)
        return
    }

    if calculatedHash == expectedHash {
        fmt.Println("true")
    } else {
        fmt.Println("false")
    }
}
