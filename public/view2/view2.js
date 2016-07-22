'use strict';

angular.module('myApp.view2', ['ngRoute', 'myApp'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view2', {
    templateUrl: 'view2/view2.html',
    controller: 'View2Ctrl'
  });
}])

.controller('View2Ctrl', function($scope, productService, $http, $location) {
	$scope.username = productService.getProducts()[0];
	if ($scope.username == null) {
		$location.path("/view1")
	}
	//$scope.username = "rishabh";


	function createQueryString(a) {
		var retString = {
		  "jsonrpc": "2.0",
		  "method": "query",
		  "params": {
		    "type": 1,
		    "chaincodeID":{
		        "name": "mycc"
		    },
		    "ctorMsg": {
		        "function":"query",
		        "args": [a]
		    },
		    "secureContext": "lukas"
		  },
		  "id": 1
		}
		return retString;
	}

	function createInvokeString(funcStr, args) {
		var retString = {
		  "jsonrpc": "2.0",
		  "method": "invoke",
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



	$scope.goToClaimToken = function() {
		$location.path("/view3")
	};

	var sendQuery = createQueryString($scope.username);

	$http({
	  url: 'http://9.232.130.202:3000/chaincode',
	  method: 'POST',
	  data: sendQuery
	})
	//if Query Gets response successfully
	.then(function(queryResponse) {
		$scope.assets = []
		var userQueryResult = angular.fromJson(queryResponse.data.result.message);
		$scope.name = userQueryResult.name;
		$scope.points = userQueryResult.points;
		$scope.tokenValue = userQueryResult.tokenvalue;
		$scope.totalPoints = (parseInt($scope.points, 10) + parseInt($scope.tokenValue, 10));
		if (userQueryResult.assets == "") {
			$scope.assets = [{"id": "N/A", "colour": "N/A"}]
			$scope.assetString = "None"
		}
		else {
			var assetsStr = userQueryResult.assets;
			var assetsList = assetsStr.split(',');
			for ( var i = 0 ; i < assetsList.length; i++ ) {
				var queryString = createQueryString(assetsList[i]);
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

	$scope.createTrade = function(asset) {
		//$scope.name = asset;
		var args = [asset.user, asset.id, asset.tradecolour];
		var funcStr = "setTradeStatus";
		var invokeStr = createInvokeString(funcStr, args);
		$http({
		  url: 'http://9.232.130.202:3000/chaincode',
		  method: 'POST',
		  data: invokeStr
		})
		//if Query Gets response successfully
		.then(function(queryResponse) {
			$location.path("/view6");
		}, 
		//if Query doesn't get a response
		function(queryResponse) { // optional
		        $scope.query = "Could not connect to blockchain"
		});
	}

	$scope.createSale = function(asset) {
		//$scope.name = asset;
		var args = [asset.user, asset.id, asset.value];
		var funcStr = "setSellStatus";
		var invokeStr = createInvokeString(funcStr, args);
		$http({
		  url: 'http://9.232.130.202:3000/chaincode',
		  method: 'POST',
		  data: invokeStr
		})
		//if Query Gets response successfully
		.then(function(queryResponse) {
			$location.path("/view6")
			//$state.reload;
		}, 
		//if Query doesn't get a response
		function(queryResponse) { // optional
		        $scope.query = "Could not connect to blockchain"
		});
	}

});