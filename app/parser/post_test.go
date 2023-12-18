package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

const testPostHTML = `
<html>
<body>
<div class="tgme_widget_message_wrap js-widget_message_wrap">
    <div class="tgme_widget_message text_not_supported_wrap js-widget_message" data-post="telegram/1"
         data-view="eyJjIjotMTA5Njg3MzUxMiwicCI6IjQzMTZnIiwidCI6MTcwMjgyMjEyOCwiaCI6ImNjNzg5ODgxY2ZkMTIzNTY4NiJ9">
        <div class="tgme_widget_message_user"><a href="https://t.me/telegram"><i
                class="tgme_widget_message_user_photo bgcolor0" data-content="н"><img
                src="https://cdn4.cdn-telegram.org/file/img1.jpg"></i></a>
        </div>
        <div class="tgme_widget_message_bubble">
            <div class="tgme_widget_message_author accent_color">
				<a class="tgme_widget_message_owner_name" href="https://t.me/telegram">
					<span dir="auto">Test Channel</span>
				</a>
			</div>
            <div class="tgme_widget_message_grouped_wrap js-message_grouped_wrap" data-margin-w="2"
                 data-margin-h="2" style="width:453px;">
                <div class="tgme_widget_message_grouped js-message_grouped" style="padding-top:84.989%">
                    <div class="tgme_widget_message_grouped_layer js-message_grouped_layer"
                         style="width:453px;height:385px">
                        <a class="tgme_widget_message_photo_wrap grouped_media_wrap blured js-message_photo"
                           style="left:0px;top:0px;width:134px;height:201px;margin-right:2px;margin-bottom:2px;background-image:url('https://cdn4.cdn-telegram.org/file/img2.jpg')"
                           data-ratio="0.66589861751152" href="https://t.me/telegram/1">
                            <div class="grouped_media_helper" style="top:0;bottom:0;left:98px;right:98px;">
                                <div class="tgme_widget_message_photo grouped_media"
                                     style="left:0;right:0;top:-1px;bottom:0px;"></div>
                            </div>
                        </a>
                    </div>
                </div>
            </div>
            <div class="tgme_widget_message_text js-message_text" dir="auto">
                <div class="tgme_widget_message_text js-message_text" dir="auto">Test text
                </div>
            </div>
			<div class="tgme_widget_message_user">
				<a href="https://t.me/dukaliti">
					<i class="tgme_widget_message_user_photo bgcolor0" data-content="н">
						<img src="https://cdn4.cdn-telegram.org/file/img3.jpg">
					</i>
				</a>
			</div>
        </div>
    </div>
    <div class="tgme_widget_message_roundvideo_player js-message_roundvideo_player">
        <div class="tgme_widget_message_roundvideo_wrap">
            <i class="tgme_widget_message_roundvideo_thumb"
               style="background-image:url('https://cdn4.cdn-telegram.org/file/123')"></i>
            <video class="tgme_widget_message_roundvideo js-message_roundvideo"
                   src="https://cdn4.cdn-telegram.org/file/video1.mp4"
                   width="100%" height="100%" preload muted autoplay loop playsinline></video>
            <div class="tgme_widget_message_roundvideo_muted"></div>
            <div class="tgme_widget_message_roundvideo_border"></div>
        </div>
    </div>
	<div class="tgme_widget_message_footer compact js-message_footer">
		<div class="tgme_widget_message_info short js-message_info">
			<span class="tgme_widget_message_views">28.6K</span><span class="copyonly"> views</span><span
				class="tgme_widget_message_meta"><a class="tgme_widget_message_date"
													href="https://t.me/dukaliti/4335"><time
				datetime="2023-12-15T16:29:45+00:00" class="time">20:29</time></a></span>
		</div>
	</div>
</div>
</body>
</html>
`

func getSelection() *goquery.Selection {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(testPostHTML))
	if err != nil {
		panic(err)
	}
	return doc.Find("body")
}

func getEmptySelection() *goquery.Selection {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(``))
	if err != nil {
		panic(err)
	}
	return doc.Find("body")
}

func TestGetImages(t *testing.T) {
	s := getSelection()
	images := GetImages(s)
	assert.Equal(t, 3, len(images))
	assert.Equal(t, "https://cdn4.cdn-telegram.org/file/img2.jpg", images[0])
	assert.Equal(t, "https://cdn4.cdn-telegram.org/file/img1.jpg", images[1])
	assert.Equal(t, "https://cdn4.cdn-telegram.org/file/img3.jpg", images[2])

	// Empty post
	s = getEmptySelection()
	images = GetImages(s)
	assert.Equal(t, 0, len(images))
}

func TestGetVideos(t *testing.T) {
	s := getSelection()
	videos := GetVideos(s)
	assert.Equal(t, 1, len(videos))
	assert.Equal(t, "https://cdn4.cdn-telegram.org/file/video1.mp4", videos[0])

	// Empty post
	s = getEmptySelection()
	videos = GetVideos(s)
	assert.Equal(t, 0, len(videos))
}

func TestGetCreated(t *testing.T) {
	s := getSelection()
	created := GetPostCreated(s)
	assert.Equal(t, "15 Dec 23 16:29 +0000", created.Format(time.RFC822Z))

	// Empty post should return time.Now()
	s = getEmptySelection()
	created = GetPostCreated(s)
	assert.Equal(t, time.Now().Format(time.DateTime), created.Format(time.DateTime))
}

func TestGetPostLink(t *testing.T) {
	s := getSelection()
	link := GetPostLink(s)
	assert.Equal(t, "https://t.me/s/telegram/1", link)

	// Empty post should return ""
	s = getEmptySelection()
	link = GetPostLink(s)
	assert.Equal(t, "", link)
}

func TestGetPostTextHTML(t *testing.T) {
	s := getSelection()
	text := GetPostTextHTML(s)
	assert.Equal(t, "<p>Test text</p>", text)

	// Empty post should return ""
	s = getEmptySelection()
	text = GetPostTextHTML(s)
	assert.Equal(t, "", text)
}

func TestGetPostTitle(t *testing.T) {
	tbl := []struct {
		inp string
		out string
	}{
		{"", ""},
		{"<p> Test text </p>", "Test text"},
		{"<p> <a href=\"https://t.me/s/telegram\">telegram</a> 👍 👍</p>", "telegram 👍 👍"},
		{" title <br/>not-title", "title"},
		{" title <br>not-title", "title"},
	}
	for _, tb := range tbl {
		title := GetPostTitle(tb.inp)
		assert.Equal(t, tb.out, title)
	}
}
