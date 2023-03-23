package stardict

// "github.com/ilius/go-stardict"

// func TestDic(t *testing.T) {
// 	root := path.Join("tmp", "dic")

// 	if dict, err := stardict.Open(root); err == nil {
// 		for _, d := range dict {
// 			t.Logf("%s(%d)", d.GetBookName(), d.GetWordCount())
// 		}
// 	} else {
// 		t.Fatal(err)
// 	}

// 	// init dictionary with path to dictionary files and name of dictionary
// 	dict, err := stardict.NewDictionary(
// 		path.Join(root, "stardict-cdict-gb-2.4.2"),
// 		"cdict-gb",
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for _, kw := range []string{"one", "two", "three"} {
// 		senses := dict.Translate(kw) // get translations
// 		for i, seq := range senses { // for each translation analyze returned parts
// 			t.Logf("Sense %d\n", i+1)
// 			for j, p := range seq.Parts { // write each part contents to user
// 				t.Logf("Part %d:\n%c\n%s\n", j+1, p.Type, p.Data)
// 			}
// 		}

// 		if len(senses) == 0 {
// 			t.Log("Nothing found :(")
// 		}
// 	}
// }
