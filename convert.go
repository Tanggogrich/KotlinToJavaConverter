package main

func Convert(folder string) {
	/*
		Scan method takes the folder and scans each .kt files until the end.
		After it creates a lexical tree that will be converted into java code.
		Return a list of java files.
	*/
	var files, err = Scan(folder)
	if err != nil {
		panic(err)
	}
	err = Compile(files)
	/*
		Create the new directory for Java files and write them inside new directory.
	*/
	CreateJavaDir(files)
}
