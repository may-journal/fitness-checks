package main

import _ "embed"

// fileConfig is the check's keys in .fitnessrc.json, which --describe lists.
type fileConfig struct {
	// ProseBudget holds this check's hard caps. Each limit falls back to its
	// default when zero.
	ProseBudget budget `json:"proseBudget"`
	// Prose holds this check's settings for each kind of text.
	Prose prose `json:"prose"`
}

// budget holds the hard caps a repo sets under proseBudget. Each limit
// falls back to its default when zero.
type budget struct {
	// MaxSentenceWords is the most words one sentence may hold.
	MaxSentenceWords int `json:"maxSentenceWords"`
	// MaxParagraphSentences is the most sentences one paragraph may hold.
	MaxParagraphSentences int `json:"maxParagraphSentences"`
	// MaxSectionParagraphs is the most paragraphs one section may hold.
	MaxSectionParagraphs int `json:"maxSectionParagraphs"`
	// MaxListItemWords is the most words one list item may hold.
	MaxListItemWords int `json:"maxListItemWords"`
	// MaxListItems is the most items one list may hold.
	MaxListItems int `json:"maxListItems"`
	// MaxWords is the most prose words one section may hold.
	MaxWords int `json:"maxWords"`
}

// prose holds the settings a repo sets under prose, for each kind of text.
type prose struct {
	// Comments holds the `overlong` thresholds for code comments.
	Comments struct {
		// MaxWordsPerToken is the most comment words a function body may
		// hold for each code token before `overlong` fails it.
		MaxWordsPerToken int `json:"maxWordsPerToken"`
		// MinBodyComments is the fewest comments a function body needs
		// before `overlong` judges it.
		MinBodyComments int `json:"minBodyComments"`
	} `json:"comments"`
}

// defaults are the tight technical-doc limits a repo gets when it sets none.
var defaults = func() fileConfig {
	c := fileConfig{ProseBudget: budget{
		MaxSentenceWords:      23,
		MaxParagraphSentences: 4,
		MaxSectionParagraphs:  3,
		MaxListItemWords:      23,
		MaxListItems:          8,
		MaxWords:              300,
	}}
	c.Prose.Comments.MaxWordsPerToken = 3
	c.Prose.Comments.MinBodyComments = 2
	return c
}()

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte
