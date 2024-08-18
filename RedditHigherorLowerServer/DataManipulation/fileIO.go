package datamanipulation

import (
	"RedditHigherorLowerServer/DataManipulation/envfile"
	"log"
	"strings"
)

const DATA_ENTRIES_MAX = -1

func WriteData() error {
	dataDir := strings.Trim(envfile.EnvKeys["ROOT_DIR"], "\"") + "data\\"

	log.Println(dataDir)

	return nil
}
