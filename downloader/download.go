package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

func Download(url string, dest string) errorq

	req, _ := http.NewRequest("GET", url, nil)
    resp, _ := http.DefaultClient.Do(req)
    defer resp.Body.Close()
	tempPath := 
    f, _ := os.OpenFile(temp_path, os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()

    buf := make([]byte, 32*1024)
    var downloaded int64
    for {
        n, err := resp.Body.Read(buf)
        if err != nil {
            if err == io.EOF {
                break
            }
            log.Fatalf("Error while downloading: %v", err)
        }
        if n > 0 {
            f.Write(buf[:n])
            downloaded += int64(n)
            fmt.Printf("\rDownloading... %.2f%%", float64(downloaded)/float64(resp.ContentLength)*100)
        }
    }
    os.Rename(temp_path, "wordpress.zip")