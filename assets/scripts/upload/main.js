(function() {

  var $uploadText;
  var dropdown;

  function updateLanguageDropdown() {
    var lang = window.app.languageForText($uploadText.val());
    console.log('lang is', lang);
    dropdown.setLanguage(lang);
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
    });
    $('#submit-button').click(function() {
      window.app.createPost(dropdown.getLanguage(), $uploadText.val());
    });
  });

})();
