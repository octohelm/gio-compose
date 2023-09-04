package modifier

import (
	"github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/text"
)

func ProvideTextStyle(modifiers ...modifier.Modifier[*text.Style]) modifier.Modifier[any] {
	return &textStyleProvider{modifiers: modifiers}
}

type textStyleProvider struct {
	modifier.Discord

	target    compose.VNode
	modifiers []modifier.Modifier[*text.Style]
}

func (t textStyleProvider) Wrap(node compose.VNode) compose.Component {
	t.target = node
	return t
}

func (t textStyleProvider) Build(c compose.BuildContext) compose.VNode {
	if t.target == nil {
		return compose.ProvideTextStyle(t.modifiers...).Children(
			c.ChildVNodes()...,
		)
	}

	return compose.ProvideTextStyle(t.modifiers...).Children(
		t.target.Children(c.ChildVNodes()...),
	)
}

func TextStyle(style text.Style) modifier.Modifier[*text.Style] {
	return &textStyleModifier{style: style}
}

type textStyleModifier struct {
	style text.Style
}

func (t *textStyleModifier) Modify(s *text.Style) {
	s.SetStyle(t.style)
}
