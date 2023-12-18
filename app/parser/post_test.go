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

            <i class="tgme_widget_message_bubble_tail">
                <svg class="bubble_icon" width="9px" height="20px" viewBox="0 0 9 20">
                    <g fill="none">
                        <path class="background" fill="#ffffff"
                              d="M8,1 L9,1 L9,20 L8,20 L8,18 C7.807,15.161 7.124,12.233 5.950,9.218 C5.046,6.893 3.504,4.733 1.325,2.738 L1.325,2.738 C0.917,2.365 0.89,1.732 1.263,1.325 C1.452,1.118 1.72,1 2,1 L8,1 Z"></path>
                        <path class="border_1x" fill="#d7e3ec"
                              d="M9,1 L2,1 C1.72,1 1.452,1.118 1.263,1.325 C0.89,1.732 0.917,2.365 1.325,2.738 C3.504,4.733 5.046,6.893 5.95,9.218 C7.124,12.233 7.807,15.161 8,18 L8,20 L9,20 L9,1 Z M2,0 L9,0 L9,20 L7,20 L7,20 L7.002,18.068 C6.816,15.333 6.156,12.504 5.018,9.58 C4.172,7.406 2.72,5.371 0.649,3.475 C-0.165,2.729 -0.221,1.464 0.525,0.649 C0.904,0.236 1.439,0 2,0 Z"></path>
                        <path class="border_2x"
                              d="M9,1 L2,1 C1.72,1 1.452,1.118 1.263,1.325 C0.89,1.732 0.917,2.365 1.325,2.738 C3.504,4.733 5.046,6.893 5.95,9.218 C7.124,12.233 7.807,15.161 8,18 L8,20 L9,20 L9,1 Z M2,0.5 L9,0.5 L9,20 L7.5,20 L7.5,20 L7.501,18.034 C7.312,15.247 6.64,12.369 5.484,9.399 C4.609,7.15 3.112,5.052 0.987,3.106 C0.376,2.547 0.334,1.598 0.894,0.987 C1.178,0.677 1.579,0.5 2,0.5 Z"></path>
                        <path class="border_3x"
                              d="M9,1 L2,1 C1.72,1 1.452,1.118 1.263,1.325 C0.89,1.732 0.917,2.365 1.325,2.738 C3.504,4.733 5.046,6.893 5.95,9.218 C7.124,12.233 7.807,15.161 8,18 L8,20 L9,20 L9,1 Z M2,0.667 L9,0.667 L9,20 L7.667,20 L7.667,20 L7.668,18.023 C7.477,15.218 6.802,12.324 5.64,9.338 C4.755,7.064 3.243,4.946 1.1,2.983 C0.557,2.486 0.52,1.643 1.017,1.1 C1.269,0.824 1.626,0.667 2,0.667 Z"></path>
                    </g>
                </svg>
            </i>
            <div class="tgme_widget_message_author accent_color"><a class="tgme_widget_message_owner_name"
                                                                    href="https://t.me/telegram"><span
                    dir="auto">Test Channel</span></a></div>


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

			<div class="tgme_widget_message_user"><a href="https://t.me/dukaliti"><i
				class="tgme_widget_message_user_photo bgcolor0" data-content="н"><img
				src="https://cdn4.cdn-telegram.org/file/img3.jpg"></i></a>
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
