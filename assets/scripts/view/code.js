(function() {

  var SCROLL_TOP_TIMEOUT = 10;
  var SCROLL_TOP_SPACE = 30;
  var LINE_MARKER_PADDING = 10;

  function CodeView() {
    this._$code = $('#code');
    this._$firstSelectedRow = null;
    this._lineSet = parseURLLineSet();
    this._generateCodeRows();
    this._evenOutNumberWidths();

    if (this._$firstSelectedRow !== null) {
      setTimeout(function() {
        $(window).scrollTop(this._$firstSelectedRow.offset().top - SCROLL_TOP_SPACE);
      }.bind(this), SCROLL_TOP_TIMEOUT);
    }
  }

  CodeView.prototype._generateCodeRows = function() {
    var firstSelected = this._lineSet.first();
    var lines = window.app.postInfo.content.split('\n');
    for (var i = 0, len = lines.length; i < len; ++i) {
      var line = lines[i];
      var $row = $('<tr class="code-line"><td class="code-line-marker"></td>' +
        '<td class="code-text-container selectable"><br class="code-line-break">');
      $row.children('.code-line-marker').attr({'line-number': i+1});

      var $codeText = $('<pre class="code-text selectable"></pre>').text(line);
      $row.children('.code-text-container').prepend($codeText);
      this._$code.append($row);

      if (this._lineSet.includes(i+1)) {
        $row.addClass('highlighted-line');
      }
      this._registerRowClick(i+1, $row);

      if (i === firstSelected) {
        this._$firstSelectedRow = $row;
      }
    }
  };

  CodeView.prototype._evenOutNumberWidths = function() {
    var markers = $('.code-line-marker');
    var width = 0;
    for (var i = 0, len = markers.length; i < len; ++i) {
      width = Math.max(width, $(markers[i]).width());
    }
    for (var i = 0, len = markers.length; i < len; ++i) {
      var m = $(markers[i]);
      var w = m.width();
      var difference = width - w;
      if (difference > 0) {
        m.css({paddingLeft: LINE_MARKER_PADDING+difference});
      }
    }
  };

  CodeView.prototype._registerRowClick = function(index, $row) {
    $row.click(function(e) {
      if (this._lineSet.includes(index)) {
        this._lineSet.remove(index);
        $row.removeClass('highlighted-line');
      } else {
        this._lineSet.add(index);
        $row.addClass('highlighted-line');
      }
      this._lineSetChanged();
    }.bind(this));
  };

  CodeView.prototype._lineSetChanged = function() {
    window.location.hash = '#' + this._lineSet;
  };

  function parseURLLineSet() {
    if (window.location.hash !== '') {
      try {
        var lineCount = window.app.postInfo.content.split('\n').length;
        return window.app.LineSet.parse(window.location.hash.substr(1), lineCount);
      } catch (e) {
      }
    }
    return new window.app.LineSet({});
  }

  window.app.CodeView = CodeView;

})();
