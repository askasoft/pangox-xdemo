(function($) {
	var page = 1;
	if (location.hash) {
		var ps = $.queryParams(location.hash.substring(1));
		if (ps['page'] && parseInt(ps['page']) > 0) {
			page = parseInt(ps['page']);
		}
	}

	function xlx_preview() {
		$.ajax({
			url: $('#xls_dnload').attr('href'),
			xhr: function() {
				var xhr = new XMLHttpRequest();
				xhr.onreadystatechange = function() {
					if (xhr.readyState == 2) {
						if (xhr.status == 200) {
							xhr.responseType = "blob";
						} else {
							xhr.responseType = "text";
						}
					}
				};
				return xhr;
			},
			beforeSend: main.loadmask,
			success: function(data) {
				var fr = new FileReader();
				fr.onload = function(evt) {
					var wb = XLSX.read(evt.target.result, { cellDates: true });
					var cnt = wb.SheetNames.length;
					if (page > cnt) {
						page = 1;
					}

					var $ul = $('<ul class="nav nav-tabs">');
					wb.SheetNames.forEach(function(sn, i) {
						var $a = $('<a class="nav-link" data-bs-toggle="tab">').attr('href', '#sh_'+sn).text(sn);
						if (i == (page - 1)) {
							$a.addClass('active');
						}
						$ul.append($('<li class="nav-item">').append($a));
					});

					var $tc = $('<div class="tab-content table-responsive my-4">');
					wb.SheetNames.forEach(function(sn, i) {
						var $d = $('<div class="tab-pane">').attr('id', 'sh_'+sn);
						if (i == (page - 1)) {
							$d.addClass('active');
						}
						$d.html(XLSX.utils.sheet_to_html(wb.Sheets[sn], { cellDates: true }));
						$d.children('table').addClass('table table-bordered table-striped');
						$tc.append($d);
					});

					$('#xls_preview').append($ul, $tc);
				};
				fr.readAsArrayBuffer(data);
			},
			error: main.ajax_error,
			complete: main.unloadmask
		});
	}

	$(window).on('load', xlx_preview);
})(jQuery);
