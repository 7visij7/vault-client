package common

import (
	"os"
	"io"
	"fmt"
	// "bufio"
	"regexp"
	"strings"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v3"
	// "smib-vault-client/types"
	"smib-vault-client/pkg/base64"
	"smib-vault-client/pkg/errors"
)

func Unmarshal(yamlFile []byte, secrets  map[string]string) (map[string]string, error) {
    err := yaml.Unmarshal(yamlFile, &secrets)
	if err != nil {
        return secrets, err
    }
    return secrets, nil
}

func ReadFile(filename string) (byteValue []byte) {
	yamlFile, err := os.Open(filename)
    errors.CheckError(err)
	defer yamlFile.Close()
	
	byteValue, _ = ioutil.ReadAll(yamlFile)
	return byteValue
}

func GetSecrets(filename string) map[string]string{
    yamlFile := ReadFile(filename)

	var secrets map[string]string
    secrets, err := Unmarshal(yamlFile, secrets)
    errors.CheckError(err)

	for key, val := range secrets {
		CheckVariable(fmt.Sprintf("%v", secrets[key]), key, filename)
		secrets[key] = base64.CheckBase64(val)
	}
	
	return secrets
}

func Marshal(secret map[string]string, filename string) {
	buf, err := yaml.Marshal(secret)
	errors.CheckError(err)

	err = ioutil.WriteFile(filename, buf, 0644)
	errors.CheckError(err)
}


func ConvertInterfaceToString(secretIn map[string]interface{}) map[string]string {
	secretOut := make(map[string]string)
	for key, value := range secretIn {
		checkedValue := base64.CheckBase64(fmt.Sprintf("%v", value))
		secretOut[key] = checkedValue
	}
	return secretOut
} 

func AddAnchorToFile(path string, filename string, secrets map[string]string){
	for key, _ := range secrets {
		secrets[key] = fmt.Sprintf("VAULT:%s#%s", path, key)
	}
	Marshal(secrets ,filename)
}

func FindYamlFiles(dir_path string) []string {
	files := []string{}
	filepath.Walk(dir_path, func(path string, f os.FileInfo, err error) error {
		f, err = os.Stat(path)
		errors.CheckError(err)
		f_mode := f.Mode()
		if f_mode.IsRegular() {
			if filepath.Ext(f.Name()) == ".yaml" {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}

func FindFilesByName(dir_path string, name string) map[string]string {
	f := FindYamlFiles(dir_path)
	files := make(map[string]string)
	for _, path := range f {	
		f, err := os.Stat(path)
		errors.CheckError(err)
		if f.Name() == name {
			catalog := filepath.Base(filepath.Dir(path))
			files[catalog] = path
		}
	}
	return files
}

func CopyFile(src string, dst string) {
	source, err := os.Open(src)
	errors.CheckError(err)
	defer source.Close()

	destination, err := os.Create(dst)
	errors.CheckError(err)
	defer destination.Close()

	_, err = io.Copy(destination, source)
	errors.CheckError(err)
}

func CreateTmpFile(filename string) string{
	filenameTmp := fmt.Sprintf("/tmp/%s.tmp", filename)
	// fmt.Println(filenameTmp)
	CopyFile(filename, filenameTmp)
	return filenameTmp
}

func DeleteFile(filename string) {
	err := os.Remove(filename)
    errors.CheckError(err)
}

func CleanTmpFiles(filename string, filenameTmp string) { 
	CopyFile(filenameTmp, filename)
	DeleteFile(filenameTmp)
}

func BranchRename(inputString string) string{
	m := regexp.MustCompile("[-_./*]")
	Str := "-"
	result := m.ReplaceAllString(inputString, Str)
	return result
}

// func GetBranch() string{
// 	if types.CI_BRANCH == "" {
// 		file, err := os.Open(".git/HEAD")
// 		errors.CheckError(err)
// 		defer file.Close()
		 
// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() {
// 			text := scanner.Text()
	
// 			vaultRegexp := regexp.MustCompile("(?P<ref>ref: refs/heads/)(?P<branch>[-_a-zA-Z0-9/+:.]{1,})")
// 			result := vaultRegexp.FindAllStringSubmatch(text, 1)
			
// 			if len(result[0][2]) > 0 {
// 				return BranchRename(result[0][2])
// 			}
// 		}
// 	fmt.Println("Can not get information about branch. Check evrioment variables CI_BRANCH, or file .git/HEAD")
// 	os.Exit(1)
// 	}	

// 	return BranchRename(types.CI_BRANCH)
// }

func CheckSliceContainsElement(elems []string, comparable string) bool {
    for _, value := range elems {
        if value == comparable {
            return true
        }
    }
    return false
}

func CheckVariable(value string, key string, filename string) {
	if strings.Contains(value, "VAULT") {
		fmt.Printf("Error. Check value for variable %s in  %s \n", key, filename)
		os.Exit(1)
	}
}