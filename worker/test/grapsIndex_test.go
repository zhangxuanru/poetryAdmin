package test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Category"
	"poetryAdmin/worker/core/grasp/poetry/Content"
	"poetryAdmin/worker/core/parse"
	"regexp"
	"testing"
	"time"
)

//单元测试--抓取首页
func TestGrabsIndex(t *testing.T) {
	var SubReceiveMsg parse.SubscribeMsg
	Convey("测试抓取所有", t, func() {
		SubReceiveMsg = parse.SubscribeMsg{
			PubTile:  "抓取所有",
			AddDate:  time.Now().Unix(),
			Status:   define.TaskStatusImplemented,
			TaskMark: define.GrabPoetryAll,
		}
		parse.NewDispatch(SubReceiveMsg).Execution()
		time.Sleep(20 * time.Second)
	})
}

//测试诗文类型详情页
func TestCategory(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	home := &define.HomeFormat{
		Identifier: "test",
		Data: define.DataMap{
			1: &define.TextHrefFormat{
				Href:         "https://so.gushiwen.org/gushi/tangshi.aspx",
				Text:         "唐诗三百",
				ShowPosition: 1,
			},
		},
	}
	Category.NewCategory().GraspByIndexData(home)
	time.Sleep(120 * time.Second)
}

//测试诗文详情页
func TestContent(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	poetry := &define.PoetryAuthorList{
		AuthorName:      "柳宗元",
		PoetryTitle:     "江雪",
		PoetrySourceUrl: "/shiwenv_58313be2d918.aspx",
		GenreTitle:      "五言绝句",
		Category: &define.TextHrefFormat{
			Text:         "唐诗三百",
			Href:         "https://so.gushiwen.org/gushi/tangshi.aspx",
			ShowPosition: 1,
		},
	}
	Content.NewContent().GraspContentData(poetry)

	time.Sleep(70 * time.Second)

	//Content.NewContent().GraspContentSaveData("/shiwenv_73add8822103.aspx", nil)
}

func TestA(t *testing.T) {

	//去掉<div class="contyishang">和<div class="dingpai">的HTML内容
	str := `
<div class="contyishang">
<div style="height:30px; font-weight:bold; font-size:16px; margin-bottom:10px; clear:both;">
<h2><span style="float:left;">典故：不为五斗米折腰</span></h2>
<a style="float:left; margin-top:7px; margin-left:5px;" href="javascript:PlayZiliaoquan(1607)"><img id="speakerimgZiliaoquan1607" src="https://song.gushiwen.org/siteimg/speaker.png" alt="" width="16" height="16"/></a>
</div>
<p>　　中国古代有不少因维护人格，保持气节而不食的故事，<a href="https://so.gushiwen.org/authorv_07d17f8539d7.aspx" target="_blank">陶渊明</a>“不为五斗米折腰”就是其中最具代表性的一例。</p>
<p>　　东晋后期的大<a href="https://so.gushiwen.org/authorv_07d17f8539d7.aspx">诗</a>人陶渊明，是名人之后，他的曾祖父是赫赫有名的东晋大司马。年轻时的陶渊明本有“大济于苍生”之志，可是，在国家濒临崩溃的动乱年月里，陶渊明的一腔抱负根本无法实现。加之他性格耿直，清明廉正，不愿卑躬屈膝攀附权贵，因而和污浊黑暗的现实社会发生了尖锐的矛盾，产生了格格不入的感情。</p>
<p>　　为了生存，陶渊明最初做过州里的小官，可由于看不惯官场上的那一套恶劣作风，不久便辞职回家了。后来，为了生活他还陆续做过一些地位不高的官职，过着时隐时仕的生活。陶渊明最后一次做官，是义熙元年（405年）。那一年，已过“不惑之年”（四十一岁）的陶渊明在朋友的劝说下，再次出任彭泽县令。有一次，县里派督邮来了解情况。有人告诉陶渊明说：那是上面派下来的人，应当穿戴整齐、恭恭敬敬地去迎接。陶渊明听后长长叹了一口气：“我不愿为了小小县令的五斗薪俸，就低声下气去向这些家伙献殷勤。”说完，就辞掉官职，回家去了。陶渊明当彭泽县令，不过八十多天。他这次弃职而去，便永远脱离了官场。</p>
<p>　　此后，他一面读书为文，一面参加农业劳动。后来由于农田不断受灾，房屋又被火烧，家境越来越恶化。但他始终不愿再为官受禄，甚至连江州刺使送来的米和肉也坚拒不受。朝廷曾征召他任著作郎，也被他拒绝了。</p>
<p>　　陶渊明是在贫病交加中离开人世的。他原本可以活得舒适些，至少衣食不愁，但那要以付出人格和气节为代价。陶渊明因“不为五斗米折腰”，而获得了心灵的自由，获得了人格的尊严，写出了一代文风并流传百世的诗文。在为后人留下宝贵文学财富的同时，也留下了弥足珍贵的精神财富。他因“不为五斗米折腰”的高风亮节，成为中国后代有志之士的楷模。 <a title="收起" href="javascript:ziliaoClose(1607)">▲</a></p>
</div>
<div class="dingpai">
<a id="dingzl1607" href="javascript:dingzl('1607','https://so.gushiwen.org/author.aspx?id=645')">有用</a><a id="paizl1607" style=" margin-left:10px;" href="javascript:paizl('1607','https://so.gushiwen.org/author.aspx?id=645')">没用</a>
<a style="width:34px; height:18px; line-height:19px; margin-top:2px; float:right; color:#aeaeae;" href="/jiucuo.aspx?u=%e8%b5%84%e6%96%991607%e3%80%8a%e5%85%b8%e6%95%85%ef%bc%9a%e4%b8%8d%e4%b8%ba%e4%ba%94%e6%96%97%e7%b1%b3%e6%8a%98%e8%85%b0%e3%80%8b" target="_blank">完善</a>
</div>
<div class="cankao">
<p style=" color:#919090;margin:0px; font-size:12px;line-height:160%;">本节内容由匿名网友上传，原作者已无法考证。<a style=" color:#919090;" href="https://www.gushiwen.org/">本站</a>免费发布仅供学习参考，其观点不代表本站立场。站务邮箱：service@gushiwen.org</p>
</div>`

	mustCompile := regexp.MustCompile(`(?msU)<div class="contyishang">.*</div>`)
	s := mustCompile.ReplaceAllString(str, "")
	compile := regexp.MustCompile(`(?msU)<div class="dingpai">.*</div>`)
	s = compile.ReplaceAllString(s, "")
	nr := regexp.MustCompile(`(?m)[\r\n|\t]`)
	s = nr.ReplaceAllString(s, "")

	logrus.Infoln(s)

	return

	src := "https://song.gushiwen.org/authorImg/taoyuanming.jpg"
	//src := "https://song.gushiwen.org/machine/ziliao/1601/ok.mp3"
	fileName, err2 := data.NewUploadStore().Upload(src)
	logrus.Infoln(fileName)
	logrus.Infoln(err2)

	return
	file := "D:/server/gitData/goPath/poetryAdmin/worker/test/index.html"
	bytes, err := tools.ReadFile(file)
	logrus.Info("err:", err)
	query, e := tools.NewDocumentFromReader(string(bytes))
	logrus.Info("err:", e)

	query.Find(".right>.sons").Eq(2).Find(".cont>a").Each(func(j int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		logrus.Infoln("href:", href, "text:", selection.Text())
	})
}

func TestB(T *testing.T) {
	ch := make(chan bool, 5)
	end := make(chan bool)
	go func() {
		for {
			select {
			case <-ch:
				logrus.Infoln("ch.......")
			case <-end:
				if len(ch) > 0 {
					logrus.Infoln("还有数据.....")
					continue
				}
				time.Sleep(1 * time.Second)
				goto GoEnd
			}
		}
	GoEnd:
		logrus.Info("end......")
		return
	}()

	for i := 0; i < 10; i++ {
		go func(i int) {
			if i > 4 {
				end <- true
			}
		}(i)
		go func(i int) {
			ch <- true
		}(i)
	}
	time.Sleep(10 * time.Second)
}
