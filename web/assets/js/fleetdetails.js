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
		var formData = $('#addProfitForm').serializeArray();
		formData.push({ name: "command", value: "addProfit" });
		
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
	
	$('a.add-loss-submit').click(function() {
		var formData = $('#addLossForm').serializeArray();
		formData.push({ name: "command", value: "addLoss" });
		
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
	
	$('a.fleet-member-list-toggle').click(function() {
		$('div.fleet-member-list[member='+$(this).attr('member')+']').toggle();
	});
	
	$('a.fleet-member-list-save').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: $('form.fleet-member-list-form[member='+$(this).attr('member')+']').serialize(),
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
			url: '/fleet/'+$(this).attr('fleet') + '/members/'+$(this).attr('member')
		});
	});
	
	$('a.add-member-submit').click(function() {
		$.ajax({
			accepts: "application/json",
			cache: false,
			data: $('#addMemberForm').serialize(),
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
			type: "POST",
			url: '/fleet/'+$(this).attr('fleet') + '/members'
		});
	});
	
	$('a.fleet-member-list-remove').click(function() {		
		$.ajax({
			accepts: "application/json",
			cache: false,
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
			type: "DELETE",
			url: '/fleet/'+$(this).attr('fleet')+'/members/'+$(this).attr('member')
		});
	});
});