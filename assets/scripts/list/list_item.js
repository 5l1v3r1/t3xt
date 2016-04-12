(function() {

  function ListItem(info) {
    this._$element = $('<li class="list-item"></li>');
    this._$info = $('<div class="list-item-info"></div>');
    this._$lineCount = $('<label class="line-count"></label>');
    this._$dateLabel = $('<label class="date-label"></label>');
    this._$codeBlock = $('<pre class="preview-block"></pre>');

    this._$info.append(this._$lineCount);
    this._$info.append(this._$dateLabel);
    this._$element.append(this._$info);
    this._$element.append(this._$codeBlock);

    this._$lineCount.text(info.lines + ' line' + (info.lines !== 1 ? 's' : ''));

    this._$dateLabel.text(window.app.formatTime(info.postTime));
    this._$codeBlock.text(info.head);

    this._$element.click(this._view.bind(this, info.secretId));
  }

  ListItem.prototype.element = function() {
    return this._$element;
  };

  ListItem.prototype._view = function(shareId) {
    window.location = '/view/' + shareId;
  };

  window.app.ListItem = ListItem;

})();
