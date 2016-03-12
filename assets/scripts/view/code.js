(function() {

  function CodeView() {
    this._$code = $('#code');
    this._lineSet = parseURLLineSet();
    this._generateCodeRows();
  }

  CodeView.prototype._generateCodeRows = function() {
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
    }
  };

  CodeView.prototype._registerRowClick = function(index, $row) {
    $row.click(function() {
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
