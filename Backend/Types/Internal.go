package Types

type Config struct {

	GinMode string

	Server struct {

		Port int 

	}

	Database struct {

		URI string 
		Name string 

	}

	Versions struct {

		Server string
		Frontend string

	}

}

