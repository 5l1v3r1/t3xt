(function() {

  // CLASSIFY_DELAY specifies how many milliseconds to
  // wait before classifying code after a user types.
  var CLASSIFY_DELAY = 500;
  var POST_URL = '/classify';

  function Classifier() {
    this._currentText = null;
    this._isFirst = true;

    this._fetching = false;
    this._changedWhileFetching = false;
    this._nextRequestTimeout = null;

    this.onClassify = null;
  }

  Classifier.prototype.textChanged = function(t) {
    this._currentText = t;
    if (this._fetching) {
      this._changedWhileFetching = true;
      return;
    } else if (this._isFirst) {
      // We will make the first request immediately,
      // since it is common to open the homepage and
      // immediately paste in a large block of code.
      this._isFirst = false;
      this._fetch();
      return;
    }

    if (this._nextRequestTimeout) {
      clearTimeout(this._nextRequestTimeout);
    }
    this._nextRequestTimeout = setTimeout(this._fetch.bind(this), CLASSIFY_DELAY);
  };

  Classifier.prototype._fetch = function() {
    this._fetching = true;
    this._nextRequestTimeout = null;

    requestClassification(this._currentText, function(data) {
      this._fetching = false;
      if (this._changedWhileFetching) {
        this._changedWhileFetching = false;
        this.textChanged(this._currentText);
      }
      this.onClassify(data);
    }.bind(this));
  };

  function requestClassification(text, callback) {
    if (text === '') {
      callback('Plain Text');
    } else {
      $.post(POST_URL, text, callback);
    }
  }

  window.app.Classifier = Classifier;

})();
