package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// VariablePathPrefix the prefix to all secrets that contain variable metadata
const VariablePathPrefix = "/jay-internal/dev/miniconfig/variables"

// AllowedSecretsPathPrefix the prefix to all secrets that contain secret metadata. These contain a list of paths
// that are allowed to be maintained by this tool.
const AllowedSecretsPathPrefix = "/jay-internal/dev/miniconfig/secrets"

// GetVariablesFromSecretsManager returns all variable definitions from secrets manager
func GetVariablesFromSecretsManager() ([]Variable, error) {
	/*
	   	   secrets manager allows for 10k of data per secret.
	   	    variables are name value pairs
	   	    name = the variable name, value = description.
	   	        TODO: convert the value to an object with a description property, so we can add more data later.

	   	    if there are a lot of variables, they won't all fit into one secret.
	   	    this finds all secrets with the VariablePathPrefix prefix,
	           and loads them all.
	*/

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var prefixes = []*string{aws.String(VariablePathPrefix)}

	filter := &secretsmanager.Filter{Key: aws.String(secretsmanager.FilterNameStringTypeName), Values: prefixes}
	var filters []*secretsmanager.Filter
	filters = append(filters, filter)

	input := &secretsmanager.ListSecretsInput{Filters: filters}

	service := secretsmanager.New(sess)
	result, err := service.ListSecrets(input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// get and iterate the secrets that begin with the prefix
	var variables []Variable
	for _, s := range result.SecretList {
		fmt.Println(*s.Name)

		// for each secret, get the string value from secrets manager
		var valueRequest = &secretsmanager.GetSecretValueInput{SecretId: s.Name}
		value, errValue := service.GetSecretValue(valueRequest)
		if errValue != nil {
			fmt.Println(errValue)
			return nil, errValue
		}

		fmt.Println("\t" + *value.SecretString)

		// load the json into a map, then
		// iterate the map and convert each item to
		// a Variable.
		var doc map[string]string
		errDeserialize := json.Unmarshal([]byte(*value.SecretString), &doc)
		if errDeserialize != nil {
			fmt.Println(errDeserialize)
			return nil, errDeserialize
		}

		for key, value := range doc {
			// TODO: throw exception if same variable
			// is encountered twice
			v := Variable{Name: key, Description: value}
			variables = append(variables, v)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Println(variables)

	return variables, nil
}
