package parser

import "github.com/PuerkitoBio/goquery"

// FixEmoji removes all emoji tags
func FixEmoji(s *goquery.Selection) *goquery.Selection {
	// Clone the selection to avoid modifying the original
	s = s.Clone()

	// replace <i class="emoji">emoji</i> with just `emoji`
	s.Find("i.emoji").Each(func(_ int, s *goquery.Selection) {
		s.ReplaceWithHtml(s.Text())
	})
	// replace <tg-emoji>emoji</tg-emoji> with just `emoji`
	s.Find("tg-emoji").Each(func(_ int, s *goquery.Selection) {
		s.ReplaceWithHtml(s.Text())
	})
	return s
}

// FixLinks removes all attributes from links
func FixLinks(s *goquery.Selection) *goquery.Selection {
	// Clone the selection to avoid modifying the original
	s = s.Clone()

	// Clean up links from attributes: target, rel, onclick
	s.Find("a").Each(func(_ int, s *goquery.Selection) {
		s.RemoveAttr("target")
		s.RemoveAttr("rel")
		s.RemoveAttr("onclick")
	})
	return s
}

// RemoveUnsafeTags removes all tags except <a>, <i>, <b>, <br>
func RemoveUnsafeTags(s *goquery.Selection) *goquery.Selection {
	// Clone the selection to avoid modifying the original
	s = s.Clone()

	// Remove all tags except <a>, <i>, <b>, <br>
	s.Find("*").Each(func(i int, s *goquery.Selection) {
		if !s.Is("a, i, b, br") {
			s.ReplaceWithHtml(s.Text())
		}
	})
	return s
}
