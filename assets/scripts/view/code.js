(function() {

  function CodeView() {
    this._$code = $('#code');
    this._generateCodeRows(parseURLLineSet());
  }

  CodeView.prototype._generateCodeRows = function(lineSet) {
    var lines = window.app.postInfo.content.split('\n');
    for (var i = 0, len = lines.length; i < len; ++i) {
      var line = lines[i];
      var $row = $('<tr class="code-line"><td class="code-line-marker"></td>' +
        '<td class="code-text-container selectable"><br class="code-line-break">');
      $row.children('.code-line-marker').attr({'line-number': i+1});

      var $codeText = $('<pre class="code-text selectable"></pre>').text(line);
      $row.children('.code-text-container').prepend($codeText);
      this._$code.append($row);

      if (lineSet !== null && lineSet.includes(i)) {
        // TODO: select the line here.
      }
    }
  };

  function parseURLLineSet() {
    if (window.location.hash !== '') {
      try {
        return window.app.LineSet.parse(window.location.hash.substr(1));
      } catch (e) {
      }
    }
    return null;
  }

  window.app.CodeView = CodeView;

})();
