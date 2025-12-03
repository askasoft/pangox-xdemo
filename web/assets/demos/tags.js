(function($) {
	function tags_label_align() {
		$('#tags_form').removeClass('label-left label-right').addClass($(this).val());
	}

	//----------------------------------------------------
	// init
	//
	function tags_init() {
		$('#tags_form').on('change', 'input[name=label]', tags_label_align);
	}

	$(window).on('load', tags_init);
})(jQuery);
