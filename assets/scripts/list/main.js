(function() {

  $(function() {
    var $postList = $('#post-list');
    for (var i = 0, len = window.app.listData.posts.length; i < len; ++i) {
      var postInfo = window.app.listData.posts[i];
      var listItem = new window.app.ListItem(postInfo);
      $postList.append(listItem.element());
    }
  });

})();
