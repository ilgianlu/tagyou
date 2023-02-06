package conf

// valid backend values
const BACKEND_SQLITE = "sqlite"
const BACKEND_BADGER = "badger"

var AVAILABLE_BACKENDS = []string{
	BACKEND_SQLITE,
	BACKEND_BADGER,
}

// whick backend db to use, defaults to sqlite
var BACKEND_PERSISTENCE = BACKEND_SQLITE

// full or relative dir where db files will reside
var DB_PATH = "./"

// db name for sqlite single file, ignored on badger
var DB_NAME = "sqlite.db3"
