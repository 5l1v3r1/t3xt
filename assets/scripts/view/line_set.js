(function() {

  function LineSet(map) {
    this._map = map;
  }

  LineSet.parse = function(str, maxLine) {
    var map = {};
    var parts = str.split(',');
    for (var i = 0, len = parts.length; i < len; ++i) {
      var part = parts[i];
      var rangeMatch = /^([0-9]*)-([0-9]*)$/.exec(part);
      if (rangeMatch !== null) {
        var rangeStart = parseInt(rangeMatch[1]);
        var rangeEnd = parseInt(rangeMatch[2]);
        if (rangeEnd > maxLine || isNaN(rangeEnd) || isNaN(rangeStart)) {
          throw new Error('invalid range: ' + part);
        }
        for (var j = rangeStart; j <= rangeEnd; ++j) {
          map[j] = true;
        }
      } else {
        var line = parseInt(part);
        if (isNaN(line) || line < 0 || line > maxLine) {
          throw new Error('invalid line number: ' + part);
        }
        map[line] = true;
      }
    }
    return new LineSet(map);
  };

  LineSet.prototype.toString = function() {
    var sortedLines = [];
    var keys = Object.keys(this._map);
    for (var i = 0, len = keys.length; i < len; ++i) {
      sortedLines[i] = parseInt(keys[i]);
    }
    sortedLines.sort(function(a, b) {
      if (a < b) {
        return -1;
      } else if (a === b) {
        return 0;
      }
      return 1;
    });
    var rangeStrs = [];
    var rangeStart = -1;
    var rangeEnd = -1;
    for (var i = 0, len = sortedLines.length; i <= len; ++i) {
      var line = (i === len ? 0 : sortedLines[i]);
      if (rangeEnd === -1) {
        rangeStart = line;
        rangeEnd = line;
      } else if (rangeEnd+1 === line) {
        ++rangeEnd;
      } else {
        if (rangeStart === rangeEnd) {
          rangeStrs.push(rangeStart);
        } else {
          rangeStrs.push(rangeStart + '-' + rangeEnd);
        }
        rangeStart = line;
        rangeEnd = line;
      }
    }
    return rangeStrs.join(',');
  };

  LineSet.prototype.includes = function(line) {
    return this._map[line];
  };

  LineSet.prototype.add = function(line) {
    this._map[line] = true;
  };

  LineSet.prototype.remove = function(line) {
    delete this._map[line];
  };

  window.app.LineSet = LineSet;

})();
