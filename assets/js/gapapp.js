/*
  The MIT License (MIT)

  Copyright (c) 2015 Charles Liu

  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:

  The above copyright notice and this permission notice shall be included in
  all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
  THE SOFTWARE.
*/

/*
* This file is used internally for my gap.go project, this file requires zepto
*/


function updateSessionCnt() {
	$.ajax({
		type: 'GET',
		url: '/sessionCnt',
		
		dataType: 'json',
		timeout: 300,
		context: $('body'),
		
		success: function(data){
			$("#sessionCnt").text(data.sessionCnt)
		},
		error: function(xhr, type){
			console.log("getting session count error")
		}
	});
}

$("#main").on('click',function(){
	var passcode = $("#passcode").val();
	var username = $("#username").val();
	$.ajax({
		type: 'GET',
		url: '/newSession?secret='+passcode+'&user='+username,
		
		dataType: 'json',
		timeout: 300,
		context: $('body'),
		
		success: function(data){
			var status = data.status;
			if (status === "ok") {
				var username = data.username;
				var password = data.password;
				
				
				var msg = "<b>Success! </b> <br><p>PPTP username: <code>"+ username + "</code> <br>password: <code>"+ password +"</code></p>";
				$(".info").fadeOut(400,function(){
					$(".info").html(msg);
					$(".info").fadeIn(400);	
					$(".info").removeClass("info");
				});
				
				
				setTimeout(300,function(){moveDown(".main");});
			} else {
				$(".info").fadeOut(400,function(){
					$(".info").html("<b>Error:</b> " + "<code>" + status + "</code>");
					$(".info").fadeIn(400);	
				});
			}
			//alert(data.status);
		},
		error: function(xhr, type){
			console.log("getting session count error")
		}
	});
});

//

$("#passcode").on("focus", function(e){
	$(".info").fadeOut(400,function(){
		$(".info").html("You should know the passcode by heart.");
		$(".info").fadeIn(400);	
	});
});

$("#passcode").on("focusout", function(e){
	$(".info").fadeOut(400,function(){
		if ($("#passcode").val() === "") {
			$(".info").html("Maybe you could ask someone for the passcode.");
		} else {
			if ($("#username").val() === "") {
				$(".info").html("Well done. Now go fill your username in.");
			} else {
				$(".info").html("Well done. You seem all set now.");
			}
		}
		
		$(".info").fadeIn(400);	
	});
});

$("#username").on("focus", function(e){
	$(".info").fadeOut(400,function(){
		if ($("#passcode").val() === ""){
			$(".info").html("Why don't you fill the passcode in first? Just curious.");
		} else {
			$(".info").html("Username should be in the format of <code>atticus-finch</code>. Use your english name.");
		}
		
		$(".info").fadeIn(400);	
	});
});

isUsernameValid = function(){
	var username = $("#username").val();
	var re = /^[a-z]{2,}-[a-z]{2,}$/;
	var isValid = re.test(username);
	return isValid;
}

$("#username").on("focusout", function(e){
	$(".info").fadeOut(400,function(){
		if (isUsernameValid()){
			if ($("#passcode").val() === ""){
				$(".info").html("Great, now go fill the passcode in.");
			} else {
				$(".info").html("You seem all set!");
			}
			
		} else {
			$(".info").html("Your username does not seem right: it should be in the format of <code>firstname-lastname</code>.");
		}
		
		$(".info").fadeIn(400);	
	});
});

$("#main").on("mouseover", function(e){
	$(".info").fadeOut(400,function(){
		$(".info").html("To potential hackers: please do not exploit this service.");
		$(".info").fadeIn(400);	
	});
});

$("#main").on("mouseout", function(e){
	$(".info").fadeOut(400,function(){
		if ($("#username").val() === "" || $("#passcode").val() === ""){
			$(".info").html("Better go fill these values in, huh?");
		} else {
			
			$(".info").html("Hmmm... Why give up?");
		}
		
		$(".info").fadeIn(400);	
	});
});


onePageScroll(".main", {
   sectionContainer: "section",     // sectionContainer accepts any kind of selector in case you don't want to use section
   easing: "ease",                  // Easing options accepts the CSS3 easing animation such "ease", "linear", "ease-in",
                                    // "ease-out", "ease-in-out", or even cubic bezier value such as "cubic-bezier(0.175, 0.885, 0.420, 1.310)"
   animationTime: 1000,             // AnimationTime let you define how long each section takes to animate
   pagination: true,                // You can either show or hide the pagination. Toggle true for show, false for hide.
   updateURL: false,                // Toggle this true if you want the URL to be updated automatically when the user scroll to each page.
   beforeMove: function(index) {},  // This option accepts a callback function. The function will be called before the page moves.
   afterMove: function(index) {},   // This option accepts a callback function. The function will be called after the page moves.
   loop: false,                     // You can have the page loop back to the top/bottom when the user navigates at up/down on the first/last page.
   keyboard: true,                  // You can activate the keyboard controls
   responsiveFallback: false,        // You can fallback to normal page scroll by defining the width of the browser in which
                                    // you want the responsive fallback to be triggered. For example, set this to 600 and whenever
                                    // the browser's width is less than 600, the fallback will kick in.
   direction: "vertical" 	
});



setInterval(updateSessionCnt, 3000);