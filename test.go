package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
)

func main() {
	keyID := os.Getenv("SDC_KEY_ID")
	endpoint := os.Getenv("SDC_URL")
	accountName := os.Getenv("SDC_ACCOUNT")

	privateKey, err := ioutil.ReadFile(os.Getenv("SDC_KEY_FILE"))
	if err != nil {
		log.Fatal(err)
	}

	sshKeySigner, err := authentication.NewPrivateKeySigner(keyID, privateKey, accountName)
	if err != nil {
		log.Fatal(err)
	}

	client, err := triton.NewClient(endpoint, accountName, sshKeySigner)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	machines, err := client.Machines().ListMachines(context.Background(), &triton.ListMachinesInput{})
	if err != nil {
		log.Fatalf("ListMachines died with %v", err)
	}

	numMachines := 0
	for _, machine := range machines {
		numMachines++
		fmt.Println(fmt.Sprintf("-- Machine: %v", machine.Name))
	}

	fmt.Println("Total: ", numMachines)
}
