(function() {

  $(function() {
    var $postList = $('#post-list');
    for (var i = 0, len = window.app.listData.posts.length; i < len; ++i) {
      var postInfo = window.app.listData.posts[i];
      var listItem = new window.app.ListItem(postInfo);
      $postList.append(listItem.element());
    }
    var lastId = window.app.listData.posts[0].id;
    var firstId = window.app.listData.posts[window.app.listData.posts.length-1].id;
    if (window.app.listData.hasNext) {
      $('.button-next').removeClass('page-button-disabled').click(function() {
        window.location = '/list?after=' + (lastId + 1);
      });
    }
    if (window.app.listData.hasLast) {
      $('.button-prev').removeClass('page-button-disabled').click(function() {
        window.location = '/list?before=' + (firstId - 1);
      });
    }
  });

})();
