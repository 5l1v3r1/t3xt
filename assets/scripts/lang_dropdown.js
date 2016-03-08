(function() {

  function LangDropdown() {
    this._dropdown = new window.dropdownjs.Dropdown(200);
    this._$element = $(this._dropdown.element());

    this._dropdown.setOptions(allLanguages(), 0);

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

  function allLanguages() {
    var obj = {};
    languagesUnderNode(obj, window.app.codeIdentificationTree.TreeRoot);
    return Object.keys(obj);
  }

  function languagesUnderNode(res, node) {
    if (node.Leaf) {
      res[node.LeafClassification] = true;
    } else {
      languagesUnderNode(res, node.TrueBranch);
      languagesUnderNode(res, node.FalseBranch);
    }
  }

  window.app.LangDropdown = LangDropdown;

})();
