(function() {

  function Metadata() {
    this._$language = $('#post-language');
    this._$date = $('#post-date');

    this._$date.text(window.app.formatTime(window.app.postInfo.postTime));
    this._$language.text(window.app.postInfo.language);
  }

  window.app.Metadata = Metadata;

})();
