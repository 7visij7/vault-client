# Vault-client
> Vault-client application help to work with Vault server. The main feature of application it can read secret by anchors. 
>
 > Anchor looks like: 
+ DATABASE_PASSWORD: VAULT:smbi/smbi-service-test#DATABASE_PASSWORD
> App developed for services which load settings with start and can't reread it after in running process. Therefore it is the way to not use vault sidecar injector. Vault-client allow to read, write, copy(recursively), delete(recursively) on Vault server. 
___
## Build
> To build application execute next command:
```Bash
go build -o vault-client
```
---
## Required variables
+ VAULT_ADDR - Vault server url
+ USERNAME - Vault user
+ PASSWORD - Password of vault user
+ ENCRYPT_KEY - secret key AES. 
+ PROJECT_NAME - KV store in Vault
___
## How it works:

> Write secret:

To use "write" command you should prepare file with secret like this example.yaml:
```
secret1: "foo"
secret2: "bar"
secret3: "secret"
```

Next parametr which you shouls specify is vault "path" where stored secrets. Exp: "dev/vault-client". 
With flag "replace" application changes value to vault path where secrets stored.
Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be:

+ "{VAULT_PROJECT_NAME}/{path}".

To the run command: 
```Bash
$ vault-client write --fileaname example.yaml --path dev/vault-client --replace
```
After specify flag "replace" in file example.yaml values of secret change to anchors:
```
secret1: VAULT:dev/vault-client#secret1
secret2: VAULT:dev/vault-client#secret2
secret3: VAULT:dev/vault-client#secret3
```

> Read secret:
  
Read command allow to get list of all secrets located on vault "path".

You could store received secrets to file. For this, please use flag "filename".

For the security of viewing information, there is opportunity to convert secret to encrypted by using AES, as a encryption key use global variables ENCRYPT_KEY. For this, please use flag "encrypt".
Or with flag "base64" data will be encode to base64.

Use only one flag: "encrypt" or "base64".

Flag separator let use different symbols which will be displaed between keys and values. By default separator is ': '"

Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 

"{VAULT_PROJECT_NAME}/{path}".

To the run command:
```Bash
$ vault-client read --path dev/vault-client --filename SuperPuper.secret --base64
```


> Read secrets:

Copy command allow to dublicate secrets located in "source" to other KV which you have to define as "destination".

If flag "soft" setted, and secret destination has already exist, coping does not run.

Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 

    "{VAULT_PROJECT_NAME}/{path}".

To the run command: 
```
$ vault-client copy  --source smib/vault-client --destination smbu/vault-client
```

> Delete secrets:
> 
Delete command allow to remove all secrets located on vault "path".

Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 

	"{VAULT_PROJECT_NAME}/{path}".

To the run command: 
```
$ vault-client delete --path smib/vault-client`
```

> Read secrets from Vault by anchors:

Discovery command allow to find all yaml files in directoy "helm" and put secrets to Vault.

To use "discovery" you should specify flags "ris", "enviroment" and "servicename".

Flag "enviroment" must be one of [dev, stage, production ... etc] values depends on in what catalog seek secrets.

With flag "replace" application changes value to vault path where secrets stored.

Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 

	"{VAULT_PROJECT_NAME}/{ris}/{servicename}#{key}".

To the run command:
```
	$ vault-client discovery --ris smib --servicename testing-new  --replace --enviroment dso
```

> Write secrets to Vault from files values.yaml
Helm command allow to find all yaml files in directoy, and replace values from Vault where you set anchor like this:

Variale_One: VAULT:kafka/kafka-devops#KAFKA_PASSWORD

Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server.

By default vault path will be formed:
+ "{VAULT_PROJECT_NAME}/{ris}/{servicename}#{key}",
 	+   ris -name of information system, 
 	+   servicename - name of +service, 
 	+   key - name of variable.
 	+   
Also you can specify directory where seek yaml files, for this use flag "path".
By default value of varaible will be encode to Base64, if you want to get value as is - use flag "raw". 

To the run command:
```
$ vault-client helm --path ./secret`
```