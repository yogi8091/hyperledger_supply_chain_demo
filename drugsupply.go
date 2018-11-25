/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
        "time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the asset structure, with  properties.  Structure tags are used by encoding/json library
type Asset struct {

	AssetId			string
	ManufactureName string
	BatchNo 		string
	ItemName 		string
	Quantity 		uint64
	ExpiryDate		string
	Owner           string
	OwnerType		string

}


/*
 * The Init method is called when the Smart Contract "drugChain" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "drugChain"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createAsset" {
		return s.createAsset(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "changeAssetOwner" {
		return s.changeAssetOwner(APIstub, args)
	} else if function == "getItemHistory"{
		return s.getItemHistory(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) getItemHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetId := args[0]

	fmt.Printf("- asset to be retrieved: %s\n", assetId)

	resultsIterator, err := APIstub.GetHistoryForKey(assetId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true

		}
		buffer.WriteString("]")

		fmt.Printf("- getHistoryForItem returning:\n%s\n", buffer.String())

		return shim.Success(buffer.Bytes())
}

func (s *SmartContract) createAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

i, err := strconv.ParseUint(args[4], 10, 64)
if err != nil {
     return shim.Error(err.Error())
}
	
	var asset = Asset{ManufactureName: args[1], BatchNo: args[2], ItemName: args[3], Quantity: i, ExpiryDate: args[5], Owner: args[6], OwnerType: args[7]}

	assetAsBytes, _ := json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(nil)
}







func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	assets := []Asset{
		Asset{ManufactureName: "ReddyLabs", BatchNo: "1/2018", ItemName: "paracetmol", Quantity: 100, ExpiryDate: "01/01/2020", Owner: "ReddyLabs", OwnerType: "Manufacturer"},
		Asset{ManufactureName: "ReddyLabs", BatchNo: "1/2018", ItemName: "paracetmol", Quantity: 100, ExpiryDate: "01/01/2020", Owner: "ReddyLabs", OwnerType: "Manufacturer"},
		Asset{ManufactureName: "ReddyLabs", BatchNo: "1/2018", ItemName: "paracetmol", Quantity: 100, ExpiryDate: "01/01/2020", Owner: "ReddyLabs", OwnerType: "Manufacturer"},
		Asset{ManufactureName: "ReddyLabs", BatchNo: "1/2018", ItemName: "paracetmol", Quantity: 100, ExpiryDate: "01/01/2020", Owner: "ReddyLabs", OwnerType: "Manufacturer"},
		Asset{ManufactureName: "ReddyLabs", BatchNo: "1/2018", ItemName: "paracetmol", Quantity: 100, ExpiryDate: "01/01/2020", Owner: "ReddyLabs", OwnerType: "Manufacturer"},
		Asset{ManufactureName: "ReddyLabs", BatchNo: "1/2018", ItemName: "paracetmol", Quantity: 100, ExpiryDate: "01/01/2020", Owner: "ReddyLabs", OwnerType: "Manufacturer"},
	}

	i := 0
	for i < len(assets) {
		fmt.Println("i is ", i)
		assetAsBytes, _ := json.Marshal(assets[i])
		APIstub.PutState("ASSET"+strconv.Itoa(i), assetAsBytes)
		fmt.Println("Added", assets[i])
		i = i + 1
	}

	return shim.Success(nil)
}



func (s *SmartContract) changeAssetOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	assetAsBytes, _ := APIstub.GetState(args[0])
	asset := Asset{}

	json.Unmarshal(assetAsBytes, &asset)
	asset.Owner = args[1]

	assetAsBytes, _ = json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

