// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

// Event is the alert event.
type Event struct {
	ID     string  `json:"id"`
	Index  *Index  `json:"index"`
	Metric *Metric `json:"metric"`
	Rule   *Rule   `json:"rule"`
}

// NewEvent returns a new event from metric and index.
func NewEvent(m *Metric, idx *Index, rule *Rule) *Event {
	ev := &Event{Metric: m, Index: idx, Rule: rule}
	ev.generateID()
	return ev
}

// generateID generates a sha1 string id for the event.
func (ev *Event) generateID() {
	slug := fmt.Sprintf("%s:%d", ev.Metric.Name, ev.Metric.Stamp)
	hash := sha1.New()
	hash.Write([]byte(slug))
	ev.ID = hex.EncodeToString(hash.Sum(nil))
}

// EventWrapper is a wrapper of Event for tmp usage.
type EventWrapper struct {
	*Event
	Project               *Project `json:"project"`
	User                  *User    `json:"user"`
	RuleTranslatedComment string   `json:"ruleTranslatedComment"`
}

// NewWrapperOfEvent creates an event wrapper from given event.
func NewWrapperOfEvent(ev *Event) *EventWrapper {
	return &EventWrapper{Event: ev}
}

// TranslateRuleComment translates rule comment variables with metric name and
// rule pattern.
//
//	m := &Metric{Name: "timer.count_ps.foo"}
//	r := &Rule{Pattern: "timer.count_ps.*", Comment: "$1 timing"}
//	ev := &Event{Metric:m, Rule:r}
//	ev.TranslateRuleComment()  // ev.RuleTranslatedComment => "foo timing"
//
func (ew *EventWrapper) TranslateRuleComment() {
	patternParts := strings.Split(ew.Rule.Pattern, ".")
	metricParts := strings.Split(ew.Metric.Name, ".")
	if len(patternParts) != len(metricParts) { // Unexcepted input metric and pattern.
		ew.RuleTranslatedComment = ew.Rule.Comment // Use original comment
		return
	}
	i := 0
	s := ew.Rule.Comment
	for j, patternPart := range patternParts {
		if patternPart == "*" {
			i++
			repl := fmt.Sprintf("$%d", i)
			s = strings.Replace(s, repl, metricParts[j], 1)
		}
	}
	ew.RuleTranslatedComment = s
}
