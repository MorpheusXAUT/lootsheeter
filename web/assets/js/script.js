$('a[data-toggle=collapse]').click(function() {
	if ($(this).hasClass('active') === true) {
		$(this).removeClass('active');
	} else {
		$(this).addClass('active');
	}
});