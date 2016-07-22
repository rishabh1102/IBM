'use strict';

angular.module('myApp.view1', ['ngRoute','myApp'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view1', {
    templateUrl: 'view1/view1.html',
    controller: 'View1Ctrl'
  });
}])

.controller('View1Ctrl', function($scope, $http, $location, productService) {

	$scope.createQuery = function(functionString) {
		var returnString = ""
		if (functionString == "checkLoginDetails") {
			returnString = {
			  "jsonrpc": "2.0",
			  "method": "query",
			  "params": {
			    "type": 1,
			    "chaincodeID":{
			        "name": "mycc"
			    },
			    "ctorMsg": {
			        "function":functionString,
			        "args": [$scope.username, $scope.password]
			    },
			    "secureContext": "lukas"
			  },
			  "id": 1
			}	
		} else {
			returnString = {
			  "jsonrpc": "2.0",
			  "method": "query",
			  "params": {
			    "type": 1,
			    "chaincodeID":{
			        "name": "mycc"
			    },
			    "ctorMsg": {
			        "function":functionString,
			        "args": [$scope.username]
			    },
			    "secureContext": "lukas"
			  },
			  "id": 1
			}
		}
		return returnString;
	}

	$scope.registerUser = function() {

		$scope.query = $scope.createQuery("entryExist")

		$scope.invoke = {
		  "jsonrpc": "2.0",
		  "method": "invoke",
		  "params": {
		    "type": 1,
		    "chaincodeID":{
		        "name": "mycc"
		    },
		    "ctorMsg": {
		        "function":"createUser",
		        "args":[$scope.username, $scope.password, $scope.name]
		    },
		    "secureContext": "lukas"
		  },
		  "id": 1
		}

		//Http Post request
	    $http({
		  url: 'http://9.232.130.202:3000/chaincode',
		  method: 'POST',
		  data: $scope.query
		})
		//if Query Gets response successfully
		.then(function(queryResponse) {

				//Json queryResponse. negUserAvail is negation of User Availability
	            var negUserAvail = angular.fromJson(queryResponse.data.result).message
	            //Username is available
	            if (negUserAvail == "FALSE") {
	             	$http({
					  url: 'http://9.232.130.202:3000/chaincode',
					  method: 'POST',
					  data: $scope.invoke
					})
					//if Invoke gets successful Response
					.then(function(invokeResponse) {
						//Query to check if account was created
						//var xyz = $location.path("/view2")
						productService.addProduct($scope.username)
						$location.path("/view6")
					},
					//if Invoke doesn't connect successfully
					function(invokeResponse) {
						$scope.deploy = "Could not connect to blockchain"
					});
					
	            } else {
	            	$scope.deploy = "This username cannot be taken, Please choose new one and try again"
	            }
	    }, 
	    //if Query doesn't get a response
	    function(queryResponse) { // optional
	            $scope.deploy = "Could not connect to blockchain"
	    });
	}

	$scope.loginUser = function() {
		$scope.query = $scope.createQuery("checkLoginDetails")
		//Http Post request
	    $http({
		  url: 'http://9.232.130.202:3000/chaincode',
		  method: 'POST',
		  data: $scope.query
		})
		//if Query Gets response successfully
		.then(function(queryResponse) {
				//Json queryResponse. negUserAvail is negation of User Availability
	            var loginSuccess = angular.fromJson(queryResponse.data.result).message
	            //Token is available
	            if (loginSuccess == "TRUE") {
	             	productService.addProduct($scope.username)
	             	$location.path("/view2")
	            } else {
	            	$scope.deploy = "Username or Password is incorrect!"
	            }
	    }, 
	    //if Query doesn't get a response
	    function(queryResponse) { // optional
	            $scope.deploy = "Could not connect to blockchain"
	    });
	}
});