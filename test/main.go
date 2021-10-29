package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	os.Setenv("FABRIC_CFG_PATH", "/root/workspace/src/fabric-samples/config")
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	os.Setenv("CORE_PEER_LOCALMSPID", "Org"+"0"+"MSP")
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", "/root/workspace/src/fabric-samples/test-network/organizations/peerOrganizations/org"+"1"+".example.com/peers/peer0.org"+"1"+".example.com/tls/ca.crt")
	os.Setenv("CORE_PEER_MSPCONFIGPATH", "/root/workspace/src/fabric-samples/test-network/organizations/peerOrganizations/org"+"1"+".example.com/users/Admin@org"+"1"+".example.com/msp")
	os.Setenv("CORE_PEER_ADDRESS", "localhost:"+"7051")

	cmd := exec.Command("peer", "channel", "list")
	cmd.Dir = "/root/workspace/src/fabric-samples/test-network"
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
