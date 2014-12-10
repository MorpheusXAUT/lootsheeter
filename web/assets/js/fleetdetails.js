$(document).ready(function(e) {
	$('a.fleet-details-toggle').click(function() {
		$('div.fleet-details').toggle();
	});
	
	$('a.fleet-details-save').click(function() {
		console.log($('#fleetDetails').serialize());
	});
		
	$('a.add-profit-submit').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=addprofit', $('#addProfitForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				console.log('success');
			} else {
				console.log(data);
			}
		});
	});
	
	$('a.add-loss-submit').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=addloss', $('#addLossForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				console.log('success');
			} else {
				console.log(data);
			}
		});
	});
	
	$('a.add-profit-submit').click(function() {
		$.getJSON('/fleet/'+$(this).attr('fleet')+'/edit?command=addprofit', $('#addProfitForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				console.log('success');
			} else {
				console.log(data);
			}
		});
	});
	
	$('a.fleet-member-list-toggle').click(function() {
		$('div.fleet-member-list[member='+$(this).attr('member')+']').toggle();
	});
	
	$('a.fleet-member-list-remove').click(function() {
		$('tr.fleet-member-list-row[member='+$(this).attr('member')+']').remove();
	});
	
	$('a.fleet-member-list-save').click(function() {
		console.log($('#fleetMemberList').serialize());
	});
	
	$('a.add-member-submit').click(function() {
		$.getJSON("/fleet/'+$(this).attr('fleet')+'/edit?command=addmember", $('#addMemberForm').serialize(), function(data) {
			if (data.result === "success" && data.error === null) {
				console.log('success');
			} else {
				console.log(data);
			}
		});
	});
});