(function($) {
	$(window).on('load', function() {
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
				var fr = reader = new FileReader();
				fr.onload = function(evt) {
					var wb = XLSX.read(evt.target.result);

					var $ul = $('<ul class="nav nav-tabs">');
					$.each(wb.SheetNames, function(i, n) {
						var $a = $('<a class="nav-link" data-bs-toggle="tab">').attr('href', '#sh_'+n).text(n);
						if (i == 0) {
							$a.addClass('active');
						}
						$ul.append($('<li class="nav-item">').append($a));
					});

					var $tc = $('<div class="tab-content table-responsive my-4">');
					$.each(wb.SheetNames, function(i, n) {
						var $d = $('<div class="tab-pane">').attr('id', 'sh_'+n);
						if (i == 0) {
							$d.addClass('active');
						}
						$d.html(XLSX.utils.sheet_to_html(wb.Sheets[n]));
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
	});
})(jQuery);
