(function() {

  var $uploadText;
  var dropdown;

  function updateLanguageDropdown() {
    var lang = window.app.languageForText($uploadText.val());
    dropdown.setLanguage(lang.language);
  }

  $(function() {
    $uploadText = $('#upload-text');
    dropdown = new window.app.LangDropdown();
    $uploadText.on('input propertychange', function() {
      if (dropdown.manuallySet()) {
        $uploadText.off('input propertychange');
        return;
      }
      updateLanguageDropdown();
    }.bind(this));
  });

})();
