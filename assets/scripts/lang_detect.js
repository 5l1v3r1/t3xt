(function() {

  var CHAR_CLASS_LETTER = 0;
  var CHAR_CLASS_NUMBER = 1;
  var CHAR_CLASS_SPACE = 2;
  var CHAR_CLASS_SYMBOL = 3;

  function languageForText(str) {
    var freqs = normalizedWordFrequencies(str);
    var searchNode = window.app.codeIdentificationTree.TreeRoot;
    while (!searchNode.Leaf) {
      var keyword = searchNode.Keyword;
      if (freqs[keyword] > searchNode.Threshold) {
        searchNode = searchNode.TrueBranch;
      } else {
        searchNode = searchNode.FalseBranch;
      }
    }
    return {
      language: searchNode.LeafClassification,
      confidence: searchNode.LeafConfidence
    };
  }

  function normalizedWordFrequencies(str) {
    var res = {};
    var words = [];
    var totalCount = 0;

    var maps = [heterogeneousWords(str), homogeneousWords(str)];
    for (var i = 0, len = maps.length; i < len; ++i) {
      var map = maps[i];
      var keys = Object.keys(map);
      for (var j = 0, len1 = keys.length; j < len1; ++j) {
        var key = keys[j];
        var value = map[key];
        res[key] = value;
        words.push(key);
        totalCount += value;
      }
    }

    for (var i = 0, len = words.length; i < len; ++i) {
      var word = words[i];
      res[word] /= totalCount;
    }

    return res;
  }

  function heterogeneousWords(str) {
    var spaceTokens = str.split(/\s+/);
    var res = {};
    for (var i = 0, len = spaceTokens.length; i < len; ++i) {
      var word = spaceTokens[i];
      if (!isHomogeneous(word)) {
        if (res.hasOwnProperty(word)) {
          ++res[word];
        } else {
          res[word] = 1;
        }
      }
    }
    return res;
  }

  function homogeneousWords(str) {
    var currentClass = -1;
    var currentWord = '';
    var res = {};
    for (var i = 0, len = str.length; i < len; ++i) {
      var chr = str.charCodeAt(i);
      var charClass = characterClass(chr);
      if (charClass === currentClass) {
        currentWord += str[i];
      } else {
        if (currentWord.length > 0 && currentClass !== CHAR_CLASS_SPACE) {
          if (res.hasOwnProperty(currentWord)) {
            ++res[currentWord];
          } else {
            res[currentWord] = 1;
          }
        }
        currentWord = str[i];
        currentClass = charClass;
      }
    }
    if (currentClass !== CHAR_CLASS_SPACE && currentWord.length > 0) {
      if (res.hasOwnProperty(currentWord)) {
        ++res[currentWord];
      } else {
        res[currentWord] = 1;
      }
    }
    return res;
  }

  function isHomogeneous(word) {
    if (word.length === 0) {
      return true;
    }
    var chClass = characterClass(word.charCodeAt(0));
    for (var i = 1, len = word.length; i < len; ++i) {
      if (characterClass(word.charCodeAt(i)) != chClass) {
        return false;
      }
    }
    return true;
  }

  var LOWER_LETTER_START = 'a'.charCodeAt(0);
  var LOWER_LETTER_END = 'z'.charCodeAt(0);
  var UPPER_LETTER_START = 'A'.charCodeAt(0);
  var UPPER_LETTER_END = 'Z'.charCodeAt(0);
  var NUMBERS_START = '0'.charCodeAt(0);
  var NUMBERS_END = '9'.charCodeAt(0);

  var SPACE_CHARS = [' '.charCodeAt(0), '\t'.charCodeAt(0), '\n'.charCodeAt(0),
    '\r'.charCodeAt(0)];

  function characterClass(chr) {
    if ((chr >= LOWER_LETTER_START && chr <= LOWER_LETTER_END) ||
        (chr >= UPPER_LETTER_START && chr <= UPPER_LETTER_START)) {
      return CHAR_CLASS_LETTER;
    } else if (SPACE_CHARS.indexOf(chr) >= 0) {
      return CHAR_CLASS_SPACE;
    } else if (chr >= NUMBERS_START && chr <= NUMBERS_END) {
      return CHAR_CLASS_NUMBER;
    } else {
      return CHAR_CLASS_SYMBOL;
    }
  }

  window.app.languageForText = languageForText;

})();
