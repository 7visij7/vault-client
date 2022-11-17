package vault

import (
	"os"
	"fmt"
	// "log"
	"net/http"
	"strings"
	"crypto/tls"
	"smib-vault-client/types"
	"smib-vault-client/pkg/aes"
	"smib-vault-client/pkg/errors"
	"smib-vault-client/pkg/base64"
	"github.com/hashicorp/vault/api"
)

var (
	ListPath = ListSecret{}
)

type ListSecret struct {
	Items []ListSecretParams
}

type ListSecretParams struct {
	Path	string
}

func (list *ListSecret) AddItem(item ListSecretParams) {
	list.Items = append(list.Items, item)
}

func CheckFlag(filename string) (bool){
	if filename == ""{
		return false
	}
	return true
}

func userpassLogin() (string) {
	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},},}
	client, err := api.NewClient(&api.Config{Address: types.VAULT_ADDR, HttpClient: httpClient})
	errors.CheckError(err)

	options := map[string]interface{}{
		"password": types.PASSWORD,
	}
	path := fmt.Sprintf("auth/ldap/login/%s", types.USERNAME)
	secret, err := client.Logical().Write(path, options)
	errors.CheckError(err)

	token := secret.Auth.ClientToken
	return token
}

func connectVault(vault_addr string, token string)  (*api.Client) {
	http_client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},},}

	config := &api.Config{
		Address: vault_addr,
		HttpClient: http_client,
	}
	
	client, err := api.NewClient(config)
	errors.CheckError(err)

	client.SetToken(token)
	return client
}

func GetVaultSecret(vault_path string) *api.Secret {
	client := connectVault(types.VAULT_ADDR, userpassLogin())
	secret, err := client.Logical().Read(vault_path)
	errors.CheckError(err)

	CheckVaultSecret(secret, vault_path)
	
	return secret
}

func CheckVaultSecret(secret *api.Secret, vault_path string) {
	if secret == nil {
		fmt.Printf("Error: Can not get access to %s. Please check path to backend secret Vault.\n" , vault_path )
		os.Exit(1)
	}
	if secret.Data["data"] == nil {
		fmt.Printf("No secret is stored on path %s\n", vault_path )
		os.Exit(1)
	}
}

func Read(path string, filename string, encrypt bool, base64Flag bool, separator string) {
	vault_path := fmt.Sprintf("%s/data/%s", types.PROJECT_NAME, path)
	secrets := GetVaultSecret(vault_path)
	CheckVaultSecret(secrets, vault_path)

	var result string
	data := secrets.Data["data"].(map[string]interface{})

	for k, v := range data {
		if encrypt{
			v = aes.Encrypt(v.(string))
		}
		if base64Flag{
			v = base64.Base64Encode(v.(string))
		}
		result += fmt.Sprintf("%s%s%s\n", k, separator, v)
	}

	if CheckFlag(filename) {
		err := os.WriteFile(filename, []byte(result), 0666)
		errors.CheckError(err)
	} else {
		fmt.Println(result)
	}
}

func Delete(path string) {
	client := connectVault(types.VAULT_ADDR, userpassLogin())
	List(path)

	for _, j := range ListPath.Items{
		vault_path := fmt.Sprintf("%s/metadata/%s", types.PROJECT_NAME, j.Path)
		_, err := client.Logical().Delete(vault_path)
		errors.CheckError(err)	
		fmt.Printf("Successfully deleted secret %s\n", j.Path)
	}
}

func Write(path string, secrets map[string]string) {
	client := connectVault(types.VAULT_ADDR, userpassLogin())
	inputData := map[string]interface{}{
        "data": map[string]interface{}{},
    }
	for k, v := range secrets{
		inputData["data"].(map[string]interface{})[k] = v
	}
	vault_path := fmt.Sprintf("%s/data/%s", types.PROJECT_NAME, path)
	_, err := client.Logical().Write(vault_path, inputData)
	errors.CheckError(err)
	fmt.Println("Successfully writed secrets to", vault_path)
}

func Copy(source string, destination string, soft bool) {
	List(source)
	client := connectVault(types.VAULT_ADDR, userpassLogin())
	for _, src := range ListPath.Items{
		secrets := GetVaultSecret(fmt.Sprintf("%s/data/%s", types.PROJECT_NAME, src.Path))
		dstPath := strings.Replace(src.Path, source, destination, 1)
		if (CheckPathEndPoint(dstPath, client) && soft) {
			continue
		}
		_, err := client.Logical().Write(fmt.Sprintf("%s/data/%s", types.PROJECT_NAME, dstPath), secrets.Data)
		errors.CheckError(err)
		fmt.Printf("Successfully copied secrets from %s to %s.\n", src.Path, dstPath)
	} 
}

func CheckPathEndPoint(path string, client *api.Client) bool {	
	keys, err := client.Logical().Read(fmt.Sprintf("%s/data/%s", types.PROJECT_NAME, path))
	errors.CheckError(err)
	if keys == nil {
		return false
	}
	return true
}

func GetList(path string, client *api.Client) {
	metadata := fmt.Sprintf("%s/metadata/%s", types.PROJECT_NAME, path)
	secrets, _ := client.Logical().List(metadata)
	if secrets == nil {
		return
	}

	data := secrets.Data["keys"].([]interface{})
	for _, v := range data {
		subPath := fmt.Sprintf("%s/%s", path, v)
		ListPath.AddItem(ListSecretParams{Path:subPath})
		if !(CheckPathEndPoint(subPath, client)) { 
			GetList(subPath, client) 
		}
	}
}

func List(path string) {
	client := connectVault(types.VAULT_ADDR, userpassLogin())

	if CheckPathEndPoint(path, client) {
		ListPath.AddItem(ListSecretParams{Path:path})
	}

	GetList(path, client)
	if (len(ListPath.Items) == 0) {
		fmt.Printf("No secret is stored on path %s.\nCheck the vault path.", path)
	}
}