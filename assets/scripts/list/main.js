(function() {

  $(function() {
    var $postList = $('#post-list');
    for (var i = 0, len = window.app.postList.length; i < len; ++i) {
      var postInfo = window.app.postList[i];
      var listItem = new window.app.ListItem(postInfo);
      $postList.append(listItem.element());
    }
  });

})();
