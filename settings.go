package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"reflect"
)

var (
	
	SettingsFile = "config.json"

	App = &AppSettings{
		Host:   "",
		Port:   "8000",
		APIKey: "none",
	}
)

type AppSettings struct {
	Host   string `json:"host" env:"HOST" flag:"a" help:"net interface address (127.0.0.1 for localhost only)"`
	Port   string `json:"port" env:"PORT" flag:"p" help:"port of main http server"`
	APIKey string `json:"apikey" env:"APIKEY" flag:"apikey" help:"key for API access"`
}

func (as *AppSettings) Init() {

	//load from json
	if f, err := os.Open(SettingsFile); err == nil {
		dec := json.NewDecoder(f)
		if err := dec.Decode(as); err != nil {
			log.Printf("%s: %s", SettingsFile, err)
		} else {
			log.Printf("%s: success reading", SettingsFile)
		}
		f.Close()
	} else {
		log.Printf("%s: %s", SettingsFile, err)
	}

	typ := reflect.TypeOf(as).Elem()
	for i := 0; i < typ.NumField(); i++ {

		f := typ.Field(i)
		ff := reflect.ValueOf(as).Elem().FieldByName(f.Name)

		//env
		if t, ok := f.Tag.Lookup("env"); ok {
			if se, ok := os.LookupEnv(t); ok {
				ff.Set(reflect.ValueOf(se))
			}
		}

		//flags
		help := ""
		if t, ok := f.Tag.Lookup("help"); ok {
			help = t
		}
		if t, ok := f.Tag.Lookup("flag"); ok {
			switch f.Type.Kind() {
			case reflect.String:
				flag.StringVar(ff.Addr().Interface().(*string), t, ff.String(), help)
			case reflect.Int:
				flag.IntVar(ff.Addr().Interface().(*int), t, int(ff.Int()), help)
			case reflect.Uint:
				flag.UintVar(ff.Addr().Interface().(*uint), t, uint(ff.Uint()), help)
			case reflect.Float64:
				flag.Float64Var(ff.Addr().Interface().(*float64), t, ff.Float(), help)
			case reflect.Bool:
				flag.BoolVar(ff.Addr().Interface().(*bool), t, ff.Bool(), help)
			default:
				panic("incorrect type of settings field")
			}
		}
	}

	newkey := flag.Bool("newkey", false, "generate new API key and store in "+SettingsFile)

	flag.Parse()

	if *newkey {

		buf := make([]byte, 16)
		io.ReadFull(rand.Reader, buf)
		App.APIKey = base64.RawStdEncoding.EncodeToString(buf)

		if f, err := os.Create(SettingsFile); err == nil {
			enc := json.NewEncoder(f)
			enc.SetIndent("","    ")
			if err := enc.Encode(App); err != nil {
				log.Printf("Creating %s: %s", SettingsFile, err)
			} else {
				log.Printf("Creating %s: success", SettingsFile)
			}
			f.Close()
		} else {
			log.Printf("Creating %s: %s", SettingsFile, err)
		}

		os.Exit(0)
	}

}
