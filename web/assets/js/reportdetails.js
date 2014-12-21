$(document).ready(function(e) {
	$('a.report-details-finish').click(function() {
		$.getJSON('/report/'+$(this).attr('report')+'/edit?command=finish', function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
	
	$('a.report-player-paid').click(function() {
		$.getJSON('/report/'+$(this).attr('report')+'/edit?command=playerpaid', 'playerName='+$(this).attr('player')+'', function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
});