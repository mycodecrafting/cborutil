package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/json-iterator/go"
)

// encoding/json does not support map[interface {}]interface {}
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	decodeCmd := flag.NewFlagSet("decode", flag.ExitOnError)
	decodeCbor := decodeCmd.String("cbor", "", "CBOR encoded data in hex or base64 format")
	decodePath := decodeCmd.String("path", "", "path")
	decodeIsBase64 := decodeCmd.Bool("base64", false, "Use Base64 encoding instead of hex")
	decodeAsJSON := decodeCmd.Bool("json", false, "Return result in JSON format?")
	decodeCmd.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Printf("  %s -cbor <payload> [ -path data.path -base64 -json ]\n\n", os.Args[0])
		fmt.Println("Options:")
		decodeCmd.PrintDefaults()
		os.Exit(1)
	}

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateCbor := updateCmd.String("cbor", "", "CBOR encoded data in hex or base64 format")
	updatePath := updateCmd.String("path", "", "path")
	updateData := updateCmd.String("data", "", "Input data")
	updateIsBase64 := updateCmd.Bool("base64", false, "Use Base64 encoding instead of hex")
	updateAsJSON := updateCmd.Bool("json", false, "Input data is in JSON format?")

	flag.Usage = func() {
		fmt.Println("CBOR Utility")
		fmt.Println("")
		fmt.Println("This utility can decode data in a CBOR payload located at a specific data path.")
		fmt.Println("This utility can also be used to update data in a CBOR payload at a specific data path.")
		fmt.Println("")
		fmt.Println("-----")
		fmt.Println("")
		fmt.Println("Decode usage:")
		fmt.Println("")
		fmt.Printf("  %s decode -cbor <payload> [ -path data.path -base64 -json ]\n\n", os.Args[0])
		fmt.Println("decode options:")
		decodeCmd.PrintDefaults()
		fmt.Println("")
		fmt.Println("-----")
		fmt.Println("")
		fmt.Println("Update usage:")
		fmt.Println("")
		fmt.Printf("  %s update -cbor <payload> -data <input> -path data.path [ -base64 -json ]\n\n", os.Args[0])
		fmt.Println("update options:")
		updateCmd.PrintDefaults()
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		flag.Usage()
	}

	switch os.Args[1] {
	case "decode":
		decodeCmd.Parse(os.Args[2:])
		if *decodeCbor == "" {
			decodeCmd.Usage()
			os.Exit(1)
		}
		data, err := DecodePath(*decodeCbor, *decodePath, *decodeIsBase64)
		if err != nil {
			panic(err)
		}
		if *decodeAsJSON {
			jsonB, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}
			data = string(jsonB)
		}
		fmt.Println(data)
	case "update":
		updateCmd.Parse(os.Args[2:])
		var data string
		var err error
		if *updateAsJSON {
			var v interface{}
			json.Unmarshal([]byte(*updateData), &v)
			data, err = UpdatePath(*updateCbor, *updatePath, v, *updateIsBase64)
			if err != nil {
				panic(err)
			}
		} else {
			data, err = UpdatePath(*updateCbor, *updatePath, *updateData, *updateIsBase64)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println(data)
	default:
		flag.Usage()
	}
}
