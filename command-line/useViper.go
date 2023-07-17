package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func aliasNormalizeFunc(f *pflag.FlagSet, n string) pflag.NormalizedName {
	switch n {
	case "pass":
		fallthrough // body of the next case will be executed
	case "ps":
		n = "password"
	}
	return pflag.NormalizedName(n)
}

func main() {
	pflag.StringP("name", "n", "Mike", "Name parameter")
	pflag.StringP("password", "p", "hardToGuess", "Password")
	pflag.CommandLine.SetNormalizeFunc(aliasNormalizeFunc)

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	name := viper.GetString("name")
	password := viper.GetString("password")

	fmt.Println(name, password)

	viper.BindEnv("GOMAXPROCS")
	val := viper.Get("GOMAXPROCS")
	if val != nil {
		fmt.Println("GOMAXPROCS:", val)
	}

	viper.Set("GOMAXPROCS", 16)
	val = viper.Get("GOMAXPROCS")
	fmt.Println("GOMAXPROCS:", val)
}
