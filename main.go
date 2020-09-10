package main

import (
	"errors"
	"reflect"
	"strings"
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/vault/api"
	"github.com/namsral/flag"
)

type LoginRequest struct {
	Role   string `json:"role_id"`
	Secret string `json:"secret_id"`
}

type LoginAuthResponse struct {
	ClientToken string `json:"client_token"`
}

type LoginResponse struct {
	Auth LoginAuthResponse `json:"auth"`
}

type Config struct {
	ApproleId string
	ApproleSecret string
	HostSignerPath string
	VaultAddr string
	PublicKeyPath string
	CertPath string
}

func approleLogin(vaultAddr string, roleId string, secretId string) (string, error) {
	credentials := &LoginRequest{
		Role: roleId,
		Secret: secretId}

	var loginData []byte
	loginData, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	res, err := http.Post(
		fmt.Sprintf("%s/v1/auth/approle/login", vaultAddr),
		"application/json",
		strings.NewReader(string(loginData)))

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	loginBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	loginRes := LoginResponse{}
	if err := json.Unmarshal([]byte(loginBody), &loginRes); err != nil {
		return "", err
	}

	return loginRes.Auth.ClientToken, nil
}

func signHostCert(vault *api.Logical, publicKeyPath string, publicCertDest string, hostSignerPath string) (error) {
	key, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	cert, err := vault.Write(hostSignerPath, map[string]interface{}{
		"public_key": string(key),
		"cert_type": "host"})
	if err != nil {
		return err
	}

	b := []byte(cert.Data["signed_key"].(string))

	err = ioutil.WriteFile(publicCertDest, b, 0640)
	if err != nil {
		return err
	}

	return nil
}

func getVaultClient(vaultAddr string, vaultToken string) (*api.Logical, error) {
	config := &api.Config{
		Address: vaultAddr,
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
		
	client.SetToken(vaultToken)

	return client.Logical(), nil
}

func getConfig() (*Config, error) {
	c := &Config{}

	flag.StringVar(&c.ApproleId, "approle-id", "", "Vault approle role ID")
	flag.StringVar(&c.ApproleSecret, "approle-secret", "", "Vault approle role secret")
	flag.StringVar(&c.HostSignerPath, "host-signer-path", "", "Vault ssh host-signer path to sign certs against")	
	flag.StringVar(&c.VaultAddr, "vault-addr", "", "Vault URL")
	flag.StringVar(&c.PublicKeyPath, "public-key-path", "", "Path to the public key to sign")
	flag.StringVar(&c.CertPath, "cert-path", "", "Path to write the signed cert to")

	flag.Parse()

	v := reflect.ValueOf(*c)
	typeOfS := v.Type()
	for i := 0; i< v.NumField(); i++ {
		if v.Field(i).Interface() == "" {
			return nil, errors.New(fmt.Sprintf("Error: %s is required.", typeOfS.Field(i).Name))
		}	
	}

	return c, nil
}

func main() {
	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	vaultToken, err := approleLogin(config.VaultAddr, config.ApproleId, config.ApproleSecret)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	vault, err := getVaultClient(config.VaultAddr, vaultToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	err = signHostCert(vault, config.PublicKeyPath, config.CertPath, config.HostSignerPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
