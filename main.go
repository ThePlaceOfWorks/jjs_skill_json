package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/klauspost/compress/zstd"
)

const CREDITS = "// agi // 2024"

func main() {
	if len(os.Args) < 2 {
		menuMode()
	} else {
		cliMode()
	}
}

func menuMode() {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Println()
		fmt.Println("ZSTD Tool " + CREDITS)
		fmt.Println("--------------------")
		fmt.Println("1) Decode (b64 -> json)")
		fmt.Println("2) Encode (json -> b64)")
		fmt.Println("Q) Cikis")
		fmt.Println()
		fmt.Print("Secim: ")
		
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))
		
		switch choice {
		case "1":
			doDecode(reader)
		case "2":
			doEncode(reader)
		case "q":
			fmt.Println("Gorusuruz.")
			return
		default:
			fmt.Println("Gecersiz secim")
		}
	}
}

func doDecode(reader *bufio.Reader) {
	fmt.Print("Girdi (.b64): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	fmt.Print("Cikti [output.json]: ")
	output, _ := reader.ReadString('\n')
	output = strings.TrimSpace(output)
	if output == "" {
		output = "output.json"
	}
	
	data, _ := os.ReadFile(input)
	zstdData, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
	
	zreader, _ := zstd.NewReader(nil)
	jsonBytes, _ := zreader.DecodeAll(zstdData, nil)
	zreader.Close()
	
	var jsonData interface{}
	json.Unmarshal(jsonBytes, &jsonData)
	
	outFile, _ := os.Create(output)
	defer outFile.Close()
	
	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	encoder.Encode(jsonData)
	
	fmt.Printf("Tamam: %s -> %s\n", input, output)
	fmt.Print("[Enter]")
	reader.ReadString('\n')
}

func doEncode(reader *bufio.Reader) {
	fmt.Print("Girdi (.json): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	fmt.Print("Cikti [output.b64]: ")
	output, _ := reader.ReadString('\n')
	output = strings.TrimSpace(output)
	if output == "" {
		output = "output.b64"
	}
	
	jsonBytes, _ := os.ReadFile(input)
	var jsonData interface{}
	json.Unmarshal(jsonBytes, &jsonData)
	
	compactJSON, _ := json.Marshal(jsonData)
	
	enc, _ := zstd.NewWriter(nil)
	zstdData := enc.EncodeAll(compactJSON, nil)
	enc.Close()
	
	b64 := base64.StdEncoding.EncodeToString(zstdData)
	os.WriteFile(output, []byte(b64), 0644)
	
	fmt.Printf("Tamam: %s -> %s\n", input, output)
	fmt.Print("[Enter]")
	reader.ReadString('\n')
}

func cliMode() {
	mode := strings.ToLower(os.Args[1])
	input := os.Args[2]
	output := ""
	if len(os.Args) > 3 {
		output = os.Args[3]
	}
	
	if mode == "decode" || mode == "d" {
		if output == "" { output = "output.json" }
		
		data, _ := os.ReadFile(input)
		zstdData, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
		
		zreader, _ := zstd.NewReader(nil)
		jsonBytes, _ := zreader.DecodeAll(zstdData, nil)
		zreader.Close()
		
		var jsonData interface{}
		json.Unmarshal(jsonBytes, &jsonData)
		
		outFile, _ := os.Create(output)
		defer outFile.Close()
		
		encoder := json.NewEncoder(outFile)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		encoder.Encode(jsonData)
		
		fmt.Printf("Tamam: %s -> %s\n", input, output)
		
	} else if mode == "encode" || mode == "e" {
		if output == "" { output = "output.b64" }
		
		jsonBytes, _ := os.ReadFile(input)
		var jsonData interface{}
		json.Unmarshal(jsonBytes, &jsonData)
		
		compactJSON, _ := json.Marshal(jsonData)
		
		enc, _ := zstd.NewWriter(nil)
		zstdData := enc.EncodeAll(compactJSON, nil)
		enc.Close()
		
		b64 := base64.StdEncoding.EncodeToString(zstdData)
		os.WriteFile(output, []byte(b64), 0644)
		
		fmt.Printf("Tamam: %s -> %s\n", input, output)
	}
}
