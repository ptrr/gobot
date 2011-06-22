$(function(){
	addChannelBox();
	doConnect();
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
		alert(name);
		$.ajax({
			url: "/init?name="+name,
  			context: this,
  			success: function(){
				$(this).removeClass("connect");
				$(this).addClass("disconnect");
    			$(this).text("Disconnect")
				$("#"+name+" span").text("Connected");
				$("#"+name+" span").removeAttr('style');
				$("#"+name+" span").css('color', 'green');
  			}
		});
	});
}