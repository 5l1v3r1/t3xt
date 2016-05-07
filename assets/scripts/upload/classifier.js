(function() {

  var MIN_REQUEST_PERIOD = 3000;
  var POST_URL = '/classify';

  function Classifier() {
    this._currentText = null;

    this._fetching = false;
    this._changedWhileFetching = false;
    this._lastRequestTime = null;
    this._nextRequestTimeout = null;

    this.onClassify = null;
  }

  Classifier.prototype.textChanged = function(t) {
    this._currentText = t;
    if (this._fetching) {
      this._changedWhileFetching = true;
      return;
    } else if (this._nextRequestTimeout !== null) {
      return;
    } else if (this._lastRequestTime === null) {
      this._fetch();
      return;
    }

    var timeSince = Math.max(0, new Date().getTime()-this._lastRequestTime);
    if (timeSince >= MIN_REQUEST_PERIOD) {
      this._fetch();
    } else {
      var timeout = MIN_REQUEST_PERIOD - timeSince;
      this._nextRequestTimeout = setTimeout(this._fetch.bind(this), timeout);
    }
  };

  Classifier.prototype._fetch = function() {
    this._fetching = true;
    this._nextRequestTimeout = null;
    this._lastRequestTime = new Date().getTime();

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
