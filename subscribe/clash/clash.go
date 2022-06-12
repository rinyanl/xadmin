package clash

import (
	"context"
	"fmt"
	"io/ioutil"
	"xadmin/app/db"
	"xadmin/conf"

	"os"

	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Test []string `yaml:"array.test,flow"`
}
type Proxies struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}
type ProxyGroups struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"`
	Proxies  []string `yaml:"proxies"`
	Url      string   `yaml:"url,omitempty"`
	Interval int      `yaml:"interval,omitempty"`
}
type YamlConfig struct {
	Port               int           `yaml:"port"`
	Ipv6               bool          `yaml:"ipv6"`
	SocksPort          int           `yaml:"socks-port"`
	AllowLan           bool          `yaml:"allow-lan"`
	LogLevel           string        `yaml:"log-level"`
	ExternalController string        `yaml:"external-controller"`
	Proxies            []Proxies     `yaml:"proxies"`
	ProxyGroups        []ProxyGroups `yaml:"proxy-groups"`
	Rules              []string      `yaml:"rules"`
}

var Yf = YamlConfig{}

// 生成文件
func CreateVip(email string, vipType string, dir string) error {
	y := YamlConfig{
		Port:               7890,
		Ipv6:               true,
		SocksPort:          7891,
		AllowLan:           false,
		LogLevel:           "info",
		ExternalController: "127.0.0.1:9090",
		ProxyGroups: []ProxyGroups{
			{
				Name: "Proxies",
				Type: "select",
			},
			{
				Name: "Domestic",
				Type: "select",
				Proxies: []string{
					"DIRECT",
					"Proxies",
				},
			},
			{
				Name:     "Auto",
				Type:     "fallback",
				Url:      "http://www.gstatic.com/generate_204",
				Interval: 30000,
			},
		},
		Rules: ClashRules(),
	}
	switch vipType {
	case "vip1":
		// 引入 vip1 的 线路
		v1p, err := GetProxiesMap(email, vipType)
		if err != nil {
			return err
		}
		y.Proxies = v1p
		// 自动添加所有
		for _, k := range y.Proxies {
			y.ProxyGroups[0].Proxies = append(y.ProxyGroups[0].Proxies, k.Name)
			y.ProxyGroups[2].Proxies = append(y.ProxyGroups[2].Proxies, k.Name)
			// fmt.Println(y.ProxyGroups[0].Proxies)
		}
		// 赋值保存
		Yf = y
		err = SaveYamlFile(dir)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("ok")
}

//  信息
func GetProxiesMap(email string, vipType string) ([]Proxies, error) {

	switch vipType {
	case "vip1":

		v1 := []Proxies{}

		// 获取订阅配置
		sub, _, err := ClashRulelist()
		if err != nil {
			return v1, err
		}

		for _, v := range sub {
			puv := Proxies{
				Name:     v.Name,
				Type:     "trojan",
				Server:   v.Server,
				Port:     v.Port,
				Password: email,
			}
			v1 = append(v1, puv)
		}

		return v1, nil
	}

	ep := []Proxies{}
	return ep, nil
}

// 保存文件
func SaveYamlFile(dir string) error {
	data, _ := yaml.Marshal(&Yf)
	// fmt.Printf("%s\n", string(data))

	exPath, _ := os.Getwd()
	err := ioutil.WriteFile(exPath+"/app/assets/clash/"+dir+"/config.yaml", data, os.ModePerm)
	if err != nil {
		fmt.Printf("%v  错误", err)
		return err
	}
	return nil
}

// 读取文件内容、暂时没用
func ReadYaml() {

	exPath, _ := os.Getwd()
	file, err := os.Open(exPath + "/app/assets/clash/test.yaml")

	if err != nil {
		fmt.Println("打开文件错误", err)
		return
	}

	// 读取文件内容
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("读取失败")
		return
	}

	err = yaml.Unmarshal(data, &Yf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("student=%+v\n", Yf)
}

func ClashRules() []string {
	rules := []string{
		"DOMAIN,,DIRECT",
		"DOMAIN-SUFFIX,services.googleapis.cn, Proxies",
		"DOMAIN-SUFFIX,xn  ngstr lra8j.com, Proxies",
		"DOMAIN,safebrowsing.urlsec.qq.com,DIRECT",
		"DOMAIN,safebrowsing.googleapis.com,DIRECT",
		"DOMAIN,developer.apple.com, Proxies",
		"DOMAIN-SUFFIX,digicert.com, Proxies",
		"DOMAIN,ocsp.apple.com, Proxies",
		"DOMAIN,ocsp.comodoca.com, Proxies",
		"DOMAIN,ocsp.usertrust.com, Proxies",
		"DOMAIN,ocsp.sectigo.com, Proxies",
		"DOMAIN,ocsp.verisign.net, Proxies",
		"DOMAIN-SUFFIX,apple dns.net, Proxies",
		"DOMAIN,testflight.apple.com, Proxies",
		"DOMAIN,sandbox.itunes.apple.com, Proxies",
		"DOMAIN,itunes.apple.com, Proxies",
		"DOMAIN-SUFFIX,apps.apple.com, Proxies",
		"DOMAIN-SUFFIX,blobstore.apple.com, Proxies",
		"DOMAIN,cvws.icloud content.com, Proxies",
		"DOMAIN-SUFFIX,mzstatic.com,DIRECT",
		"DOMAIN-SUFFIX,itunes.apple.com,DIRECT",
		"DOMAIN-SUFFIX,icloud.com,DIRECT",
		"DOMAIN-SUFFIX,icloud content.com,DIRECT",
		"DOMAIN-SUFFIX,me.com,DIRECT",
		"DOMAIN-SUFFIX,aaplimg.com,DIRECT",
		"DOMAIN-SUFFIX,cdn20.com,DIRECT",
		"DOMAIN-SUFFIX,cdn apple.com,DIRECT",
		"DOMAIN-SUFFIX,akadns.net,DIRECT",
		"DOMAIN-SUFFIX,akamaiedge.net,DIRECT",
		"DOMAIN-SUFFIX,edgekey.net,DIRECT",
		"DOMAIN-SUFFIX,mwcloudcdn.com,DIRECT",
		"DOMAIN-SUFFIX,mwcname.com,DIRECT",
		"DOMAIN-SUFFIX,apple.com,DIRECT",
		"DOMAIN-SUFFIX,apple cloudkit.com,DIRECT",
		"DOMAIN-SUFFIX,apple mapkit.com,DIRECT",
		"DOMAIN-SUFFIX,cn,DIRECT",
		"DOMAIN-KEYWORD, cn,DIRECT",
		"DOMAIN-SUFFIX,126.com,DIRECT",
		"DOMAIN-SUFFIX,126.net,DIRECT",
		"DOMAIN-SUFFIX,127.net,DIRECT",
		"DOMAIN-SUFFIX,163.com,DIRECT",
		"DOMAIN-SUFFIX,360buyimg.com,DIRECT",
		"DOMAIN-SUFFIX,36kr.com,DIRECT",
		"DOMAIN-SUFFIX,acfun.tv,DIRECT",
		"DOMAIN-SUFFIX,air matters.com,DIRECT",
		"DOMAIN-SUFFIX,aixifan.com,DIRECT",
		"DOMAIN-KEYWORD,alicdn,DIRECT",
		"DOMAIN-KEYWORD,alipay,DIRECT",
		"DOMAIN-KEYWORD,taobao,DIRECT",
		"DOMAIN-SUFFIX,amap.com,DIRECT",
		"DOMAIN-SUFFIX,autonavi.com,DIRECT",
		"DOMAIN-KEYWORD,baidu,DIRECT",
		"DOMAIN-SUFFIX,bdimg.com,DIRECT",
		"DOMAIN-SUFFIX,bdstatic.com,DIRECT",
		"DOMAIN-SUFFIX,bilibili.com,DIRECT",
		"DOMAIN-SUFFIX,bilivideo.com,DIRECT",
		"DOMAIN-SUFFIX,caiyunapp.com,DIRECT",
		"DOMAIN-SUFFIX,clouddn.com,DIRECT",
		"DOMAIN-SUFFIX,cnbeta.com,DIRECT",
		"DOMAIN-SUFFIX,cnbetacdn.com,DIRECT",
		"DOMAIN-SUFFIX,cootekservice.com,DIRECT",
		"DOMAIN-SUFFIX,csdn.net,DIRECT",
		"DOMAIN-SUFFIX,ctrip.com,DIRECT",
		"DOMAIN-SUFFIX,dgtle.com,DIRECT",
		"DOMAIN-SUFFIX,dianping.com,DIRECT",
		"DOMAIN-SUFFIX,douban.com,DIRECT",
		"DOMAIN-SUFFIX,doubanio.com,DIRECT",
		"DOMAIN-SUFFIX,duokan.com,DIRECT",
		"DOMAIN-SUFFIX,easou.com,DIRECT",
		"DOMAIN-SUFFIX,ele.me,DIRECT",
		"DOMAIN-SUFFIX,feng.com,DIRECT",
		"DOMAIN-SUFFIX,fir.im,DIRECT",
		"DOMAIN-SUFFIX,frdic.com,DIRECT",
		"DOMAIN-SUFFIX,g cores.com,DIRECT",
		"DOMAIN-SUFFIX,godic.net,DIRECT",
		"DOMAIN-SUFFIX,gtimg.com,DIRECT",
		"DOMAIN,cdn.hockeyapp.net,DIRECT",
		"DOMAIN-SUFFIX,hongxiu.com,DIRECT",
		"DOMAIN-SUFFIX,hxcdn.net,DIRECT",
		"DOMAIN-SUFFIX,iciba.com,DIRECT",
		"DOMAIN-SUFFIX,ifeng.com,DIRECT",
		"DOMAIN-SUFFIX,ifengimg.com,DIRECT",
		"DOMAIN-SUFFIX,ipip.net,DIRECT",
		"DOMAIN-SUFFIX,iqiyi.com,DIRECT",
		"DOMAIN-SUFFIX,jd.com,DIRECT",
		"DOMAIN-SUFFIX,jianshu.com,DIRECT",
		"DOMAIN-SUFFIX,knewone.com,DIRECT",
		"DOMAIN-SUFFIX,le.com,DIRECT",
		"DOMAIN-SUFFIX,lecloud.com,DIRECT",
		"DOMAIN-SUFFIX,lemicp.com,DIRECT",
		"DOMAIN-SUFFIX,licdn.com,DIRECT",
		"DOMAIN-SUFFIX,luoo.net,DIRECT",
		"DOMAIN-SUFFIX,meituan.com,DIRECT",
		"DOMAIN-SUFFIX,meituan.net,DIRECT",
		"DOMAIN-SUFFIX,mi.com,DIRECT",
		"DOMAIN-SUFFIX,miaopai.com,DIRECT",
		"DOMAIN-SUFFIX,microsoft.com,DIRECT",
		"DOMAIN-SUFFIX,microsoftonline.com,DIRECT",
		"DOMAIN-SUFFIX,miui.com,DIRECT",
		"DOMAIN-SUFFIX,miwifi.com,DIRECT",
		"DOMAIN-SUFFIX,mob.com,DIRECT",
		"DOMAIN-SUFFIX,netease.com,DIRECT",
		"DOMAIN-SUFFIX,office.com,DIRECT",
		"DOMAIN-SUFFIX,office365.com,DIRECT",
		"DOMAIN-KEYWORD,officecdn,DIRECT",
		"DOMAIN-SUFFIX,oschina.net,DIRECT",
		"DOMAIN-SUFFIX,ppsimg.com,DIRECT",
		"DOMAIN-SUFFIX,pstatp.com,DIRECT",
		"DOMAIN-SUFFIX,qcloud.com,DIRECT",
		"DOMAIN-SUFFIX,qdaily.com,DIRECT",
		"DOMAIN-SUFFIX,qdmm.com,DIRECT",
		"DOMAIN-SUFFIX,qhimg.com,DIRECT",
		"DOMAIN-SUFFIX,qhres.com,DIRECT",
		"DOMAIN-SUFFIX,qidian.com,DIRECT",
		"DOMAIN-SUFFIX,qihucdn.com,DIRECT",
		"DOMAIN-SUFFIX,qiniu.com,DIRECT",
		"DOMAIN-SUFFIX,qiniucdn.com,DIRECT",
		"DOMAIN-SUFFIX,qiyipic.com,DIRECT",
		"DOMAIN-SUFFIX,qq.com,DIRECT",
		"DOMAIN-SUFFIX,qqurl.com,DIRECT",
		"DOMAIN-SUFFIX,rarbg.to,DIRECT",
		"DOMAIN-SUFFIX,ruguoapp.com,DIRECT",
		"DOMAIN-SUFFIX,segmentfault.com,DIRECT",
		"DOMAIN-SUFFIX,sinaapp.com,DIRECT",
		"DOMAIN-SUFFIX,smzdm.com,DIRECT",
		"DOMAIN-SUFFIX,snapdrop.net,DIRECT",
		"DOMAIN-SUFFIX,sogou.com,DIRECT",
		"DOMAIN-SUFFIX,sogoucdn.com,DIRECT",
		"DOMAIN-SUFFIX,sohu.com,DIRECT",
		"DOMAIN-SUFFIX,soku.com,DIRECT",
		"DOMAIN-SUFFIX,speedtest.net,DIRECT",
		"DOMAIN-SUFFIX,sspai.com,DIRECT",
		"DOMAIN-SUFFIX,suning.com,DIRECT",
		"DOMAIN-SUFFIX,taobao.com,DIRECT",
		"DOMAIN-SUFFIX,tencent.com,DIRECT",
		"DOMAIN-SUFFIX,tenpay.com,DIRECT",
		"DOMAIN-SUFFIX,tianyancha.com,DIRECT",
		"DOMAIN-SUFFIX,tmall.com,DIRECT",
		"DOMAIN-SUFFIX,tudou.com,DIRECT",
		"DOMAIN-SUFFIX,umetrip.com,DIRECT",
		"DOMAIN-SUFFIX,upaiyun.com,DIRECT",
		"DOMAIN-SUFFIX,upyun.com,DIRECT",
		"DOMAIN-SUFFIX,veryzhun.com,DIRECT",
		"DOMAIN-SUFFIX,weather.com,DIRECT",
		"DOMAIN-SUFFIX,weibo.com,DIRECT",
		"DOMAIN-SUFFIX,xiami.com,DIRECT",
		"DOMAIN-SUFFIX,xiami.net,DIRECT",
		"DOMAIN-SUFFIX,xiaomicp.com,DIRECT",
		"DOMAIN-SUFFIX,ximalaya.com,DIRECT",
		"DOMAIN-SUFFIX,xmcdn.com,DIRECT",
		"DOMAIN-SUFFIX,xunlei.com,DIRECT",
		"DOMAIN-SUFFIX,yhd.com,DIRECT",
		"DOMAIN-SUFFIX,yihaodianimg.com,DIRECT",
		"DOMAIN-SUFFIX,yinxiang.com,DIRECT",
		"DOMAIN-SUFFIX,ykimg.com,DIRECT",
		"DOMAIN-SUFFIX,youdao.com,DIRECT",
		"DOMAIN-SUFFIX,youku.com,DIRECT",
		"DOMAIN-SUFFIX,zealer.com,DIRECT",
		"DOMAIN-SUFFIX,zhihu.com,DIRECT",
		"DOMAIN-SUFFIX,zhimg.com,DIRECT",
		"DOMAIN-SUFFIX,zimuzu.tv,DIRECT",
		"DOMAIN-SUFFIX,zoho.com,DIRECT",
		"DOMAIN-KEYWORD,amazon, Proxies",
		"DOMAIN-KEYWORD,google, Proxies",
		"DOMAIN-KEYWORD,gmail, Proxies",
		"DOMAIN-KEYWORD,youtube, Proxies",
		"DOMAIN-KEYWORD,facebook, Proxies",
		"DOMAIN-SUFFIX,fb.me, Proxies",
		"DOMAIN-SUFFIX,fbcdn.net, Proxies",
		"DOMAIN-KEYWORD,twitter, Proxies",
		"DOMAIN-KEYWORD,instagram, Proxies",
		"DOMAIN-KEYWORD,dropbox, Proxies",
		"DOMAIN-KEYWORD,pronhub, Proxies",
		"DOMAIN-SUFFIX,cn.pronhub.com, Proxies",
		"DOMAIN-SUFFIX,pronhub.com, Proxies",
		"DOMAIN-SUFFIX,twimg.com, Proxies",
		"DOMAIN-KEYWORD,blogspot, Proxies",
		"DOMAIN-SUFFIX,youtu.be, Proxies",
		"DOMAIN-KEYWORD,whatsapp, Proxies",
		"DOMAIN-KEYWORD,admarvel,REJECT",
		"DOMAIN-KEYWORD,admaster,REJECT",
		"DOMAIN-KEYWORD,adsage,REJECT",
		"DOMAIN-KEYWORD,adsmogo,REJECT",
		"DOMAIN-KEYWORD,adsrvmedia,REJECT",
		"DOMAIN-KEYWORD,adwords,REJECT",
		"DOMAIN-KEYWORD,adservice,REJECT",
		"DOMAIN-SUFFIX,appsflyer.com,REJECT",
		"DOMAIN-KEYWORD,domob,REJECT",
		"DOMAIN-SUFFIX,doubleclick.net,REJECT",
		"DOMAIN-KEYWORD,duomeng,REJECT",
		"DOMAIN-KEYWORD,dwtrack,REJECT",
		"DOMAIN-KEYWORD,guanggao,REJECT",
		"DOMAIN-KEYWORD,lianmeng,REJECT",
		"DOMAIN-SUFFIX,mmstat.com,REJECT",
		"DOMAIN-KEYWORD,mopub,REJECT",
		"DOMAIN-KEYWORD,omgmta,REJECT",
		"DOMAIN-KEYWORD,openx,REJECT",
		"DOMAIN-KEYWORD,partnerad,REJECT",
		"DOMAIN-KEYWORD,pingfore,REJECT",
		"DOMAIN-KEYWORD,supersonicads,REJECT",
		"DOMAIN-KEYWORD,uedas,REJECT",
		"DOMAIN-KEYWORD,umeng,REJECT",
		"DOMAIN-KEYWORD,usage,REJECT",
		"DOMAIN-SUFFIX,vungle.com,REJECT",
		"DOMAIN-KEYWORD,wlmonitor,REJECT",
		"DOMAIN-KEYWORD,zjtoolbar,REJECT",
		"DOMAIN-SUFFIX,9to5mac.com, Proxies",
		"DOMAIN-SUFFIX,abpchina.org, Proxies",
		"DOMAIN-SUFFIX,adblockplus.org, Proxies",
		"DOMAIN-SUFFIX,adobe.com, Proxies",
		"DOMAIN-SUFFIX,akamaized.net, Proxies",
		"DOMAIN-SUFFIX,alfredapp.com, Proxies",
		"DOMAIN-SUFFIX,amplitude.com, Proxies",
		"DOMAIN-SUFFIX,ampproject.org, Proxies",
		"DOMAIN-SUFFIX,android.com, Proxies",
		"DOMAIN-SUFFIX,angularjs.org, Proxies",
		"DOMAIN-SUFFIX,aolcdn.com, Proxies",
		"DOMAIN-SUFFIX,apkpure.com, Proxies",
		"DOMAIN-SUFFIX,appledaily.com, Proxies",
		"DOMAIN-SUFFIX,appshopper.com, Proxies",
		"DOMAIN-SUFFIX,appspot.com, Proxies",
		"DOMAIN-SUFFIX,arcgis.com, Proxies",
		"DOMAIN-SUFFIX,archive.org, Proxies",
		"DOMAIN-SUFFIX,armorgames.com, Proxies",
		"DOMAIN-SUFFIX,aspnetcdn.com, Proxies",
		"DOMAIN-SUFFIX,att.com, Proxies",
		"DOMAIN-SUFFIX,awsstatic.com, Proxies",
		"DOMAIN-SUFFIX,azureedge.net, Proxies",
		"DOMAIN-SUFFIX,azurewebsites.net, Proxies",
		"DOMAIN-SUFFIX,bing.com, Proxies",
		"DOMAIN-SUFFIX,bintray.com, Proxies",
		"DOMAIN-SUFFIX,bit.com, Proxies",
		"DOMAIN-SUFFIX,bit.ly, Proxies",
		"DOMAIN-SUFFIX,bitbucket.org, Proxies",
		"DOMAIN-SUFFIX,bjango.com, Proxies",
		"DOMAIN-SUFFIX,bkrtx.com, Proxies",
		"DOMAIN-SUFFIX,blog.com, Proxies",
		"DOMAIN-SUFFIX,blogcdn.com, Proxies",
		"DOMAIN-SUFFIX,blogger.com, Proxies",
		"DOMAIN-SUFFIX,blogsmithmedia.com, Proxies",
		"DOMAIN-SUFFIX,blogspot.com, Proxies",
		"DOMAIN-SUFFIX,blogspot.hk, Proxies",
		"DOMAIN-SUFFIX,bloomberg.com, Proxies",
		"DOMAIN-SUFFIX,box.com, Proxies",
		"DOMAIN-SUFFIX,box.net, Proxies",
		"DOMAIN-SUFFIX,cachefly.net, Proxies",
		"DOMAIN-SUFFIX,chromium.org, Proxies",
		"DOMAIN-SUFFIX,cl.ly, Proxies",
		"DOMAIN-SUFFIX,cloudflare.com, Proxies",
		"DOMAIN-SUFFIX,cloudfront.net, Proxies",
		"DOMAIN-SUFFIX,cloudmagic.com, Proxies",
		"DOMAIN-SUFFIX,cmail19.com, Proxies",
		"DOMAIN-SUFFIX,cnet.com, Proxies",
		"DOMAIN-SUFFIX,cocoapods.org, Proxies",
		"DOMAIN-SUFFIX,comodoca.com, Proxies",
		"DOMAIN-SUFFIX,crashlytics.com, Proxies",
		"DOMAIN-SUFFIX,culturedcode.com, Proxies",
		"DOMAIN-SUFFIX,d.pr, Proxies",
		"DOMAIN-SUFFIX,danilo.to, Proxies",
		"DOMAIN-SUFFIX,dayone.me, Proxies",
		"DOMAIN-SUFFIX,db.tt, Proxies",
		"DOMAIN-SUFFIX,deskconnect.com, Proxies",
		"DOMAIN-SUFFIX,disq.us, Proxies",
		"DOMAIN-SUFFIX,disqus.com, Proxies",
		"DOMAIN-SUFFIX,disquscdn.com, Proxies",
		"DOMAIN-SUFFIX,dnsimple.com, Proxies",
		"DOMAIN-SUFFIX,docker.com, Proxies",
		"DOMAIN-SUFFIX,dribbble.com, Proxies",
		"DOMAIN-SUFFIX,droplr.com, Proxies",
		"DOMAIN-SUFFIX,duckduckgo.com, Proxies",
		"DOMAIN-SUFFIX,dueapp.com, Proxies",
		"DOMAIN-SUFFIX,dytt8.net, Proxies",
		"DOMAIN-SUFFIX,edgecastcdn.net, Proxies",
		"DOMAIN-SUFFIX,edgekey.net, Proxies",
		"DOMAIN-SUFFIX,edgesuite.net, Proxies",
		"DOMAIN-SUFFIX,engadget.com, Proxies",
		"DOMAIN-SUFFIX,entrust.net, Proxies",
		"DOMAIN-SUFFIX,eurekavpt.com, Proxies",
		"DOMAIN-SUFFIX,evernote.com, Proxies",
		"DOMAIN-SUFFIX,fabric.io, Proxies",
		"DOMAIN-SUFFIX,fast.com, Proxies",
		"DOMAIN-SUFFIX,fastly.net, Proxies",
		"DOMAIN-SUFFIX,fc2.com, Proxies",
		"DOMAIN-SUFFIX,feedburner.com, Proxies",
		"DOMAIN-SUFFIX,feedly.com, Proxies",
		"DOMAIN-SUFFIX,feedsportal.com, Proxies",
		"DOMAIN-SUFFIX,fiftythree.com, Proxies",
		"DOMAIN-SUFFIX,firebaseio.com, Proxies",
		"DOMAIN-SUFFIX,flexibits.com, Proxies",
		"DOMAIN-SUFFIX,flickr.com, Proxies",
		"DOMAIN-SUFFIX,flipboard.com, Proxies",
		"DOMAIN-SUFFIX,g.co, Proxies",
		"DOMAIN-SUFFIX,gabia.net, Proxies",
		"DOMAIN-SUFFIX,geni.us, Proxies",
		"DOMAIN-SUFFIX,gfx.ms, Proxies",
		"DOMAIN-SUFFIX,ggpht.com, Proxies",
		"DOMAIN-SUFFIX,ghostnoteapp.com, Proxies",
		"DOMAIN-SUFFIX,git.io, Proxies",
		"DOMAIN-KEYWORD,github, Proxies",
		"DOMAIN-SUFFIX,globalsign.com, Proxies",
		"DOMAIN-SUFFIX,gmodules.com, Proxies",
		"DOMAIN-SUFFIX,godaddy.com, Proxies",
		"DOMAIN-SUFFIX,golang.org, Proxies",
		"DOMAIN-SUFFIX,gongm.in, Proxies",
		"DOMAIN-SUFFIX,goo.gl, Proxies",
		"DOMAIN-SUFFIX,goodreaders.com, Proxies",
		"DOMAIN-SUFFIX,goodreads.com, Proxies",
		"DOMAIN-SUFFIX,gravatar.com, Proxies",
		"DOMAIN-SUFFIX,gstatic.com, Proxies",
		"DOMAIN-SUFFIX,gvt0.com, Proxies",
		"DOMAIN-SUFFIX,hockeyapp.net, Proxies",
		"DOMAIN-SUFFIX,hotmail.com, Proxies",
		"DOMAIN-SUFFIX,icons8.com, Proxies",
		"DOMAIN-SUFFIX,ifixit.com, Proxies",
		"DOMAIN-SUFFIX,ift.tt, Proxies",
		"DOMAIN-SUFFIX,ifttt.com, Proxies",
		"DOMAIN-SUFFIX,iherb.com, Proxies",
		"DOMAIN-SUFFIX,imageshack.us, Proxies",
		"DOMAIN-SUFFIX,img.ly, Proxies",
		"DOMAIN-SUFFIX,imgur.com, Proxies",
		"DOMAIN-SUFFIX,imore.com, Proxies",
		"DOMAIN-SUFFIX,instapaper.com, Proxies",
		"DOMAIN-SUFFIX,ipn.li, Proxies",
		"DOMAIN-SUFFIX,is.gd, Proxies",
		"DOMAIN-SUFFIX,issuu.com, Proxies",
		"DOMAIN-SUFFIX,itgonglun.com, Proxies",
		"DOMAIN-SUFFIX,itun.es, Proxies",
		"DOMAIN-SUFFIX,ixquick.com, Proxies",
		"DOMAIN-SUFFIX,j.mp, Proxies",
		"DOMAIN-SUFFIX,js.revsci.net, Proxies",
		"DOMAIN-SUFFIX,jshint.com, Proxies",
		"DOMAIN-SUFFIX,jtvnw.net, Proxies",
		"DOMAIN-SUFFIX,justgetflux.com, Proxies",
		"DOMAIN-SUFFIX,kat.cr, Proxies",
		"DOMAIN-SUFFIX,klip.me, Proxies",
		"DOMAIN-SUFFIX,libsyn.com, Proxies",
		"DOMAIN-SUFFIX,linkedin.com, Proxies",
		"DOMAIN-SUFFIX,linode.com, Proxies",
		"DOMAIN-SUFFIX,lithium.com, Proxies",
		"DOMAIN-SUFFIX,littlehj.com, Proxies",
		"DOMAIN-SUFFIX,live.com, Proxies",
		"DOMAIN-SUFFIX,live.net, Proxies",
		"DOMAIN-SUFFIX,livefilestore.com, Proxies",
		"DOMAIN-SUFFIX,llnwd.net, Proxies",
		"DOMAIN-SUFFIX,macid.co, Proxies",
		"DOMAIN-SUFFIX,macromedia.com, Proxies",
		"DOMAIN-SUFFIX,macrumors.com, Proxies",
		"DOMAIN-SUFFIX,mashable.com, Proxies",
		"DOMAIN-SUFFIX,mathjax.org, Proxies",
		"DOMAIN-SUFFIX,medium.com, Proxies",
		"DOMAIN-SUFFIX,mega.co.nz, Proxies",
		"DOMAIN-SUFFIX,mega.nz, Proxies",
		"DOMAIN-SUFFIX,megaupload.com, Proxies",
		"DOMAIN-SUFFIX,microsofttranslator.com, Proxies",
		"DOMAIN-SUFFIX,mindnode.com, Proxies",
		"DOMAIN-SUFFIX,mobile01.com, Proxies",
		"DOMAIN-SUFFIX,modmyi.com, Proxies",
		"DOMAIN-SUFFIX,msedge.net, Proxies",
		"DOMAIN-SUFFIX,myfontastic.com, Proxies",
		"DOMAIN-SUFFIX,name.com, Proxies",
		"DOMAIN-SUFFIX,nextmedia.com, Proxies",
		"DOMAIN-SUFFIX,nsstatic.net, Proxies",
		"DOMAIN-SUFFIX,nssurge.com, Proxies",
		"DOMAIN-SUFFIX,nyt.com, Proxies",
		"DOMAIN-SUFFIX,nytimes.com, Proxies",
		"DOMAIN-SUFFIX,omnigroup.com, Proxies",
		"DOMAIN-SUFFIX,onedrive.com, Proxies",
		"DOMAIN-SUFFIX,onenote.com, Proxies",
		"DOMAIN-SUFFIX,ooyala.com, Proxies",
		"DOMAIN-SUFFIX,openvpn.net, Proxies",
		"DOMAIN-SUFFIX,openwrt.org, Proxies",
		"DOMAIN-SUFFIX,orkut.com, Proxies",
		"DOMAIN-SUFFIX,osxdaily.com, Proxies",
		"DOMAIN-SUFFIX,outlook.com, Proxies",
		"DOMAIN-SUFFIX,ow.ly, Proxies",
		"DOMAIN-SUFFIX,paddleapi.com, Proxies",
		"DOMAIN-SUFFIX,parallels.com, Proxies",
		"DOMAIN-SUFFIX,parse.com, Proxies",
		"DOMAIN-SUFFIX,pdfexpert.com, Proxies",
		"DOMAIN-SUFFIX,periscope.tv, Proxies",
		"DOMAIN-SUFFIX,pinboard.in, Proxies",
		"DOMAIN-SUFFIX,pinterest.com, Proxies",
		"DOMAIN-SUFFIX,pixelmator.com, Proxies",
		"DOMAIN-SUFFIX,pixiv.net, Proxies",
		"DOMAIN-SUFFIX,playpcesor.com, Proxies",
		"DOMAIN-SUFFIX,playstation.com, Proxies",
		"DOMAIN-SUFFIX,playstation.com.hk, Proxies",
		"DOMAIN-SUFFIX,playstation.net, Proxies",
		"DOMAIN-SUFFIX,playstationnetwork.com, Proxies",
		"DOMAIN-SUFFIX,pushwoosh.com, Proxies",
		"DOMAIN-SUFFIX,rime.im, Proxies",
		"DOMAIN-SUFFIX,servebom.com, Proxies",
		"DOMAIN-SUFFIX,sfx.ms, Proxies",
		"DOMAIN-SUFFIX,shadowsocks.org, Proxies",
		"DOMAIN-SUFFIX,sharethis.com, Proxies",
		"DOMAIN-SUFFIX,shazam.com, Proxies",
		"DOMAIN-SUFFIX,skype.com, Proxies",
		"DOMAIN-SUFFIX,smartdns Proxies.com, Proxies",
		"DOMAIN-SUFFIX,smartmailcloud.com, Proxies",
		"DOMAIN-SUFFIX,sndcdn.com, Proxies",
		"DOMAIN-SUFFIX,sony.com, Proxies",
		"DOMAIN-SUFFIX,soundcloud.com, Proxies",
		"DOMAIN-SUFFIX,sourceforge.net, Proxies",
		"DOMAIN-SUFFIX,spotify.com, Proxies",
		"DOMAIN-SUFFIX,squarespace.com, Proxies",
		"DOMAIN-SUFFIX,sstatic.net, Proxies",
		"DOMAIN-SUFFIX,st.luluku.pw, Proxies",
		"DOMAIN-SUFFIX,stackoverflow.com, Proxies",
		"DOMAIN-SUFFIX,startpage.com, Proxies",
		"DOMAIN-SUFFIX,staticflickr.com, Proxies",
		"DOMAIN-SUFFIX,steamcommunity.com, Proxies",
		"DOMAIN-SUFFIX,symauth.com, Proxies",
		"DOMAIN-SUFFIX,symcb.com, Proxies",
		"DOMAIN-SUFFIX,symcd.com, Proxies",
		"DOMAIN-SUFFIX,tapbots.com, Proxies",
		"DOMAIN-SUFFIX,tapbots.net, Proxies",
		"DOMAIN-SUFFIX,tdesktop.com, Proxies",
		"DOMAIN-SUFFIX,techcrunch.com, Proxies",
		"DOMAIN-SUFFIX,techsmith.com, Proxies",
		"DOMAIN-SUFFIX,thepiratebay.org, Proxies",
		"DOMAIN-SUFFIX,theverge.com, Proxies",
		"DOMAIN-SUFFIX,time.com, Proxies",
		"DOMAIN-SUFFIX,timeinc.net, Proxies",
		"DOMAIN-SUFFIX,tiny.cc, Proxies",
		"DOMAIN-SUFFIX,tinypic.com, Proxies",
		"DOMAIN-SUFFIX,tmblr.co, Proxies",
		"DOMAIN-SUFFIX,todoist.com, Proxies",
		"DOMAIN-SUFFIX,trello.com, Proxies",
		"DOMAIN-SUFFIX,trustasiassl.com, Proxies",
		"DOMAIN-SUFFIX,tumblr.co, Proxies",
		"DOMAIN-SUFFIX,tumblr.com, Proxies",
		"DOMAIN-SUFFIX,tweetdeck.com, Proxies",
		"DOMAIN-SUFFIX,tweetmarker.net, Proxies",
		"DOMAIN-SUFFIX,twitch.tv, Proxies",
		"DOMAIN-SUFFIX,txmblr.com, Proxies",
		"DOMAIN-SUFFIX,typekit.net, Proxies",
		"DOMAIN-SUFFIX,ubertags.com, Proxies",
		"DOMAIN-SUFFIX,ublock.org, Proxies",
		"DOMAIN-SUFFIX,ubnt.com, Proxies",
		"DOMAIN-SUFFIX,ulyssesapp.com, Proxies",
		"DOMAIN-SUFFIX,urchin.com, Proxies",
		"DOMAIN-SUFFIX,usertrust.com, Proxies",
		"DOMAIN-SUFFIX,v.gd, Proxies",
		"DOMAIN-SUFFIX,v2ex.com, Proxies",
		"DOMAIN-SUFFIX,vimeo.com, Proxies",
		"DOMAIN-SUFFIX,vimeocdn.com, Proxies",
		"DOMAIN-SUFFIX,vine.co, Proxies",
		"DOMAIN-SUFFIX,vivaldi.com, Proxies",
		"DOMAIN-SUFFIX,vox cdn.com, Proxies",
		"DOMAIN-SUFFIX,vsco.co, Proxies",
		"DOMAIN-SUFFIX,vultr.com, Proxies",
		"DOMAIN-SUFFIX,w.org, Proxies",
		"DOMAIN-SUFFIX,w3schools.com, Proxies",
		"DOMAIN-SUFFIX,webtype.com, Proxies",
		"DOMAIN-SUFFIX,wikiwand.com, Proxies",
		"DOMAIN-SUFFIX,wikileaks.org, Proxies",
		"DOMAIN-SUFFIX,wikimedia.org, Proxies",
		"DOMAIN-SUFFIX,wikipedia.com, Proxies",
		"DOMAIN-SUFFIX,wikipedia.org, Proxies",
		"DOMAIN-SUFFIX,windows.com, Proxies",
		"DOMAIN-SUFFIX,windows.net, Proxies",
		"DOMAIN-SUFFIX,wire.com, Proxies",
		"DOMAIN-SUFFIX,wordpress.com, Proxies",
		"DOMAIN-SUFFIX,workflowy.com, Proxies",
		"DOMAIN-SUFFIX,wp.com, Proxies",
		"DOMAIN-SUFFIX,wsj.com, Proxies",
		"DOMAIN-SUFFIX,wsj.net, Proxies",
		"DOMAIN-SUFFIX,xda developers.com, Proxies",
		"DOMAIN-SUFFIX,xeeno.com, Proxies",
		"DOMAIN-SUFFIX,xiti.com, Proxies",
		"DOMAIN-SUFFIX,yahoo.com, Proxies",
		"DOMAIN-SUFFIX,yimg.com, Proxies",
		"DOMAIN-SUFFIX,ying.com, Proxies",
		"DOMAIN-SUFFIX,yoyo.org, Proxies",
		"DOMAIN-SUFFIX,ytimg.com, Proxies",
		"DOMAIN-SUFFIX,telegra.ph, Proxies",
		"DOMAIN-SUFFIX,telegram.org, Proxies",
		"IP-CIDR,91.108.4.0/22, Proxies,no resolve",
		"IP-CIDR,91.108.8.0/21, Proxies,no resolve",
		"IP-CIDR,91.108.16.0/22, Proxies,no resolve",
		"IP-CIDR,91.108.56.0/22, Proxies,no resolve",
		"IP-CIDR,149.154.160.0/20, Proxies,no resolve",
		"IP-CIDR6,2001:67c:4e8::/48, Proxies,no resolve",
		"IP-CIDR6,2001:b28:f23d::/48, Proxies,no resolve",
		"IP-CIDR6,2001:b28:f23f::/48, Proxies,no resolve",
		"IP-CIDR,120.232.181.162/32, Proxies,no resolve",
		"IP-CIDR,120.241.147.226/32, Proxies,no resolve",
		"IP-CIDR,120.253.253.226/32, Proxies,no resolve",
		"IP-CIDR,120.253.255.162/32, Proxies,no resolve",
		"IP-CIDR,120.253.255.34/32, Proxies,no resolve",
		"IP-CIDR,120.253.255.98/32, Proxies,no resolve",
		"IP-CIDR,180.163.150.162/32, Proxies,no resolve",
		"IP-CIDR,180.163.150.34/32, Proxies,no resolve",
		"IP-CIDR,180.163.151.162/32, Proxies,no resolve",
		"IP-CIDR,180.163.151.34/32, Proxies,no resolve",
		"IP-CIDR,203.208.39.0/24, Proxies,no resolve",
		"IP-CIDR,203.208.40.0/24, Proxies,no resolve",
		"IP-CIDR,203.208.41.0/24, Proxies,no resolve",
		"IP-CIDR,203.208.43.0/24, Proxies,no resolve",
		"IP-CIDR,203.208.50.0/24, Proxies,no resolve",
		"IP-CIDR,220.181.174.162/32, Proxies,no resolve",
		"IP-CIDR,220.181.174.226/32, Proxies,no resolve",
		"IP-CIDR,220.181.174.34/32, Proxies,no resolve",
		"DOMAIN,injections.adguard.org,DIRECT",
		"DOMAIN,local.adguard.org,DIRECT",
		"DOMAIN-SUFFIX,local,DIRECT",
		"IP-CIDR,127.0.0.0/8,DIRECT",
		"IP-CIDR,172.16.0.0/12,DIRECT",
		"IP-CIDR,192.168.0.0/16,DIRECT",
		"IP-CIDR,10.0.0.0/8,DIRECT",
		"IP-CIDR,17.0.0.0/8,DIRECT",
		"IP-CIDR,100.64.0.0/10,DIRECT",
		"IP-CIDR,224.0.0.0/4,DIRECT",
		"IP-CIDR6,fe80::/10,DIRECT",
		"GEOIP,CN,DIRECT",
		"MATCH, Proxies",
	}
	return rules
}

// 所有 clash 规则列表
func ClashRulelist() ([]conf.Proxies, int64, error) {
	c := db.Client.Database("xary").Collection("subconf")
	total, _ := c.CountDocuments(context.TODO(), bson.M{})

	carr := []conf.Proxies{}
	e := []conf.Proxies{}

	cur, err := c.Find(
		context.TODO(),
		bson.M{},
	)
	if err != nil {
		return e, 0, fmt.Errorf("获取clash 规则出错")
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		result := conf.Proxies{}
		err := cur.Decode(&result)
		if err != nil {
			return e, 0, fmt.Errorf("获取clash 规则出错")
		}
		carr = append(carr, result)
	}
	return carr, total, nil
}

func init() {
	db.ConnectDB()
}
