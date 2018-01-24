# easysettings
Simple getting an app settings from json (first), environment (second) and flags (third) into your struct with tags.
Imports only standard libraries.
Use very simple format of struct with tags for all variants of values representation:
```
type AppSettings struct {
	Host   string `json:"host" env:"HOST" flag:"a" help:"net interface address (127.0.0.1 for localhost only)"`
	Port   string `json:"port" env:"PORT" flag:"p" help:"port of main http server"`
	APIKey string `json:"apikey" env:"APIKEY" flag:"apikey" help:"key for API access"`
}
```

Sample usage:

```
go get github.com/covrom/easysettings
go build .
./easysettings.exe -h
./easysettings.exe -newkey
HOST=127.0.0.1 ./easysettings.exe -h
```
