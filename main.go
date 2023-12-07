package main

import (
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/plonk"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"

	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	var port int
	var otherPorts string
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&otherPorts, "otherPorts", "", "Comma-separated list of other ports to communicate with")
	flag.Parse()

	otherPortsSlice := strings.Split(otherPorts, ",")

	// Start TCP server
	startServer(port, otherPortsSlice)
}

func startServer(port int, otherPorts []string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Server started on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, port, otherPorts)
	}
}

func handleConnection(conn net.Conn, port int, otherPorts []string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		handleMessage(strings.TrimSpace(msg), port, otherPorts)
	}
}

func handleMessage(msg string, port int, otherPorts []string) {
	if msg == "TEST" {
		for _, p := range otherPorts {
			if p != fmt.Sprintf("%d", port) {
				sendToPort(p, "TEST_FROM_PORT\n")
			}
		}
	} else if msg == "TEST_FROM_PORT" {
		fmt.Printf("Port %d received TEST_FROM_PORT\n", port)
	} else {
		fmt.Printf("Port %d received unknown command: %s\n", port, msg)
	}
}

func sendToPort(port, message string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, message)
}

func ObtainProof() {
	var x1, x2, x3, temp, y fr.Element
	// TODO this only works for small numbers because of the moduli
	// x1.SetRandom()
	// x2.SetRandom()
	// x3.SetRandom()
	x1.SetUint64(647585485)
	x2.SetUint64(86978979045)
	x3.SetUint64(13324153467)
	temp.Add(&x1, &x2)
	y.Add(&x3, &temp)

	var circuit MpcAddCircuit

	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &circuit)
	_r1cs := ccs.(*cs.SparseR1CS)
	srs, err := test.NewKZGSRS(_r1cs)
	if err != nil {
		panic(err)
	}

	var w MpcAddCircuit
	w.C.X1 = x1.Uint64()
	w.C.X2 = x2.Uint64()
	w.C.X3 = x3.Uint64()
	w.Y = y.Uint64()

	witnessFull, err := frontend.NewWitness(&w, ecc.BN254.ScalarField())
	if err != nil {
		log.Fatal(err)
	}
	witnessPublic, err := frontend.NewWitness(&w, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		log.Fatal(err)
	}

	pk, vk, err := plonk.Setup(ccs, srs)
	if err != nil {
		log.Fatal(err)
	}

	proof, err := plonk.Prove(ccs, pk, witnessFull)
	if err != nil {
		log.Fatal(err)
	}

	err = plonk.Verify(proof, vk, witnessPublic)
	if err != nil {
		log.Fatal(err)
	}
}
