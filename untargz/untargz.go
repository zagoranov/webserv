package untargz

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func ExtractTarGz(gzipStream io.Reader) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Printf("ExtractTarGz: NewReader failed to open file")
	} else {
		tarReader := tar.NewReader(uncompressedStream)

		for true {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("ExtractTarGz: Next() failed: %s", err.Error())
			} else {

				switch header.Typeflag {
				case tar.TypeDir:
					if err := os.Mkdir(header.Name, 0755); err != nil {
						log.Printf("ExtractTarGz: Mkdir() failed: %s", err.Error())
						break
					}
				case tar.TypeReg:
					outFile, err := os.Create(header.Name)
					if err != nil {
						log.Printf("ExtractTarGz: Create() failed: %s", err.Error())
						break
					}
					if _, err := io.Copy(outFile, tarReader); err != nil {
						log.Printf("ExtractTarGz: Copy() failed: %s", err.Error())
						break
					}
					outFile.Close()

				default:
					log.Printf(
						"ExtractTarGz: uknown type: %s in %s",
						header.Typeflag,
						header.Name)
				}
			}
		}
	}
}
