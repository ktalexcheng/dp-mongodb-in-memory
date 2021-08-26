package download

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/spf13/afero"
)

var afs = afero.Afero{Fs: afero.NewOsFs()}

func init() {
	log.Namespace = "dp-mongodb-in-memory"
}

// GetMongoDB returns the path to the mongod binary for the given config
// It will download one if not found in the cache
func GetMongoDB(cfg Config) (string, error) {
	// Check the cache
	existsInCache, existsErr := afs.Exists(cfg.cachePath)
	if existsErr != nil {
		log.Error(context.Background(), "error checking cache", existsErr)
		return "", existsErr
	}
	if existsInCache {
		log.Info(context.Background(), "File found in cache", log.Data{"filename": cfg.cachePath})
		return cfg.cachePath, nil
	} else {
		return downloadMongoDB(cfg)
	}
}

// downloadMongoDB will download a mongodb tarball and
// store the mongod exec file in the cache path.
// It returns the path to the saved file
func downloadMongoDB(cfg Config) (string, error) {

	downloadStartTime := time.Now()

	downloadedFile, downloadErr := downloadFile(cfg.mongoUrl)
	if downloadErr != nil {
		log.Error(context.Background(), "error downloading file", downloadErr, log.Data{"url": cfg.mongoUrl})
		return "", downloadErr
	}

	defer func() {
		_ = downloadedFile.Close()
		_ = afs.Remove(downloadedFile.Name())
	}()

	validErr := verify(downloadedFile.Name(), cfg)
	if validErr != nil {
		log.Error(context.Background(), "error verifying integrity of MongoDB package", validErr, log.Data{"url": cfg.mongoUrl})
		return "", validErr
	}

	mongodTmpFile, mongoTmpErr := extractMongoBin(downloadedFile)
	if mongoTmpErr != nil {
		return "", mongoTmpErr
	}

	mkdirErr := afs.MkdirAll(path.Dir(cfg.cachePath), 0755)
	if mkdirErr != nil {
		log.Error(context.Background(), "error creating cache directory", mkdirErr, log.Data{"dir": path.Dir(cfg.cachePath)})
		return "", mkdirErr
	}

	renameErr := afs.Rename(mongodTmpFile, cfg.cachePath)
	if renameErr != nil {
		log.Error(context.Background(), "error copying mongod binary", renameErr, log.Data{"filename-from": mongodTmpFile, "filename-to": cfg.cachePath})
		return "", renameErr
	}

	log.Info(context.Background(), "mongod downloaded and stored in cache", log.Data{"filename": cfg.cachePath, "ellapsed": time.Since(downloadStartTime).String()})

	return cfg.cachePath, nil
}

// downloadFile downloads the file from the given url and stores it in a temporary file.
// It returns the path to the temporary file where it has been downloaded
func downloadFile(urlStr string) (afero.File, error) {
	log.Info(context.Background(), "Downloading file", log.Data{"url": urlStr})

	resp, httpGetErr := http.Get(urlStr)
	if httpGetErr != nil {
		return nil, httpGetErr
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	tgzTempFile, tmpFileErr := afs.TempFile("", "")
	if tmpFileErr != nil {
		return nil, tmpFileErr
	}

	_, copyErr := io.Copy(tgzTempFile, resp.Body)
	if copyErr != nil {
		_ = tgzTempFile.Close()
		_ = afs.Remove(tgzTempFile.Name())
		return nil, copyErr
	}

	log.Info(context.Background(), "Downloaded to temp file", log.Data{"file": tgzTempFile.Name(), "url": urlStr})
	return tgzTempFile, nil
}

// extractMongoBin extracts the mongod executable file
// from the given tarball to a temporary file.
// It returns the path to the extracted file
func extractMongoBin(tgzTempFile afero.File) (string, error) {
	_, seekErr := tgzTempFile.Seek(0, 0)
	if seekErr != nil {
		log.Error(context.Background(), "error seeking back to start of file", seekErr)
		return "", seekErr
	}

	gzReader, gzErr := gzip.NewReader(tgzTempFile)
	if gzErr != nil {
		log.Error(context.Background(), "error intializing gzip reader", gzErr, log.Data{"file": tgzTempFile.Name()})
		return "", gzErr
	}

	tarReader := tar.NewReader(gzReader)

	for {
		nextFile, tarErr := tarReader.Next()
		if tarErr == io.EOF {
			return "", fmt.Errorf("did not find a mongod binary in the tar file")
		}
		if tarErr != nil {
			log.Error(context.Background(), "error reading from tar file", tarErr, log.Data{"file": tgzTempFile.Name()})
			return "", tarErr
		}

		if strings.HasSuffix(nextFile.Name, "bin/mongod") {
			break
		}
	}

	// Extract to a temp file first, then copy to the destination, so we get
	// atomic behavior if there's multiple parallel downloaders
	mongodTmpFile, tmpFileErr := afs.TempFile("", "")
	if tmpFileErr != nil {
		log.Error(context.Background(), "error creating temp file for mongod", tmpFileErr)
		return "", tmpFileErr
	}
	defer func() {
		_ = mongodTmpFile.Close()
	}()

	_, writeErr := io.Copy(mongodTmpFile, tarReader)
	if writeErr != nil {
		log.Error(context.Background(), "error writing mongod binary", writeErr, log.Data{"filename": mongodTmpFile.Name()})
		return "", writeErr
	}

	_ = mongodTmpFile.Close()

	chmodErr := afs.Chmod(mongodTmpFile.Name(), 0755)
	if chmodErr != nil {
		log.Error(context.Background(), "error chmod-ing mongod binary", chmodErr, log.Data{"filename": mongodTmpFile.Name()})
		return "", chmodErr
	}
	return mongodTmpFile.Name(), nil
}

func verify(mongoFile string, cfg Config) error {
	//TODO implement validation against checksum/gpg key
	return nil
}
