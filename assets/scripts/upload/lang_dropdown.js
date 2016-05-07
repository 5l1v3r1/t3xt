(function() {

  function LangDropdown() {
    this._dropdown = new window.dropdownjs.Dropdown(200, '#f0f0f0');
    this._$element = $(this._dropdown.element());
    this._$element[0].id = 'lang-select';

    this._dropdown.setOptions(languageNames(), 0);

    this._manuallySet = false;
    this._dropdown.onChange = function() {
      this._manuallySet = true;
    }.bind(this);

    $('#content').prepend(this._$element);
  }

  LangDropdown.prototype.getLanguage = function() {
    return this._dropdown.getValue();
  };

  LangDropdown.prototype.setLanguage = function(lang) {
    this._dropdown.setValue(lang);
  };

  LangDropdown.prototype.manuallySet = function() {
    return this._manuallySet;
  };

  function languageNames() {
    var res = window.app.languageNames;
    var idx = res.indexOf('Plain Text');
    if (idx >= 0) {
      res.splice(idx, 1);
    }
    res.sort();
    res.splice(0, 0, 'Plain Text');
    return res;
  }

  window.app.LangDropdown = LangDropdown;

})();
