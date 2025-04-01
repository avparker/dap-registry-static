package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TBD54566975/dap-go/dap"
	"github.com/alecthomas/types/optional"
	"github.com/tbd54566975/web5-go/dids/didweb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide desired handle e.g. moegrammer")
		return
	}

	handle := os.Args[1]

	fmt.Println("Enter comma delimited money addresses:")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	input = strings.TrimSuffix(input, "\n")
	maddrs := strings.Split(input, ",")
	maddrServices := make([]didweb.CreateOption, len(maddrs))

	for i, maddr := range maddrs {
		maddrService := didweb.Service(strconv.Itoa(i), "MoneyAddress", maddr)
		maddrServices[i] = maddrService
	}

	domain := "dapme.xyz"
	url := domain + "/" + handle
	bearerDID, err := didweb.Create(url, maddrServices...)
	if err != nil {
		fmt.Println("Error creating DID:", err)
		return
	}

	didDocument, err := json.MarshalIndent(bearerDID.Document, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling DID Document:", err)
		return
	}

	docPath := filepath.Join("registry", handle, "did.json")
	writeFile(didDocument, docPath)

	r := dap.NewRegistration(handle, domain, bearerDID.URI)
	err = r.Sign(bearerDID)
	if err != nil {
		fmt.Println("Error signing registration:", err)
		return
	}

	resolutionResponse := dap.ResolutionResponse{
		DID:   bearerDID.URI,
		Proof: optional.Some(r),
	}

	jsonResolutionResponse, err := json.MarshalIndent(resolutionResponse, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling registration:", err)
		return
	}

	regPath := filepath.Join("registry", "daps", handle)
	writeFile(jsonResolutionResponse, regPath)

	fmt.Println("now open a PR with the auto-generated files")
}

// Write content to a file, creating directories if necessary
func writeFile(content []byte, filePath string) error {
	// Create directories if they don't exist
	dir := filepath.Dir(filePath)
	if !directoryExists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write bytes to the file
	_, err = file.Write(content)
	if err != nil {
		return err
	}

	fmt.Printf("wrote %s\n", filePath)

	return nil
}

// Check if a directory exists
func directoryExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
