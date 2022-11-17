package discovery

import (
	"os"
	"fmt"
	"sort"
	"bytes"
	"strings"
	"reflect"
	"io/ioutil"
	"gopkg.in/yaml.v3"
	"smib-vault-client/types"
	"smib-vault-client/pkg/vault"
	"smib-vault-client/pkg/common"
	"smib-vault-client/pkg/errors"
)

type OrderedMap struct {
	Order []string
	Map   map[interface{}]interface{}
}

var lines []string

func (om *OrderedMap) UnmarshalYAML(b []byte) error {
	yaml.Unmarshal(b, &om.Map)
    index := make(map[string]int)
    for key := range om.Map {
        k := fmt.Sprintf("%v", key)
		om.Order = append(om.Order, fmt.Sprintf("%v", k) )

		index[k] = bytes.Index(b, []byte(k))
	}
    sort.Slice(om.Order, func(i, j int) bool { return index[om.Order[i]] < index[om.Order[j]] })
    return nil
}

func (om OrderedMap) MarshalYAML() {    
    for _, key := range om.Order {
        Unpack(om.Map[key], key, "") 
    }
}

func GetSecretFromFile(path string) types.Secrets {
	byteValue := common.ReadFile(path)

	var data types.Secrets
	err := yaml.Unmarshal(byteValue, &data)
	errors.CheckError(err)
	return data
}

func Unpack(m interface{}, key string, indent string) {
    indent = fmt.Sprintf("%s", indent)
    lines = append(lines, fmt.Sprintf("%s%s: ", indent, key))
    switch m.(type) {
	case nil: 
        lines = append(lines, fmt.Sprintf("\n"))
    case bool, string, int:
        lines = append(lines, fmt.Sprintf("%v\n", m))
    default:
        lines = append(lines, fmt.Sprintf("\n"))
        indent = fmt.Sprintf("%s  ", indent)
        for index, value := range m.(map[string]interface{}) { 
            varType := reflect.ValueOf(value).Kind()
            if varType == reflect.Map {
                Unpack(value, index, indent)
                continue
            }
            lines = append(lines, fmt.Sprintf("%s%s: %v\n", indent, index, value))       
        }
    }
}

func CheckAnchorInFiles(ris string, servicename string, replace bool, environment string, file string) {
	secretsWithAchor := make(map[string]interface{})
	var yamlMap OrderedMap
	lines = nil
	ymlContent, err := ioutil.ReadFile(file)
	errors.CheckError(err)
	
	yaml.Unmarshal(ymlContent, &yamlMap)
	yamlMap.UnmarshalYAML(ymlContent)

	if yamlMap.Map["secret"] == nil { return }
	secrets := yamlMap.Map["secret"].(map[string]interface{})

	if len(secrets) > 0 {
		path := fmt.Sprintf("%s/%s", ris, servicename)
		for key, _ := range secrets {
			common.CheckVariable(fmt.Sprintf("%v", secrets[key]), key, file)
			secretsWithAchor[key] = fmt.Sprintf("VAULT:%s#%s", path, key)
		}
		yamlMap.Map["secret"] = secretsWithAchor
		secretToVault := common.ConvertInterfaceToString(secrets)
		fmt.Println("********************************")
		vault.Write(path, secretToVault)

	yamlMap.MarshalYAML()
	}
	if replace {
		err = ioutil.WriteFile(file,  []byte(strings.Join(lines, "")), 0755)
		errors.CheckError(err)
	}
}

func GetEnviroment(enviroment string) []string {
	switch enviroment {
	case "dso":
		return []string{"dso"}
	case "st":
		return []string{"st"}
	default:
		fmt.Println("Enviroment must be one of [dso, st] values depends on where application running")
		os.Exit(1)
	} 
	return []string{}
}

func Run(ris string, servicename string, replace bool, env string) {
	envAllowed := GetEnviroment(env)
	data := common.FindFilesByName("./helm", "values.yaml")
	for environment, file := range data {
		if common.CheckSliceContainsElement(envAllowed, environment){
			CheckAnchorInFiles(ris, servicename, replace, environment, file)
		}
	}
}
