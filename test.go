package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
)

func main() {
	keyID := os.Getenv("SDC_KEY_ID")
	accountName := os.Getenv("SDC_ACCOUNT")

	file, err := os.Open(os.Getenv("SDC_KEY_FILE"))
	if err != nil {
		log.Fatal(err)
	}

	privateKey := make([]byte, 5000)
	_, err = file.Read(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	sshKeySigner, err := authentication.NewPrivateKeySigner(keyID, privateKey, accountName)
	if err != nil {
		log.Fatalf("Fatal exception from NewPrivateKeySigner: %s", err)
	}

	// sshKeySigner, err := authentication.NewSSHAgentSigner(os.Getenv("SDC_KEY_ID"), os.Getenv("SDC_ACCOUNT"))
	// if err != nil {
	// 	log.Fatalf("NewSSHAgentSigner: %s", err)
	// }

	client, err := triton.NewClient(os.Getenv("SDC_URL"), os.Getenv("SDC_ACCOUNT"), sshKeySigner)
	// client, err := triton.NewClient(os.Getenv("SDC_URL"), os.Getenv("SDC_ACCOUNT"), sshKeySigner)
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
