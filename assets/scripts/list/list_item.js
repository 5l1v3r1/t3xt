(function() {

  function ListItem(info) {
    this._$element = $('<li class="list-item"></li>');
    this._$link = $('<a class="list-link"></a>');
    this._$link.attr('href', '/view/' + info.secretId);
    this._$element.append(this._$link);

    this._$info = $('<div class="list-item-info"></div>');
    this._$lineCount = $('<label class="line-count"></label>');
    this._$dateLabel = $('<label class="date-label"></label>');
    this._$codeBlock = $('<pre class="preview-block"></pre>');

    this._$info.append(this._$lineCount);
    this._$info.append(this._$dateLabel);
    this._$link.append(this._$info);
    this._$link.append(this._$codeBlock);

    this._$lineCount.text(info.lines + ' line' + (info.lines !== 1 ? 's' : ''));

    this._$dateLabel.text(window.app.formatTime(info.postTime));
    this._$codeBlock.text(info.head);
  }

  ListItem.prototype.element = function() {
    return this._$element;
  };

  window.app.ListItem = ListItem;

})();
