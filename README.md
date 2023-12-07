# MPC with ZKP

Shamir Secret Sharing with Zero Knowledge proofs proving the correct operations were performed. [gnark](https://github.com/Consensys/gnark/tree/4199bb354d89c82afce06d801424cbbeae222a1b) is used to create the ZKP. 

This uses a similar approach to https://github.com/hashcloak/shamir-rs, and adds zero knowledge proofs at the points where a calculation is done by a party. 

## Build & Run

```
go run . -port=8080 -otherPorts=8081,8082
go run . -port=8081 -otherPorts=8080,8082
go run . -port=8082 -otherPorts=8080,8081

echo "TEST" | nc 127.0.0.1 8080
```

## Test

Run tests:
```
go test
```

## Next steps

Handle command with given functionality:
- SEND_SHARES - calculate shares and send RECEIVE_SHARE to other parties
- RECEIVE_SHARE - store received share
- SUM_AND_DISTRIBUTE - sum received shares, create ZKP and send both the share and the ZKP to other parties with RECEIVE_SUM_AND_PROOF
- RECEIVE_SUM_AND_PROOF - verify proof and if it checks out, store sum. Otherwise, store an error (?)
- GIVE_RESULT - give result if there was no error

Research & fix: the zkp can only handle small numbers, not random numbers. 

Optional:
- generate https://pkg.go.dev/github.com/consensys/gnark-crypto/ecc/bn254/fr from goff

<!-- ## Tutorial

The goal is have some inital practice with gnark.

### 0. Init project

```
mkdir zkp-mpc
cd zkp-mpc
go mod init zkp-mpc
```

### 1. Shamir Secret Sharing

We choosing the Scalar Field of bn254 as the field to so Shamir Secret Sharing over. The package can be found on the Go Package Discovery site [here](https://pkg.go.dev/github.com/consensys/gnark-crypto/ecc/bn254/fr), and the source code can be found [here](https://github.com/consensys/gnark-crypto). 

Import the package into the project
```
go get github.com/consensys/gnark-crypto/ecc/bn254/fr
go mod tidy
```

Create a new file `shamir.go`. This will belong to `package main` and needs to import the package we just added. We're going to create functions `GenerateSecret`, `GetSharesSecret` and `CreatePol` here. 

Create a new file `main.go`, which also belongs to `package main` and must contain the `main` function. 

Both files will be added directly to the newly created directory `zkp-mpc`. 

`main.go`: 
```go
package main

func main() {

}
```

`shamir.go` with the function `GenerateSecret` defined:
```go
package main

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func GenerateSecret() fr.Element {
	var s fr.Element
	s.SetRandom()
	return s
}
```

### 2. Add gnark

```
go get github.com/consensys/gnark@latest
go mod tidy
```

We'll make a new file `zkp.go` and add the gnark package:

```go
package main

import (
	"github.com/consensys/gnark/frontend"
)

//..
``` -->
