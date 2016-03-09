(function() {

  window.app.formatTime = function(epochTime) {
    var date = new Date(0);
    date.setUTCSeconds(epochTime);

    var monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'];
    var monthName = monthNames[date.getMonth()];
    var dateInYear = monthName + ' ' + date.getDate();
    if (date.getFullYear() !== new Date().getFullYear()) {
      return dateInYear + ', ' + date.getFullYear();
    } else {
      return dateInYear;
    }
  };

})();
