(function() {

  function createPost(language, code) {
    var $form = $('<form method="POST" style="display: none"></form>');
    var $lang = $('<input type="hidden" name="language">').val(language);
    var $body = $('<input type="hidden" name="code">').val(code);
    $form.append($lang).append($body);
    $(document.body).append($form);
    $form.submit();
  }

  window.app.createPost = createPost;

})();
