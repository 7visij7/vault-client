package helm

import (
	"os"
	"fmt"
	"bufio"
	"regexp"
	"strings"
	"io/ioutil"
	"smib-vault-client/types"
	"smib-vault-client/pkg/vault"
	"smib-vault-client/pkg/common"
	"smib-vault-client/pkg/base64"
	"smib-vault-client/pkg/errors"
)

func CheckVariable(name interface{}) bool {
	if name == "" {
		return false
	}
	return true
}

func CheckFiles(filePath string, raw bool) {
	file, err := os.Open(filePath)
	errors.CheckError(err)
	defer file.Close()
 
	scanner := bufio.NewScanner(file)
	var lines []string
	cache := make(map[string]interface{})
	for scanner.Scan() {
		text := scanner.Text()

		vaultRegexp := regexp.MustCompile("(?P<variable>[-_a-zA-Z0-9 ]{1,}:)(?P<pointer> VAULT:)(?P<path>[-a-zA-Z0-9/]{1,})(#)(?P<key>[-_a-zA-Z0-9]{1,})")
		result := vaultRegexp.FindAllStringSubmatch(text, 1)
		
		if len(result) > 0 {
			var vaultPath string
			var data  map[string]interface {}
			variable := result[0][1]
			path := result[0][3]
			key := result[0][5]
			vaultPath = fmt.Sprintf("%s/data/%s", types.PROJECT_NAME, path)
			
			if cache[path] != nil {
				data = cache[path].(map[string]interface{})
			} else {
				secrets := vault.GetVaultSecret(vaultPath)
				data = secrets.Data["data"].(map[string]interface{})	
				cache[path] = secrets.Data["data"].(map[string]interface{})		
			}
			var value string
			if  data[key] != nil {
				if raw {
					value =  fmt.Sprintf("%v", data[key])
				} else {
					value =  base64.Base64Encode(fmt.Sprintf("%v", data[key]))
				}
				text = fmt.Sprintf("%s %s", variable, value)
				fmt.Println("Successfully set value for",  key, "from", vaultPath)
			} else {
				fmt.Println("Error: Cant find value for key:", key, "in backend:", vaultPath )
			}
		}
	    lines = append(lines, text)
	}	
	errors.CheckError(scanner.Err())
 
	err = ioutil.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	errors.CheckError(err)
 }

func Run(path string, raw bool) {
	files := common.FindYamlFiles(path)
	for _, file := range files {
		CheckFiles(file, raw)
	}
}
