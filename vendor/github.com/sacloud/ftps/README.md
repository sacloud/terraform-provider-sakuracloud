# FTPS Implementation for Go

## Information
This implementation does not implement the full FTP/FTPS specification. Only a small subset.

I have not done a security review of the code, yet. Therefore no guarantee is given. It would be nice if somebody could do a security review and report back if the implementation is vulnerable.

## Installation
	go get github.com/webguerilla/ftps

## Usage
	ftps := new(FTPS)

	ftps.TLSConfig.InsecureSkipVerify = true // often necessary in shared hosting environments
	ftps.Debug = true

	err := ftps.Connect("localhost", 21)
	if err != nil {
		panic(err)
	}

	err = ftps.Login("username", "password")
	if err != nil {
		panic(err)
	}

	directory, err := ftps.PrintWorkingDirectory()
	if err != nil {
		panic(err)
	}
	log.Printf("Current working directory: %s", directory)

	err = ftps.Quit()
	if err != nil {
		panic(err)
	}

## Credits
This work was inspired by the work of jlaffaye (https://github.com/jlaffaye/ftp) and smallfish (https://github.com/smallfish/ftp).
