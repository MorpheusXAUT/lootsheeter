$(document).ready(function(e) {
	$('a.report-details-toggle').click(function() {
		$('div.report-details').toggle();
	});
	
	$('a.report-details-save').click(function() {
		$.getJSON('/report/'+$(this).attr('report')+'/edit?command=editdetails', $('#reportDetailsForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
});