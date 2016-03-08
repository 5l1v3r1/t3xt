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

    this._$lineCount.text(info.lines + ' line' + (info.lineCount !== 1 ? 's' : ''));

    var date = new Date(0);
    date.setUTCSeconds(info.postTime);
    this._$dateLabel.text(date.toString());
    this._$codeBlock.text(info.head);
  }

  ListItem.prototype.element = function() {
    return this._$element;
  };

  window.app.ListItem = ListItem;

})();
