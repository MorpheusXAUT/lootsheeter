$(document).ready(function(e) {
	$('a[data-toggle=collapse]').click(function() {
		$(this).toggleClass('active');
	});
	
	$('a.fleet-details-toggle').click(function() {
		$('div.fleet-details').toggle();
	});
	
	$('a.fleet-details-save').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=editdetails', $('#fleetDetailsForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
	
	$('a.fleet-details-calculate').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=calculate', function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
	
	$('a.fleet-details-finish').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=finish', function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
	
	$('a.add-profit-submit').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=addprofit', $('#addProfitForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
	
	$('a.add-loss-submit').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=addloss', $('#addLossForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
	
	$('a.fleet-member-list-toggle').click(function() {
		$('div.fleet-member-list[member='+$(this).attr('member')+']').toggle();
	});
	
	$('a.fleet-member-list-save').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=editmember', $('form.fleet-member-list-form[member='+$(this).attr('member')+']').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
	
	$('a.add-member-submit').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=addmember', $('#addMemberForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
	
	$('a.fleet-member-list-remove').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=removemember', 'removeMemberID='+$(this).attr('member')+'', function(data) {
			if (data.result === "success" && data.error === null) {
				location.reload(true);
			} else {
				displayError(data.error);
			}
		});
	});
});