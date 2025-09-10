(function($) {
	var page = 1;
	if (location.hash) {
		var ps = $.queryParams(location.hash.substring(1));
		if (ps['page'] && parseInt(ps['page']) > 0) {
			page = parseInt(ps['page']);
		}
	}

	$(window).on('load', function() {
		if (page > 1) {
			var $a = $('#xlsx_preview .nav-tabs a').eq(page - 1);
			if ($a.length) {
				new bootstrap.Tab($a.get(0)).show();
			}
		}
	});
})(jQuery);
