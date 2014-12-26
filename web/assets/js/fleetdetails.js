$(document).ready(function(e) {
	$(function() {
		$('#addMemberSelectMember').filterByText($('#addMemberSelectMemberSearch'), true);
	});
	
	$('a[data-toggle=collapse]').click(function() {
		$(this).toggleClass('active');
	});
	
	$('a.fleet-details-toggle').click(function() {
		$('div.fleet-details').toggle();
	});
	
	$('a.fleet-details-save').click(function() {
		var formData = $('#fleetDetailsForm').serializeArray();
		formData.push({ name: "command", value: "editDetails" });
		
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: formData,
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
			url: '/fleet/'+$(this).attr('fleet')
		});
	});
	
	$('a.fleet-details-calculate').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=calculatePayouts",
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
			url: '/fleet/'+$(this).attr('fleet')
		});
	});
	
	$('a.fleet-details-finish').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: "command=finishFleet",
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
			url: '/fleet/'+$(this).attr('fleet')
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