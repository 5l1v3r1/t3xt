package main

import (
	"sort"

	"github.com/gopherjs/gopherjs/js"
	"github.com/unixpickle/whichlang/idtree"
	"github.com/unixpickle/whichlang/tokens"
)

var Classifier *idtree.Classifier

func main() {
	app := js.Global.Get("window").Get("app")

	idTreeJSON := app.Get("codeIdentificationTree")
	classifierText := js.Global.Get("JSON").Call("stringify", idTreeJSON).String()

	var err error
	Classifier, err = idtree.DecodeClassifier([]byte(classifierText))
	if err != nil {
		panic(err)
	}

	app.Set("languageForText", js.MakeFunc(Classify))
	app.Set("languageNames", languageNames())
}

func Classify(this *js.Object, args []*js.Object) interface{} {
	if len(args) != 1 {
		panic("expected one argument")
	}
	documentText := args[0].String()
	freqs := tokens.CountTokens(documentText).Freqs()
	return Classifier.Classify(freqs)
}

func languageNames() []string {
	seenLangs := map[string]bool{}
	searchSeenLangs(seenLangs, Classifier)
	delete(seenLangs, "Plain Text")

	langs := make([]string, 0, len(seenLangs))
	for lang := range seenLangs {
		langs = append(langs, lang)
	}
	sort.Strings(langs)
	langs = append([]string{"Plain Text"}, langs...)

	return langs
}

func searchSeenLangs(seen map[string]bool, c *idtree.Classifier) {
	if c.LeafClassification != nil {
		seen[*c.LeafClassification] = true
	} else {
		searchSeenLangs(seen, c.FalseBranch)
		searchSeenLangs(seen, c.TrueBranch)
	}
}
