package application

import (
	"github.com/ilius/ayandict/pkg/frequency"
	"github.com/ilius/qt/widgets"
)

const (
	expanding = widgets.QSizePolicy__Expanding
)

// we trim these characters when user right-clicks on a word without selecting it
const punctuation = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~،؛؟۔"

// when double-click in QTextBrowser. some punctuations next to words
// are also selected, specially non-ascii ones,
// so we trim them on right-click -> Query action or on middle-click action
const queryForceTrimChars = "‘’،؛"

var frequencyTable *frequency.FrequencyTable
