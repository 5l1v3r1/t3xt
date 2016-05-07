(function() {

  var dropdown;
  var classifier;

  function main() {
    initializeClassifier();

    dropdown = new window.app.LangDropdown();

    var $uploadText = $('#upload-text');
    $uploadText.on('input propertychange', function() {
      if (dropdown.manuallySet()) {
        $uploadText.off('input propertychange');
        return;
      }
      classifier.textChanged($uploadText.val());
    });

    $('#submit-button').click(function() {
      window.app.createPost(dropdown.getLanguage(), $uploadText.val());
    });
  }

  function initializeClassifier() {
    classifier = new window.app.Classifier();
    classifier.onClassify = function(lang) {
      if (!dropdown.manuallySet()) {
        dropdown.setLanguage(lang);
      }
    };
  }

  $(main);

})();
