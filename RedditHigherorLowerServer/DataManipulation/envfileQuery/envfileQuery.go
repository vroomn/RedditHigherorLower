package envfilequery

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Possible keys (must be initalized with parse):
//
//	DISCORD_CLIENT_ID
//	DISCORD_CLIENT_SECRET
//	ROOT_DIR
var EnvKeys map[string]string = make(map[string]string)

func Parse() {
	// Erros with local path so need to generate absolute path
	baseDir, _ := os.Executable()
	rootPath, _ := filepath.Split(filepath.ToSlash(baseDir))

	filePtr, openErr := os.Open(rootPath + ".env")
	if openErr != nil {
		log.Fatalln("Failed to open envfile, aborting... \nError:", openErr)
	}

	// Need the size of the file
	fileInfo, statErr := filePtr.Stat()
	if statErr != nil {
		log.Fatalln("Envfile stat read failed, aborting... \nError:", openErr)
	}

	// Copy data into buffer
	buf := make([]byte, fileInfo.Size())
	bytesRead, readErr := filePtr.Read(buf)
	if readErr != io.EOF && bytesRead != int(fileInfo.Size()) {
		log.Fatalln("Read returned before EOF, aborting... \nError", readErr)
	}

	var key, value []byte
	for i := 0; i < len(buf); {
		// Get the key
		if buf[i] != '=' {
			key = append(key, buf[i])
			i++
		} else {
			// Get the value associated
			i++
			for i < len(buf) && buf[i] != byte(13) {
				value = append(value, buf[i])
				i++
			}
			EnvKeys[string(key)] = string(value)
			key = nil
			value = nil
			i += 2
		}
	}
}
