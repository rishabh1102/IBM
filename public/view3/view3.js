'use strict';

angular.module('myApp.view3', ['ngRoute','myApp'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view3', {
    templateUrl: 'view3/view3.html',
    controller: 'View3Ctrl'
  });
}])

.controller('View3Ctrl', function($scope, productService, $http, $location) {
	$scope.claimButtonText = "Claim Token!"
	var input = document.getElementById('camsource');
    input.addEventListener('change', handleFiles, false);

	$scope.username = productService.getProducts()[0];
	if ($scope.username == null) {
		$location.path("/view1")
	}
	//$scope.username = "rishabh1102"

	$scope.goToProfile = function() {
		$location.path("/view2")
	};

	$scope.claimToken = function() {
		var qrResult = document.getElementById('qr-value')
		$scope.Token = qrResult.innerText
		$scope.query = {
		  "jsonrpc": "2.0",
		  "method": "query",
		  "params": {
		    "type": 1,
		    "chaincodeID":{
		        "name": "mycc"
		    },
		    "ctorMsg": {
		        "function":"queryTokenOwnership",
		        "args": [$scope.Token]
		    },
		    "secureContext": "lukas"
		  },
		  "id": 1
		}

		$scope.invoke = {
		  "jsonrpc": "2.0",
		  "method": "invoke",
		  "params": {
		    "type": 1,
		    "chaincodeID":{
		        "name": "mycc"
		    },
		    "ctorMsg": {
		        "function":"claimToken",
		        "args":[$scope.Token, $scope.username]
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
	            var tokenAvail = angular.fromJson(queryResponse.data.result).message
	            //Token is available
	            if (tokenAvail == "TRUE") {
	             	$http({
					  url: 'http://9.232.130.202:3000/chaincode',
					  method: 'POST',
					  data: $scope.invoke
					})
					//if Invoke gets successful Response
					.then(function(invokeResponse) {
						//Query to check if account was created
						$location.path("/view6")
						$scope.claimTokenStatus = "Success! Token has been claimed."
						$scope.claimButtonText = "Claim Another Token!"
					},
					//if Invoke doesn't connect successfully
					function(invokeResponse) {
						$scope.claimTokenStatus = "Could not connect to blockchain"
					});
					
	            } else {
	            	$scope.claimTokenStatus = "Token is unavailable (Does not exist or already taken)"
	            }
	    }, 
	    //if Query doesn't get a response
	    function(queryResponse) { // optional
	            $scope.claimTokenStatus = "Could not connect to blockchain"
	    });

	}

	function handleFiles(e) {
	    var canvas = document.getElementById('qr-canvas')
	    var ctx = canvas.getContext('2d');
	    var url = URL.createObjectURL(e.target.files[0]);
	    var img = new Image();
	    img.src = url;
	    img.onload = function() {
	    ctx.drawImage(img, 0, 0, 300, 300);  
	    	try { qrcode.decode();  }
	    	catch(err) { $("#qr-value").text(err); } 
	    }
	}

	

}); 