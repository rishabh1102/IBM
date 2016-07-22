'use strict';

angular.module('myApp.view5', ['ngRoute', 'myApp'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view5', {
    templateUrl: 'view5/view5.html',
    controller: 'View5Ctrl'
  });
}])

.controller('View5Ctrl', function($scope, productService, $http, $location) {
	$scope.username = productService.getProducts()[0];
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

	var sendQuery = createJsonString("query", "getLeaderBoard", ["userList"]);

	$http({
	  url: 'http://9.232.130.202:3000/chaincode',
	  method: 'POST',
	  data: sendQuery
	})
	//if Query Gets response successfully
	.then(function(queryResponse) {
		var leaderBoardList = queryResponse.data.result.message.split(",")
		//$scope.userPoints = leaderBoardList
		$scope.userPoints = []
		for ( var i = 0 ; i < leaderBoardList.length; i++ ) {
			var temp = [];
			temp.push(i+1)
			var tempStr = leaderBoardList[i].split("-");
			temp.push(tempStr[0]);
			var tempStr1 = tempStr[1].split("|");
			temp.push(tempStr1[0]);
			temp.push(tempStr1[1]);
			temp.push(parseInt(tempStr1[0], 10) + parseInt(tempStr1[1], 10));
			$scope.userPoints.push(temp)
		}
		$scope.userPoints.sort(function(a,b){return a[4]<b[4];})
	}, 
	//if Query doesn't get a response
	function(queryResponse) { // optional
	        $scope.query = "Could not connect to blockchain"
	});

	$scope.getValue = function(a) {
		return a[4];
	}
});