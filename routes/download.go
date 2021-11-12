package routes

import (
	"archive/zip"
	"aws-storage/middleware"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type FolderInfo struct {
	FolderPath string `json:"folderPath"`
}
type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func DownloadObjects(w http.ResponseWriter, r *http.Request) {
	middleware.LoadEnv()
	var output FolderInfo
	if err := json.NewDecoder(r.Body).Decode(&output); err != nil {
		log.Fatal(err)
	}
	zipFolderName := "job_" + output.FolderPath + "_files.zip"
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", string(zipFolderName)))
	folder := output.FolderPath
	dir, _ := os.Getwd()
	fmt.Println(dir)
	//cmd := exec.CommandContext(r.Context(), "/bin/bash", "script.sh", folder, dir)
	cmd := exec.Command("aws s3 cp s3://" + middleware.GetEnvWithKey("BUCKET_NAME") + "/" + string(folder) + " . --recursive") //USED FOR LOCAL TESTING
	cmd.Stderr = os.Stderr
	fmt.Println(cmd)
	/* out, err := cmd.Output()
	if err != nil {
		log.Printf("Command.Output: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	fmt.Println(out)
	files, err := ioutil.ReadDir(dir + "/" + folder)
	if err != nil {
		log.Fatal(err)
	}
	if err := ZipFiles(zipFolderName, folder, files); err != nil {
		log.Fatal(err)
	}
	http.ServeFile(w, r, zipFolderName) */
	fmt.Fprint(w, "Done downloading the job files!")

	/* sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess) */
}
func ZipFiles(filename string, foldername string, files []fs.FileInfo) error {
	newZipFile, err := os.Create(filename)
	dir, _ := os.Getwd()
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	if err != nil {
		return err
	}
	if err = zipSource(zipWriter, dir+"/"+foldername, filename); err != nil {
		return err
	}
	return nil
}

func zipSource(w *zip.Writer, source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		//fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := w.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
