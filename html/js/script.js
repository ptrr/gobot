$(function(){
	var resizeTimer = null;
	$(window).bind('resize', function() {
		if (resizeTimer) clearTimeout(resizeTimer);
		resizeTimer = setTimeout(rsz, 100);
	});

	addChannelBox();
	doDisconnect();
	doConnect();
	rsz();
});

function addChannelBox(){
	$('#plus').click(function(){
		$('#channels').append('Channel: <input type="text" name="channel[]" /><br />');
	});
}

function doConnect(){
	var name = "";
	$('.connect').click(function(){
		name = $(this).attr("id");
		name = name.replace("but_", "");
		//window.location = '/init?name='+name;
		alert(name + " Connect");
		$.ajax({
			url: "/init?name="+name,
  			context: this,
  			success: function(){
				alert("Connected");
    			$(this).hide();
				$("#"+name+" button.disconnect").show();
				$("#"+name+" span").text("Connected");
				$("#"+name+" span").removeAttr('style');
				$("#"+name+" span").css('color', '#6add4b');
  			}
		});
	});
}

function doDisconnect(){
	var name = "";
	$('.disconnect').click(function(){
		name = $(this).attr("id");
		name = name.replace("but_", "");
		//window.location = '/init?name='+name;
		alert(name + " Disconnect");
		$.ajax({
			url: "/kill?name="+name,
  			context: this,
  			success: function(){
				alert("Disconnected");
				$(this).hide();
				$("#"+name+" button.connect").show();
				$("#"+name+" span").text("Not connected");
				$("#"+name+" span").removeAttr('style');
				$("#"+name+" span").css('color', '#dd4b4b');
  			}
		});
	});
}

function rsz(){
	$("#controls").height($(window).height())
	$("header").width($(window).width() - $("#controls").width()-11);
}

 