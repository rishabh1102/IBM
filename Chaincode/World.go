package main

import (
	"errors"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//Defining Structure for the Tokens
type Token struct{
	Id string `json:"id"`
	Colour string `json:"colour"`
	User string `json:"user"`
	Sell bool `json:"sell"`
	Value int `json:"value"`
	Trade bool `json:"trade"`
	TradeColour string `json:"tradecolour"`
	SaleString string `json:"salestring"`
	TradeString string `json:"tradestring"`
}

//Defining structure for the players
type Player struct {
	Id string `json:"id"`		//Password
	Points string `json:"points"`
	Assets string `json:"assets"`
	Name string `json:"name"`
	Photo string `json:"photo"`
	TokenValue string `json:"tokenvalue"`
}


//"<token1ID>,<token2ID>,<Token3ID>...."
var tokenList = "tokenList"
//"<User1ID>:<TokenWillingID>|<TokenWillingColour>-<ColourRequired>,<User2ID>:<TokenWillingID>|<TokenWillingColour>-<ColourRequired>"
var tradeList = "tradeList"
//[<SellerID>:<SellerTokenID>|<SellerTokenColour>-<SaleValue>,...]
var userList = "userList"
var saleList = "saleList"

//Init
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	//We Initialize all the lists

	//Initializing User List
	err := stub.PutState(userList, []byte(""))
	if err != nil {
		return nil, err
	}

	//Initializing Token List
	err = stub.PutState(tokenList, []byte(""))
	if err != nil {
		return nil, err
	}

	//Initializing Trade List
	err = stub.PutState(tradeList, []byte(""))
	if err != nil {
		return nil, err
	}

	//Initializing Sale List
	err = stub.PutState(saleList, []byte(""))
	if err != nil {
		return nil, err
	}

	return nil, nil
}


func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	// Handle different functions
	if function == "init" {											
		return t.Init(stub, "init", args)
	} else if function == "delete" {								
		//return t.delete(stub, args)
	} else if function == "createToken" {							
		return t.createToken(stub, args)
	} else if function == "createUser" {								
		return t.createUser(stub, args)
	} else if function == "claimUser" {										
		//return t.setUser(stub, args)
	} else if function == "claimToken" {										
		//return t.setUser(stub, args)
	} else if function == "setTradeStatus" {									
		//return t.setTradeStatus(stub, args)
	} else if function == "setSellStatus" {										
		//return t.setSellStatus(stub, args)
	} else if function == "trade" {									
		//return t.trade(stub, args)
	} else if function == "buy" {									
		//return t.buy(stub, args)
	} else if function == "redeem" {									
		//return t.buy(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}


//args[0] = tokenID, args[1] = Token colour
func (t *SimpleChaincode) createToken(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//Id, Colour, User, Sell, Value, Trade, TradeColour
	
	//Checking Errors
	if len(args) < 2 {
		return nil, errors.New("Expecting At least 2 arguments")
	}
	if len(args[0]) <= 0 {
		return nil, errors.New("ID must be non null")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("Colour must be non null")
	}

	//Creating Json
	jsonString := `{"id": "` + args[0] + `", "colour": "` + args[1] + `", "user": "", "sell": "` + "false" + `", "sellvalue": "0", "trade": "false", "tradecolour" : "", "salestring": "", "tradestring": ""}`
	
	//Getting Token list to check
	getTokenList, err := stub.GetState(tokenList)
	if err != nil {
		return nil, errors.New("Failed to get Token ID")
	}

	//Checking if ID already exists
	if (searchString(string(getTokenList), args[0])) {
		fmt.Println("Token Already Exists!!")
		return nil, errors.New("Token Already Exists")
	}

	//Writing to Ledger
	err = stub.PutState(args[0], []byte(jsonString))
	if err != nil {
		return nil, err
	}
		
	//Updating tokenList
	updatedTokenList := addSubstringtoString(string(getTokenList), args[0])

	//Writing Back Updated List
	err = stub.PutState(tokenList, []byte(updatedTokenList))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//args[0] = userID, args[1] = Password, args[2] = Name
func (t *SimpleChaincode) createUser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//Id, Points, Assets, Name, Photo
	
	//Checking Errors
	if len(args) != 3 {
		return nil, errors.New("Expecting Only 1")
	}

	//Getting user list to check for errors
	getUserList, err := stub.GetState(userList)
	if err != nil {
		return nil, errors.New("Failed to get User List")
	}

	//Checking if user already exists
	if (searchString(string(getUserList), args[0])) {
		return nil, errors.New("User Already Exists")
	}
	
	//Creating Json
	jsonString := `{"id": "` + args[1] + `", "points": "` + "50" + `", "assets": "", "name": "` + args[2] + `", "photo": "", "tokenvalue": "0"}`
	
	//Writing to Ledger
	err = stub.PutState(args[0], []byte(jsonString))
	if err != nil {
		return nil, err
	}
	
	var updatedUserList string

	//Updating User List
	updatedUserList = addSubstringtoString(string(getUserList), args[0])

	//Writing Back Updated List
	err = stub.PutState(userList, []byte(updatedUserList))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//args[0] = userID, args[1] = name of new User, args[2] = Photo Encoding
func (t *SimpleChaincode) claimUser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//Getting user list to check for errors
	getUserList, err := stub.GetState(userList)
	if err != nil {
		return nil, errors.New("Failed to get User List")
	}

	//Checking if user exists
	if (!searchString(string(getUserList), args[0])) {
		return nil, errors.New("User Does not exist")
	}

	//Id, Points, Assets, Name, Photo
	userAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get user details")
	}
	tempUser := Player{}
	json.Unmarshal(userAsBytes, &tempUser)
	tempUser.Name = args[1]
	tempUser.Photo = args[2]

	tempUserWriteBack, _ := json.Marshal(tempUser)
	err = stub.PutState(args[0], tempUserWriteBack)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

//args[0] = tokenID and args[1] is the User Claiming
func (t *SimpleChaincode) claimToken(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//Getting Token list to check
	getTokenList, err := stub.GetState(tokenList)
	if err != nil {
		return nil, errors.New("Failed to get Token ID")
	}

	//Checking if ID exists
	if (!searchString(string(getTokenList), args[0])) {
		return nil, errors.New("This Token Does not exist")
	}

	//Getting user list to check for errors
	getUserList, err := stub.GetState(userList)
	if err != nil {
		return nil, errors.New("Failed to get User List")
	}

	if (!searchString(string(getUserList), args[1])) {
		return nil, errors.New("UserID Does not exist")
	}

	//Id, Colour, User, Sell, Value, Trade, TradeColour
	tokenAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get user details")
	}
	//Edit Token Entry
	tempToken := Token{}
	json.Unmarshal(tokenAsBytes, &tempToken)

	//Set User
	tempToken.User = args[1]

	tempTokenWriteBack, _ := json.Marshal(tempToken)
	err = stub.PutState(args[0], tempTokenWriteBack)

	if err != nil {
		return nil, err
	}

	//Edit user Entry
	tempUser, err := getPlayerFromID(stub, args[1])
	if err!= nil {
		return nil, err
	}

	//Add Token to his assets
	tempUser.Assets = addSubstringtoString(tempUser.Assets, args[0])
	tempUser.TokenValue = strconv.Itoa(updateTokenValue(stub, tempUser.Assets))
	tempUserWriteBack, _ := json.Marshal(tempUser)
	err = stub.PutState(args[1], tempUserWriteBack)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

//args[0] = UserId, args[1]  = TokenID to trade, args[2] = Colour wanted in return.
func (t *SimpleChaincode) setTradeStatus(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	//Get User Data
	tempUser, err := getPlayerFromID(stub, args[0])
	if err != nil {
		return nil, errors.New("Failed to get user details")
	}

	//Checking User Assets
	if searchString(tempUser.Assets, args[1]) == false {
		return nil, errors.New("Asset does not exist")
	}

	//Getting Token Data
	tempToken, err := getTokenFromID(stub, args[1])
	if err != nil {
		return nil, errors.New("Failed to get token details")
	}

	if tempToken.Trade == true {
		return nil, errors.New("Token already has existing Trade")
	}

	//Updating Trade List
	var putString = args[0] + ":" + args[1] + "|" + tempToken.Colour + "-" + args[2]

	//Updating Token Entry
	tempToken.Trade = true
	tempToken.TradeString = putString
	tempToken.TradeColour = args[2]

	//Writing back token to blockchain
	tempTokenWriteBack, _ := json.Marshal(tempToken)
	err = stub.PutState(args[1], tempTokenWriteBack)
	if err != nil {
		return nil, err
	}

	getTradeList, err := stub.GetState(tradeList)
	if err != nil {
		return nil, errors.New("Failed to get Trade List")
	}
	//Writing Back Updated List
	updatedTradeList := addSubstringtoString(string(getTradeList), putString)

	err = stub.PutState(tradeList, []byte(updatedTradeList))
	if err != nil {
		return nil, err
	}

	return nil,nil

}

//args[0] = UserID, args[1] = TokenID to sell, args[2] = Value
func (t *SimpleChaincode) setSellStatus(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
		
	//Get User Data
	userAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get user details")
	}

	tempUser := Player{}
	json.Unmarshal(userAsBytes, &tempUser)

	//Checking User Assets
	if !(searchString(tempUser.Assets, args[1])) {
		return nil, errors.New("Asset does not exist for user")
	}

	//Getting Token Data
	tempToken, err := getTokenFromID(stub, args[1])
	if err != nil {
		return nil, errors.New("Token entry couldn't be fetched")
	}

	if tempToken.Sell == true {
		return nil, errors.New("Token is already on Sale")
	}

	//Updating Sale List
	putString := args[0] + ":" + args[1] + "|" + tempToken.Colour + "-" + args[2]

	//Updating Token Entry
	tempToken.Sell = true
	tempToken.Value, _ = strconv.Atoi(args[2])
	tempToken.SaleString = putString
	tempTokenWriteBack, _ := json.Marshal(tempToken)
	err = stub.PutState(args[1], tempTokenWriteBack)
	if err != nil {
		return nil, err
	}

	

	getSaleList, err := stub.GetState(saleList)
	if err != nil {
		return nil, errors.New("Failed to get Sale List")
	}
	updatedSaleList := addSubstringtoString(string(getSaleList), putString)

	//Writing Back Updated List
	err = stub.PutState(saleList, []byte(updatedSaleList))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//args[0] = buyerID, args[1] = buyerTokenID, args[2]=buyerTokenColour args[3] = SellerID, args[4] = SellerTokenID, args[5] = SellerTokenColour
func (t *SimpleChaincode) trade(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 6 {
		return nil, errors.New("6 Arguements Required")
	}
	
	var tradeString = args[3] + ":" + args[4] + "|" + args[5] + "-" + args[2]
	tradeListString, err := stub.GetState(tradeList)
	if err != nil {
		return nil, err
	}

	//Searching Trade String to check if trade exists
	if (!searchString(string(tradeListString), tradeString)) {
		return nil, errors.New("Trade not found")
	}

	//Checking that both tokens belong to owners and getting Player Data
	checkA, tempBuyer, errA := userCheckTokenOwnership(stub, args[0], args[1])
	checkB, tempSeller, errB := userCheckTokenOwnership(stub, args[3], args[4])
	if errA != nil {
		return nil, errA
	}
	if errB != nil {
		return nil, errB
	}
	if (!(checkA && checkB)) {
		return nil, errors.New("Problem with token ownership")
	}


	//Getting Seller Token Data
	tempSellerToken, err := getTokenFromID(stub, args[4])
	if err != nil {
		return nil, err
	}

	//Getting Buyer Token Data
	tempBuyerToken, err := getTokenFromID(stub, args[1])
	if err != nil {
		return nil, err
	}

	//Creating new Trade List after omitting current trade
	var newTradeString = removeSubstringFromString(string(tradeListString), tradeString)

	//Clearing out Other Token's entry from the trade list
	if (tempBuyerToken.TradeString != "") {
		newTradeString = removeSubstringFromString(newTradeString, tempBuyerToken.TradeString)
	}


	//Removing Entries of traded tokens from Sale List
	if (tempSellerToken.Sell || tempBuyerToken.Sell) {
		getSellList, err := stub.GetState(saleList)
		if err != nil {
			return nil, errors.New("Couldn't find saleList")
		}
		var updatedSaleList string
		if (tempSellerToken.Sell) && (tempBuyerToken.Sell) {
			updatedSaleList = removeSubstringFromString(string(getSellList), tempSellerToken.SaleString)
			updatedSaleList = removeSubstringFromString(updatedSaleList, tempBuyerToken.SaleString)
		} else if (tempBuyerToken.Sell) {
			updatedSaleList = removeSubstringFromString(string(getSellList), tempBuyerToken.SaleString)
		} else {
			updatedSaleList = removeSubstringFromString(string(getSellList), tempSellerToken.SaleString)
		}
		err = stub.PutState(saleList, []byte(updatedSaleList))
	}

	//Making the trade
	tempSeller.Assets = removeSubstringFromString(tempSeller.Assets, args[4])
	tempSeller.Assets = addSubstringtoString(tempSeller.Assets, args[1])
	tempSeller.TokenValue = strconv.Itoa(updateTokenValue(stub, tempSeller.Assets))
	tempBuyer.Assets = removeSubstringFromString(tempBuyer.Assets, args[1])
	tempBuyer.Assets = addSubstringtoString(tempBuyer.Assets, args[4])
	tempBuyer.TokenValue = strconv.Itoa(updateTokenValue(stub, tempBuyer.Assets))
	tempSellerToken.User = args[0]
	tempSellerToken.Trade = false
	tempSellerToken.Sell = false
	tempSellerToken.SaleString = ""
	tempSellerToken.TradeString = ""
	tempSellerToken.TradeColour = ""
	tempSellerToken.Value = 0
	tempBuyerToken.User = args[3]
	tempBuyerToken.Trade = false
	tempBuyerToken.Sell = false
	tempBuyerToken.SaleString = ""
	tempBuyerToken.TradeString = ""
	tempBuyerToken.TradeColour = ""
	tempBuyerToken.Value = 0

	//Writing Back all data
	writeBack, _ := json.Marshal(tempSellerToken)
	err = stub.PutState(args[4], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(tempBuyerToken)
	err = stub.PutState(args[1], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(tempBuyer)
	err = stub.PutState(args[0], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(tempSeller)
	err = stub.PutState(args[3], writeBack)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(tradeList, []byte(newTradeString))
	if err != nil {
		return nil, err
	}

	//Trade Complete
	return nil, nil
}

//args[0] = buyerID, args[1] = SellerID, args[2] = SellerTokenID, args[3] = SellerTokenColour, args[4] = Value
func (t *SimpleChaincode) buy(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Exactly 5 Arguements Required")
	}

	//Obtaining List from Blockchain
	saleListString, err := stub.GetState(saleList)
	if err != nil {
		return nil, err
	}

	//Searching Sell String to check if Sale exists
	var sellString = args[1] + ":" + args[2] + "|" + args[3] + "-" + args[4]
	if (!searchString(string(saleListString), sellString)) {
		return nil, errors.New("Sale not found")
	}

	//Checking ownership of Token for seller
	check, sellerUser, err := userCheckTokenOwnership(stub, args[1], args[2])
	if err != nil {
		return nil, err
	}	
	if (!check) {
		return nil, errors.New("Token doesn't belong to user")
	}
	
	//Checking if Buyer has sufficient balance to buy
	buyerUser, err := getPlayerFromID(stub, args[0])
	reqBuyerPoints, _ := strconv.Atoi(args[4])
	buyerPoints, _ := strconv.Atoi(buyerUser.Points)
	if buyerPoints < reqBuyerPoints {
		return nil, errors.New("Buyer does not have enough points")
	}
	
	sellerToken, err := getTokenFromID(stub, args[2])
	if err != nil {
		return nil, err
	}

	//Making New Sell String after omitting this sale
	var newSellString = removeSubstringFromString(string(saleListString), sellString)

	//Removing Redundant Entries from Trade List
	if (sellerToken.Trade) {
		getTradeList, err := stub.GetState(tradeList)
		if err != nil {
			return nil, errors.New("Couldn't find tradeList")
		}
		updatedTradeList := removeSubstringFromString(string(getTradeList), sellerToken.TradeString)
		err = stub.PutState(saleList, []byte(updatedTradeList))
	}


	//Performing the sale
	sellerPoints, _ := strconv.Atoi(sellerUser.Points)
	sellerPoints = sellerPoints + reqBuyerPoints
	buyerPoints = buyerPoints - reqBuyerPoints
 	sellerUser.Assets = removeSubstringFromString(sellerUser.Assets, args[2])
 	sellerUser.TokenValue = strconv.Itoa(updateTokenValue(stub, sellerUser.Assets))
 	sellerUser.Points = strconv.Itoa(sellerPoints)
 	buyerUser.Assets = addSubstringtoString (buyerUser.Assets, args[2])
 	buyerUser.Points = strconv.Itoa(buyerPoints)
 	buyerUser.TokenValue = strconv.Itoa(updateTokenValue(stub, buyerUser.Assets))
 	sellerToken.User = args[0]
 	sellerToken.Trade = false
 	sellerToken.TradeString = ""
 	sellerToken.Sell = false
 	sellerToken.SaleString = ""
 	

 	//Writing Back to the blockchain
 	writeBack, _ := json.Marshal(sellerToken)
	err = stub.PutState(args[2], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(sellerUser)
	err = stub.PutState(args[1], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(buyerUser)
	err = stub.PutState(args[0], writeBack)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(saleList, []byte(newSellString))
	if err != nil {
		return nil, err
	}
	//Completed Sale
	return nil, nil
}

//The function to redeem combinations from the bank
func (t *SimpleChaincode) redeem(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {


	// Handle different functions
	if function == "init" {											
		return t.Init(stub, "init", args)
	} else if function == "delete" {								
		//return t.delete(stub, args)
	} else if function == "createToken" {							
		return t.createToken(stub, args)
	} else if function == "createUser" {								
		return t.createUser(stub, args)
	} else if function == "claimUser" {										
		return t.claimUser(stub, args)
	} else if function == "claimToken" {										
		return t.claimToken(stub, args)
	} else if function == "setTradeStatus" {									
		return t.setTradeStatus(stub, args)
	} else if function == "setSellStatus" {										
		return t.setSellStatus(stub, args)
	} else if function == "trade" {									
		return t.trade(stub, args)
	} else if function == "buy" {									
		return t.buy(stub, args)
	} else if function == "redeem" {									
		//return t.buy(stub, args)
	}
	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
	return nil, nil

}

//function to check if token belongs to owner
func userCheckTokenOwnership (stub *shim.ChaincodeStub, userID string, tokenID string) (bool, Player, error) {

	tempUser, err := getPlayerFromID(stub, userID)
	if err != nil {
		errorUser := Player{}
		return false, errorUser, err
	}
	tempUserAssets := strings.Split(tempUser.Assets, ",")
	for i := 0 ; i < len(tempUserAssets) ; i++ {
		if tempUserAssets[i] == tokenID {
			return true, tempUser, err
		}
	}
	return false, tempUser, err
}

//function to remove a certain substring from a string
func removeSubstringFromString (fullString string, substring string ) string {
	returnString := ""
	listString := strings.Split(fullString, ",")
	for i := 0 ; i < len(listString) ; i++ {
		if listString[i] != substring {
			if returnString == "" {
				returnString = returnString + listString[i]
			} else {
				returnString = returnString + "," + listString[i]
			}
		}
	}
	return returnString
}

//function which returns token on passing the ID
func getTokenFromID (stub *shim.ChaincodeStub, tokenID string) (Token, error) {
	tokenAsBytes, err := stub.GetState(tokenID)
	if err != nil {
		errorToken := Token{}
		return errorToken, errors.New("Failed to get user details")
	}
	tempToken := Token{}
	json.Unmarshal(tokenAsBytes, &tempToken)
	return tempToken, err
}

//Function which returns player on passing the ID
func getPlayerFromID (stub *shim.ChaincodeStub, playerID string) (Player, error) {
	userAsBytes, err := stub.GetState(playerID)
	if err != nil {
		errorPlayer := Player{}
		return errorPlayer, err
	}
	tempUser := Player{}
	json.Unmarshal(userAsBytes, &tempUser)
	return tempUser, err
}

//function which adds string to original string in the "a,b,c" format
func addSubstringtoString (fullString string , substring string ) string {
	if fullString == "" {
		return substring
	} else {
		return fullString + "," + substring
	}
}

//search for a given string in another string.
func searchString (fullString string , searchVal string ) bool {
	listString := strings.Split(fullString, ",")
	for i:= 0 ; i < len(listString) ; i++ {
		if listString[i] == searchVal {
			return true
		}
	}
	return false
}


func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//Checking Number of arguements
	// if len(args) != 1 {
	// 	return nil, errors.New("Only 1 Argument")
	// }
	//MAKE CHECK FOR NUMBER OF ARGUMENTS
	
	//Getting value from blockchain
	returnVal, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}

	if function == "entryExist" {	
		if(string(returnVal) == "") {
			return []byte("FALSE"), err
		} else {										
			return []byte("TRUE"), err
		}
	} else if function == "queryTokenOwnership" {
		return t.queryTokenOwnership(stub, args)
	} else if function == "checkLoginDetails" {
		return t.checkLoginDetails(stub, args)
	} else if function == "getLeaderBoard" {
		return t.getLeaderBoard(stub, args)
	} else {
		if(string(returnVal) == "") {
			return []byte("FALSE"), err
		}
		return returnVal, err
	}
}

func (t *SimpleChaincode) queryTokenOwnership(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//Getting value from blockchain
	returnVal, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}

	if string(returnVal) == "" {
		return []byte("FALSE"), err
	} else {
		tempToken := Token{}
		json.Unmarshal(returnVal, &tempToken)
		if tempToken.User == "" {
			return []byte("TRUE"), err
		} else {
			return []byte("FALSE"), err
		}
	}
	return nil, nil
}

func (t *SimpleChaincode) checkLoginDetails(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//Getting value from blockchain
	returnVal, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}

	if string(returnVal) == "" {
		return []byte("FALSE"), err
	} else {
		tempPlayer := Player{}
		json.Unmarshal(returnVal, &tempPlayer)
		if tempPlayer.Id == args[1] {
			return []byte("TRUE"), nil
		} else {
			return []byte("FALSE"), nil
		}
	}
	return nil, nil
}

func (t *SimpleChaincode) getLeaderBoard(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//Getting value from blockchain
	getUserList, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}
	tempPlayer := Player{}
	returnString := ""
	if string(getUserList) == "" {
		return []byte("No Users Enrolled"), err
	} else {
		usernameList := strings.Split(string(getUserList), ",")
		for i := 0 ; i < len(usernameList) ; i++ {
			if returnString != "" {
				returnString = returnString + ","
			}
			returnString = returnString + usernameList[i]
			returnString = returnString + "-"
			getPlayer, err := stub.GetState(usernameList[i])
			if err != nil {
				return nil, err
			}
			tempPlayer = Player{}
			json.Unmarshal(getPlayer, &tempPlayer)
			returnString = returnString + tempPlayer.Points
			returnString = returnString + "|"
			returnString = returnString + tempPlayer.TokenValue
		}
		return []byte(returnString), nil		
	}
	return nil, nil
}


func updateTokenValue(stub *shim.ChaincodeStub, assetStr string) int {
	assetList := strings.Split(assetStr, ",")
	if len(assetList) == 0 {
		return 0
	}
	tempToken := Token{}
	userTokenList := []Token{}
	for i := 0 ; i < len(assetList) ; i++ {
		tokenAsBytes, err := stub.GetState(assetList[i])
		if err != nil {
			return -1
		}
		tempToken = Token{}
		json.Unmarshal(tokenAsBytes, &tempToken)
		userTokenList = append(userTokenList, tempToken)
	}
	var points int = 0 
	var multiplier int = 0
	//index: 0 - Red, 1 - Yellow, 2 - Orange, 3 - Blue, 4 - Green
	score := []int{0,0,0,0,0}
	multiArray := []int{1,1,1,1,1}
	//multiVal := []int{1, 1, 1, 1, 1}
	for i := 0; i < len(userTokenList) ; i++ {
		switch userTokenList[i].Colour {
		case "red":
			multiArray[0] = multiArray[0] * 2
			score[0] = score[0] + 1
		case "yellow":
			multiArray[1] = multiArray[1] * 2
			score[1] = score[1] + 2
		case "orange":
			multiArray[2] = multiArray[2] * 2
			score[2] = score[2] + 3
		case "blue":
			multiArray[3] = multiArray[3] * 2
			score[3] = score[3] + 4
		case "green":
			multiArray[4] = multiArray[4] * 2
			score[4] = score[4] + 5
		}
	}
	// for i := 0 ; i < len(multiArray) ; i++ {
	// 	multiArray[i] = multiArray[i]/2
	// } 
	totalScore := 0
	for i := 0 ; i < len(multiArray) ; i++ {
		multiArray[i] = multiArray[i]/2
		if (multiArray[i] > 1) {
			multiplier = multiplier + (multiArray[i])
		}
		totalScore = totalScore + score[i]
	}
	if multiplier == 0 {
		multiplier = 1
	}
	fmt.Println(totalScore)
	fmt.Println(multiplier)
	points = multiplier*totalScore
	return points
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}


