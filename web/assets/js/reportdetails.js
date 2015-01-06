$(document).ready(function(e) {
	$('a.report-details-finish').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=finishReport",
			dataType: "json",
			error: displayAjaxError,
			success: function(reply) {
				if (reply.result === "success" && reply.error === null) {
					location.reload(true);
				} else {
					displayError(reply.error);
				}
			},
			timeout: 10000,
			type: "PUT",
			url: '/report/'+$(this).attr('report')
		});
	});
	
	$('a.report-player-paid').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=playerPaid&playerName="+$(this).attr('player'),
			dataType: "json",
			error: displayAjaxError,
			success: function(reply) {
				if (reply.result === "success" && reply.error === null) {
					location.reload(true);
				} else {
					displayError(reply.error);
				}
			},
			timeout: 10000,
			type: "PUT",
			url: '/report/'+$(this).attr('report')+'/players'
		});
	});
});