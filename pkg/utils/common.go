package utils

// we trim these characters when user right-clicks on a word without selecting it
const Punctuation = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~“”،؛؟۔￼"

// when double-click in QTextBrowser. some punctuations next to words
// are also selected, specially non-ascii ones,
// so we trim them on right-click -> Query action or on middle-click action
const QueryForceTrimChars = "‘’،؛"
