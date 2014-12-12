function displayError(error) {
	$('div.col-md').prepend('<div class="alert alert-danger alert-dismissible fade in" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><strong>Ooops!</strong> '+error+'</div>');
}