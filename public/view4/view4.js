'use strict';

angular.module('myApp.view4', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view4', {
    templateUrl: 'view4/view4.html',
    controller: 'View4Ctrl'
  });
}])

.controller('View4Ctrl', function($scope, productService, $http, $location) {
	$scope.username = productService.getProducts()[0];
	if ($scope.username == null) {
		$location.path("/view1")
	}
	//$scope.username = "cyrus";

	function createJsonString(typeStr, funcStr, args) {
		var retString = {
		  "jsonrpc": "2.0",
		  "method": typeStr,
		  "params": {
		    "type": 1,
		    "chaincodeID":{
		        "name": "mycc"
		    },
		    "ctorMsg": {
		        "function": funcStr,
		        "args": args
		    },
		    "secureContext": "lukas"
		  },
		  "id": 1
		}
		return retString;
	}

	var sendQuery = createJsonString("query", "query", [$scope.username]);

	$http({
	  url: 'http://9.232.130.202:3000/chaincode',
	  method: 'POST',
	  data: sendQuery
	})
	//if Query Gets response successfully
	.then(function(queryResponse) {
		$scope.assets = []
		var userQueryResult = angular.fromJson(queryResponse.data.result.message)
		$scope.name = userQueryResult.name
		$scope.points = userQueryResult.points
		if (userQueryResult.assets == "") {
			$scope.assets = [{"id": "N/A", "colour": "N/A"}]
			$scope.assetString = "None"
		}
		else {
			var assetsStr = userQueryResult.assets;
			var assetsList = assetsStr.split(',');
			//$scope.val = assetsList
			for ( var i = 0 ; i < assetsList.length; i++ ) {
				var queryString = createJsonString("query", "query", [assetsList[i]]);
				$http({
				  url: 'http://9.232.130.202:3000/chaincode',
				  method: 'POST',
				  data: queryString
				})
				.then(function(tokenQueryResponse) {
					var tokenQueryResult = angular.fromJson(tokenQueryResponse.data.result.message);
					if (tokenQueryResult.trade == false) {
						tokenQueryResult.tradeButton = "Trade";
					}
					else {
						tokenQueryResult.tradeStatus = "Exchange: " + tokenQueryResult.tradecolour;
					}
					if (tokenQueryResult.sell == false) {
						tokenQueryResult.sellButton = "Sell";
					}
					else {
						tokenQueryResult.sellStatus = "Value: " + tokenQueryResult.value + " pt";
					}
					$scope.assets.push(tokenQueryResult)
				},
				function(tokenQueryResponse) {
					$scope.assets = "Error Connecting to Blockchain"
				})
			}
			//$scope.name = $scope.assets
		}
	}, 
	//if Query doesn't get a response
	function(queryResponse) { // optional
	        $scope.query = "Could not connect to blockchain"
	});

	var tradeListQuery = createJsonString("query", "query", ["tradeList"]);
	//$scope.val = tradeListQuery;

	$http({
	  url: 'http://9.232.130.202:3000/chaincode',
	  method: 'POST',
	  data: tradeListQuery
	})
	//if Query Gets response successfully
	.then(function(queryResponse) {
		var tradeListStr = queryResponse.data.result.message;
		$scope.Trades = []
		var tradeList = tradeListStr.split(',');

		//Each Item in Trade list is: "<User1ID>:<TokenWillingID>|<TokenWillingColour>-<ColourRequired>""
		//Creating tradeList Structure where each item is:
		//["Index", "SellerID", "SellerTokenID", "SellerTokenColour", "SellerExchangeColour"]
		for ( var i = 0 ; i < tradeList.length; i++ ) {
					var tradeItem = []
					//Pushing Index
					tradeItem.push(i+1);
					var tempStr = tradeList[i].split(':');
					//Pushing Username
					tradeItem.push(tempStr[0]);
					var tempStr2 = tempStr[1].split('|');
					//Pushing TokenID
					tradeItem.push(tempStr2[0]);
					var tempStr3 = tempStr2[1].split('-');
					//Pushing Token Colour
					tradeItem.push(tempStr3[0]);
					//Pushing Exchange colour
					tradeItem.push(tempStr3[1]);
					$scope.Trades.push(tradeItem);
				}
		//$scope.val = $scope.TradeUser;
	}, 
	//if Query doesn't get a response
	function(queryResponse) { // optional
	        $scope.query = "Could not connect to Blockchain"
	});
	//END OF TRADE

	//BEGIN SALE LIST
	var saleListQuery = createJsonString("query", "query", ["saleList"]);
	//$scope.val = saleListQuery;

	$http({
	  url: 'http://9.232.130.202:3000/chaincode',
	  method: 'POST',
	  data: saleListQuery
	})
	//if Query Gets response successfully
	.then(function(queryResponse) {
		var saleListStr = queryResponse.data.result.message;
		$scope.Sales = []
		var saleList = saleListStr.split(',');
		//$scope.val = saleList

		//Each Item in Trade list is: "<User1ID>:<TokenWillingID>|<TokenWillingColour>-<ColourRequired>""
		//Creating saleList Structure where each item is:
		//["Index", "SellerID", "SellerTokenID", "SellerTokenColour", "SellerTokenValue"]
		for ( var i = 0 ; i < saleList.length; i++ ) {
					var saleItem = []
					//Pushing Index
					saleItem.push(i+1);
					var tempStr = saleList[i].split(':');
					//Pushing Username
					saleItem.push(tempStr[0]);
					var tempStr2 = tempStr[1].split('|');
					//Pushing TokenID
					saleItem.push(tempStr2[0]);
					var tempStr3 = tempStr2[1].split('-');
					//Pushing Token Colour
					saleItem.push(tempStr3[0]);
					//Pushing Token Value
					saleItem.push(tempStr3[1]);
					$scope.Sales.push(saleItem);
				}
		//$scope.val = $scope.Sales;
	}, 
	//if Query doesn't get a response
	function(queryResponse) { // optional
	        $scope.query = "Could not connect to Blockchain"
	});

	$scope.makeTrade = function(trade) {
		//trade is of type: [sno, "Sellerid", "SellerTokenId" "SellerSellingColour", "BuyerColour"]
		//$scope.assets is an Array of type: 
		//{"id":"t1","colour":"red","user":"rishabh","sell":true,"value":15,"trade":true,"tradecolour":"green","salestring":"rishabh:t1|red-15","tradestring":"rishabh:t1|red-green","tradeStatus":"Exchange: green","sellStatus":"Value: 15 pt"}
		var found = false;
		var foundIndex;
		
		for( var i = 0 ; i < $scope.assets.length; i++ ) {

			if (trade[4] == $scope.assets[i].colour) {
				found = true;
				foundIndex = i;
				break;
			}
		}
		if (found == true) {
			var args = [$scope.username, $scope.assets[foundIndex].id, $scope.assets[foundIndex].colour, trade[1], trade[2], trade[3]];
			var tradeQuery = createJsonString("invoke", "trade", args);
			//$scope.trade = tradeQuery;
			$http({
			  url: 'http://9.232.130.202:3000/chaincode',
			  method: 'POST',
			  data: tradeQuery
			})
			//if Query Gets response successfully
			.then(function(queryResponse) {
				$location.path("/view6");
			}, 
			//if Query doesn't get a response
			function(queryResponse) { // optional
			        $scope.query = "Could not connect to Blockchain"
			});
		}
		else {
			$scope.tradeStatusString = "Error! You do not have the tokens required";
			//Give error as you do not have the token required
		}

	}

	$scope.makeSale = function(sale) {
		//trade is of type: [sno, "Sellerid", "SellerTokenId" "SellerSellingColour", "SellerValue"]
		//$scope.assets is an Array of type: 
		//{"id":"t1","colour":"red","user":"rishabh","sell":true,"value":15,"trade":true,"tradecolour":"green","salestring":"rishabh:t1|red-15","tradestring":"rishabh:t1|red-green","tradeStatus":"Exchange: green","sellStatus":"Value: 15 pt"}
		var found = false;

		var userPoints = parseInt($scope.points, 10);
		var sellVal = parseInt(sale[4], 10);
		if (userPoints >= sellVal) {
			found = true;
		}
		if (found == true) {
			//$scope.val = "mhfg";
			var args = [$scope.username, sale[1], sale[2], sale[3], sale[4]];
			$scope.val = args;
			var saleQuery = createJsonString("invoke", "buy", args);
			
			//$scope.trade = tradeQuery;
			$http({
			  url: 'http://9.232.130.202:3000/chaincode',
			  method: 'POST',
			  data: saleQuery
			})
			//if Query Gets response successfully
			.then(function(queryResponse) {
				$location.path("/view6");
			}, 
			//if Query doesn't get a response
			function(queryResponse) { // optional
			        $scope.query = "Could not connect to Blockchain"
			});
		}
		else {
			$scope.saleStatusString = "Error! You do not have the adequate points";
		}

	}

});