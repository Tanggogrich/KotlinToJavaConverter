package structures

// DataFile describes the simple structure of found .kt-file,
// that consists the filename and raw code as a content
type DataFile struct {
	Name    string
	Content []byte
}

type DataFiles struct {
	Files []DataFile
}
