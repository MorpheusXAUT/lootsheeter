function toggleFleetmemberEdit(index, enable) {
	if (enable) {
		$('#form_fleetmember_role_'.concat(index)).addClass('hide');
		$('#form_fleetmember_role_edit_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_sitemodifier_'.concat(index)).addClass('hide');
		$('#form_fleetmember_sitemodifier_edit_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_paymentmodifier_'.concat(index)).addClass('hide');
		$('#form_fleetmember_paymentmodifier_edit_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_payout_'.concat(index)).addClass('hide');
		$('#form_fleetmember_payout_edit_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_payoutcomplete_'.concat(index)).addClass('hide');
		$('#form_fleetmember_payoutcomplete_edit_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_action_'.concat(index)).addClass('hide');
		$('#form_fleetmember_action_edit_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_action_cancel_'.concat(index)).addClass('hide');
		$('#form_fleetmember_action_cancel_edit_'.concat(index)).removeClass('hide');
	} else {
		$('#form_fleetmember_role_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_role_edit_'.concat(index)).addClass('hide');
		$('#form_fleetmember_sitemodifier_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_sitemodifier_edit_'.concat(index)).addClass('hide');
		$('#form_fleetmember_paymentmodifier_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_paymentmodifier_edit_'.concat(index)).addClass('hide');
		$('#form_fleetmember_payout_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_payout_edit_'.concat(index)).addClass('hide');
		$('#form_fleetmember_payoutcomplete_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_payoutcomplete_edit_'.concat(index)).addClass('hide');
		$('#form_fleetmember_action_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_action_edit_'.concat(index)).addClass('hide');
		$('#form_fleetmember_action_cancel_'.concat(index)).removeClass('hide');
		$('#form_fleetmember_action_cancel_edit_'.concat(index)).addClass('hide');
	}
}