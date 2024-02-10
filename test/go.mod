module github.com/psyark/estat/test

go 1.21

replace github.com/psyark/estat => ../

replace github.com/psyark/estat/test => ./

require (
	github.com/joho/godotenv v1.5.1
	github.com/psyark/estat v0.0.0-00010101000000-000000000000
	github.com/wI2L/jsondiff v0.5.0
)

require (
	github.com/tidwall/gjson v1.17.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)
