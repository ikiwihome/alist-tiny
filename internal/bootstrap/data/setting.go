package data

import (
	"github.com/alist-org/alist/v3/cmd/flags"
	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/pkg/utils/random"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var initialSettingItems []model.SettingItem

func initSettings() {
	InitialSettings()
	// check deprecated
	settings, err := op.GetSettingItems()
	if err != nil {
		utils.Log.Fatalf("failed get settings: %+v", err)
	}

	for i := range settings {
		if !isActive(settings[i].Key) && settings[i].Flag != model.DEPRECATED {
			settings[i].Flag = model.DEPRECATED
			err = op.SaveSettingItem(&settings[i])
			if err != nil {
				utils.Log.Fatalf("failed save setting: %+v", err)
			}
		}
	}

	// create or save setting
	for i := range initialSettingItems {
		item := &initialSettingItems[i]
		// err
		stored, err := op.GetSettingItemByKey(item.Key)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Log.Fatalf("failed get setting: %+v", err)
			continue
		}
		// save
		if stored != nil && item.Key != conf.VERSION {
			item.Value = stored.Value
		}
		if stored == nil || *item != *stored {
			err = op.SaveSettingItem(item)
			if err != nil {
				utils.Log.Fatalf("failed save setting: %+v", err)
			}
		} else {
			// Not save so needs to execute hook
			_, err = op.HandleSettingItemHook(item)
			if err != nil {
				utils.Log.Errorf("failed to execute hook on %s: %+v", item.Key, err)
			}
		}
	}
}

func isActive(key string) bool {
	for _, item := range initialSettingItems {
		if item.Key == key {
			return true
		}
	}
	return false
}

func InitialSettings() []model.SettingItem {
	var token string
	if flags.Dev {
		token = "dev_token"
	} else {
		token = random.Token()
	}
	initialSettingItems = []model.SettingItem{
		// site settings
		{Key: conf.VERSION, Value: conf.Version, Type: conf.TypeString, Group: model.SITE, Flag: model.READONLY},
		//{Key: conf.ApiUrl, Value: "", Type: conf.TypeString, Group: model.SITE},
		//{Key: conf.BasePath, Value: "", Type: conf.TypeString, Group: model.SITE},
		{Key: conf.SiteTitle, Value: "Cloud Space", Type: conf.TypeString, Group: model.SITE},
		{Key: conf.Announcement, Value: "", Type: conf.TypeText, Group: model.SITE},
		{Key: "pagination_type", Value: "all", Type: conf.TypeSelect, Options: "all,pagination,load_more,auto_load_more", Group: model.SITE},
		{Key: "default_page_size", Value: "30", Type: conf.TypeNumber, Group: model.SITE},
		{Key: conf.AllowIndexed, Value: "false", Type: conf.TypeBool, Group: model.SITE},
		{Key: conf.AllowMounted, Value: "true", Type: conf.TypeBool, Group: model.SITE},
		{Key: conf.RobotsTxt, Value: "User-agent: *\nDisallow: /", Type: conf.TypeText, Group: model.SITE},
		// style settings
		{Key: conf.Logo, Value: "https://cdn.jsdelivr.net/gh/alist-org/logo@main/logo.svg", Type: conf.TypeText, Group: model.STYLE},
		{Key: conf.Favicon, Value: "https://cdn.jsdelivr.net/gh/alist-org/logo@main/logo.svg", Type: conf.TypeString, Group: model.STYLE},
		{Key: conf.MainColor, Value: "#1890ff", Type: conf.TypeString, Group: model.STYLE},
		{Key: "home_icon", Value: "🏠", Type: conf.TypeString, Group: model.STYLE},
		{Key: "home_container", Value: "max_980px", Type: conf.TypeSelect, Options: "max_980px,hope_container", Group: model.STYLE},
		{Key: "settings_layout", Value: "list", Type: conf.TypeSelect, Options: "list,responsive", Group: model.STYLE},
		// preview settings
		{Key: conf.TextTypes, Value: "txt,htm,html,xml,java,properties,sql,js,md,json,conf,ini,vue,php,py,bat,gitignore,yml,go,sh,c,cpp,h,hpp,tsx,vtt,srt,ass,rs,lrc", Type: conf.TypeText, Group: model.PREVIEW, Flag: model.PRIVATE},
		{Key: conf.AudioTypes, Value: "mp3,flac,ogg,m4a,wav,opus,wma", Type: conf.TypeText, Group: model.PREVIEW, Flag: model.PRIVATE},
		{Key: conf.VideoTypes, Value: "mp4,mkv,avi,mov,rmvb,webm,flv", Type: conf.TypeText, Group: model.PREVIEW, Flag: model.PRIVATE},
		{Key: conf.ImageTypes, Value: "jpg,tiff,jpeg,png,gif,bmp,svg,ico,swf,webp", Type: conf.TypeText, Group: model.PREVIEW, Flag: model.PRIVATE},
		//{Key: conf.OfficeTypes, Value: "doc,docx,xls,xlsx,ppt,pptx", Type: conf.TypeText, Group: model.PREVIEW, Flag: model.PRIVATE},
		{Key: conf.ProxyTypes, Value: "m3u8", Type: conf.TypeText, Group: model.PREVIEW, Flag: model.PRIVATE},
		{Key: conf.ProxyIgnoreHeaders, Value: "authorization,referer", Type: conf.TypeText, Group: model.PREVIEW, Flag: model.PRIVATE},
		{Key: "external_previews", Value: `{}`, Type: conf.TypeText, Group: model.PREVIEW},
		{Key: "iframe_previews", Value: `{
	"doc,docx,xls,xlsx,ppt,pptx": {
		"Office预览":"https://view.officeapps.live.com/op/view.aspx?src=$e_url"
	},
	"pdf": {
		"PDF预览":"https://alist-org.github.io/pdf.js/web/viewer.html?file=$e_url"
	},
	"epub": {
		"EPUB预览":"https://alist-org.github.io/static/epub.js/viewer.html?url=$e_url"
	}
}`, Type: conf.TypeText, Group: model.PREVIEW},
		//		{Key: conf.OfficeViewers, Value: `{
		//	"Microsoft":"https://view.officeapps.live.com/op/view.aspx?src=$url",
		//	"Google":"https://docs.google.com/gview?url=$url&embedded=true",
		//}`, Type: conf.TypeText, Group: model.PREVIEW},
		//		{Key: conf.PdfViewers, Value: `{
		//	"pdf.js":"https://alist-org.github.io/pdf.js/web/viewer.html?file=$url"
		//}`, Type: conf.TypeText, Group: model.PREVIEW},
		{Key: "audio_cover", Value: "https://jsd.nn.ci/gh/alist-org/logo@main/logo.svg", Type: conf.TypeString, Group: model.PREVIEW},
		{Key: conf.AudioAutoplay, Value: "false", Type: conf.TypeBool, Group: model.PREVIEW},
		{Key: conf.VideoAutoplay, Value: "false", Type: conf.TypeBool, Group: model.PREVIEW},
		// global settings
		{Key: conf.HideFiles, Value: "/\\/README.md/i", Type: conf.TypeText, Group: model.GLOBAL},
		{Key: "package_download", Value: "true", Type: conf.TypeBool, Group: model.GLOBAL},
		{Key: conf.CustomizeHead, Value: `<script src="https://polyfill.alicdn.com/v3/polyfill.min.js?features=String.prototype.replaceAll"></script>
		<style>
		/*白天模式 搜索主体+毛玻璃*/
		.hope-ui-light .hope-c-PJLV-iiBaxsN-css{
		   background-color: rgba(255,255,255,0.2)!important;
		   backdrop-filter: blur(10px)!important;
		}
		
		/*白天模式 搜索栏输入框+毛玻璃*/
		.hope-ui-light .hope-c-kvTTWD-hYRNAb-variant-filled{
		   background-color: rgba(255,255,255,0.2)!important;
		   backdrop-filter: blur(10px)!important;
		}
		
		/*白天模式 搜索按钮+毛玻璃*/
		.hope-ui-light .hope-c-PJLV-ikEIIxw-css{
		   background-color: rgba(255,255,255,0.2)!important;
		   backdrop-filter: blur(10px)!important;
		   padding: var(--hope-space-1)!important;
		}
		
		/*夜间模式搜索主体+毛玻璃*/
		.hope-ui-dark .hope-c-PJLV-iiBaxsN-css{
			background-color: rgb(0 0 0 / 10%)!important;
			backdrop-filter: blur(10px)!important;
		}
		
		/*夜间模式搜索栏+毛玻璃*/
		.hope-ui-dark .hope-c-kvTTWD-hYRNAb-variant-filled{
			background-color: rgb(0 0 0 / 10%)!important;
			backdrop-filter: blur(10px)!important;
		}
		
		/*夜间模式 搜索按钮+毛玻璃*/
		.hope-ui-dark .hope-c-PJLV-ikEIIxw-css{
			background-color: rgb(0 0 0 / 10%)!important;
			backdrop-filter: blur(10px)!important;
			padding: var(--hope-space-1)!important;
		}
		</style>`, Type: conf.TypeText, Group: model.GLOBAL, Flag: model.PRIVATE},
		{Key: conf.CustomizeBody, Value: `<!-- 网页鼠标点击特效 - 爱心 -->
		<script type="text/javascript">
				 ! function (e, t, a) {
					function r() {
						for (var e = 0; e < s.length; e++) s[e].alpha <= 0 ? (t.body.removeChild(s[e].el), s.splice(e, 1)) : (s[
								e].y--, s[e].scale += .004, s[e].alpha -= .013, s[e].el.style.cssText = "left:" + s[e].x +
							"px;top:" + s[e].y + "px;opacity:" + s[e].alpha + ";transform:scale(" + s[e].scale + "," + s[e]
							.scale + ") rotate(45deg);background:" + s[e].color + ";z-index:99999");
						requestAnimationFrame(r)
					}
					function n() {
						var t = "function" == typeof e.onclick && e.onclick;
						e.onclick = function (e) {
							t && t(), o(e)
						}
					}
		 
					function o(e) {
						var a = t.createElement("div");
						a.className = "heart", s.push({
							el: a,
							x: e.clientX - 5,
							y: e.clientY - 5,
							scale: 1,
							alpha: 1,
							color: c()
						}), t.body.appendChild(a)
					}
		 
					function i(e) {
						var a = t.createElement("style");
						a.type = "text/css";
						try {
							a.appendChild(t.createTextNode(e))
						} catch (t) {
							a.styleSheet.cssText = e
						}
						t.getElementsByTagName("head")[0].appendChild(a)
					}
		 
					function c() {
						return "rgb(" + ~~(255 * Math.random()) + "," + ~~(255 * Math.random()) + "," + ~~(255 * Math
							.random()) + ")"
					}
					var s = [];
					e.requestAnimationFrame = e.requestAnimationFrame || e.webkitRequestAnimationFrame || e
						.mozRequestAnimationFrame || e.oRequestAnimationFrame || e.msRequestAnimationFrame || function (e) {
							setTimeout(e, 1e3 / 60)
						}, i(
							".heart{width: 10px;height: 10px;position: fixed;background: #f00;transform: rotate(45deg);-webkit-transform: rotate(45deg);-moz-transform: rotate(45deg);}.heart:after,.heart:before{content: '';width: inherit;height: inherit;background: inherit;border-radius: 50%;-webkit-border-radius: 50%;-moz-border-radius: 50%;position: fixed;}.heart:after{top: -5px;}.heart:before{left: -5px;}"
						), n(), r()
				}(window, document);
			
		</script>`, Type: conf.TypeText, Group: model.GLOBAL, Flag: model.PRIVATE},
		{Key: conf.LinkExpiration, Value: "0", Type: conf.TypeNumber, Group: model.GLOBAL, Flag: model.PRIVATE},
		{Key: conf.SignAll, Value: "false", Type: conf.TypeBool, Group: model.GLOBAL, Flag: model.PRIVATE},
		{Key: conf.PrivacyRegs, Value: `(?:(?:\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.){3}(?:\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])
([[:xdigit:]]{1,4}(?::[[:xdigit:]]{1,4}){7}|::|:(?::[[:xdigit:]]{1,4}){1,6}|[[:xdigit:]]{1,4}:(?::[[:xdigit:]]{1,4}){1,5}|(?:[[:xdigit:]]{1,4}:){2}(?::[[:xdigit:]]{1,4}){1,4}|(?:[[:xdigit:]]{1,4}:){3}(?::[[:xdigit:]]{1,4}){1,3}|(?:[[:xdigit:]]{1,4}:){4}(?::[[:xdigit:]]{1,4}){1,2}|(?:[[:xdigit:]]{1,4}:){5}:[[:xdigit:]]{1,4}|(?:[[:xdigit:]]{1,4}:){1,6}:)
(?U)access_token=(.*)&`,
			Type: conf.TypeText, Group: model.GLOBAL, Flag: model.PRIVATE},
		{Key: conf.OcrApi, Value: "https://api.nn.ci/ocr/file/json", Type: conf.TypeString, Group: model.GLOBAL},
		{Key: conf.FilenameCharMapping, Value: `{"/": "|"}`, Type: conf.TypeText, Group: model.GLOBAL},
		{Key: conf.ForwardDirectLinkParams, Value: "false", Type: conf.TypeBool, Group: model.GLOBAL},

		// aria2 settings
		{Key: conf.Aria2Uri, Value: "http://localhost:6800/jsonrpc", Type: conf.TypeString, Group: model.ARIA2, Flag: model.PRIVATE},
		{Key: conf.Aria2Secret, Value: "", Type: conf.TypeString, Group: model.ARIA2, Flag: model.PRIVATE},

		// single settings
		{Key: conf.Token, Value: token, Type: conf.TypeString, Group: model.SINGLE, Flag: model.PRIVATE},
		{Key: conf.SearchIndex, Value: "none", Type: conf.TypeSelect, Options: "database,database_non_full_text,bleve,none", Group: model.INDEX},
		{Key: conf.AutoUpdateIndex, Value: "false", Type: conf.TypeBool, Group: model.INDEX},
		{Key: conf.IgnorePaths, Value: "", Type: conf.TypeText, Group: model.INDEX, Flag: model.PRIVATE, Help: `one path per line`},
		{Key: conf.MaxIndexDepth, Value: "20", Type: conf.TypeNumber, Group: model.INDEX, Flag: model.PRIVATE, Help: `max depth of index`},
		{Key: conf.IndexProgress, Value: "{}", Type: conf.TypeText, Group: model.SINGLE, Flag: model.PRIVATE},

		// SSO settings
		{Key: conf.SSOLoginEnabled, Value: "false", Type: conf.TypeBool, Group: model.SSO, Flag: model.PUBLIC},
		{Key: conf.SSOLoginplatform, Type: conf.TypeSelect, Options: "Casdoor,Github,Microsoft,Google,Dingtalk", Group: model.SSO, Flag: model.PUBLIC},
		{Key: conf.SSOClientId, Value: "", Type: conf.TypeString, Group: model.SSO, Flag: model.PRIVATE},
		{Key: conf.SSOClientSecret, Value: "", Type: conf.TypeString, Group: model.SSO, Flag: model.PRIVATE},
		{Key: conf.SSOOrganizationName, Value: "", Type: conf.TypeString, Group: model.SSO, Flag: model.PRIVATE},
		{Key: conf.SSOApplicationName, Value: "", Type: conf.TypeString, Group: model.SSO, Flag: model.PRIVATE},
		{Key: conf.SSOEndpointName, Value: "", Type: conf.TypeString, Group: model.SSO, Flag: model.PRIVATE},
		{Key: conf.SSOJwtPublicKey, Value: "", Type: conf.TypeString, Group: model.SSO, Flag: model.PRIVATE},

		// qbittorrent settings
		{Key: conf.QbittorrentUrl, Value: "http://admin:adminadmin@localhost:8080/", Type: conf.TypeString, Group: model.SINGLE, Flag: model.PRIVATE},
		{Key: conf.QbittorrentSeedtime, Value: "0", Type: conf.TypeNumber, Group: model.SINGLE, Flag: model.PRIVATE},
	}
	if flags.Dev {
		initialSettingItems = append(initialSettingItems, []model.SettingItem{
			{Key: "test_deprecated", Value: "test_value", Type: conf.TypeString, Flag: model.DEPRECATED},
			{Key: "test_options", Value: "a", Type: conf.TypeSelect, Options: "a,b,c"},
			{Key: "test_help", Type: conf.TypeString, Help: "this is a help message"},
		}...)
	}
	return initialSettingItems
}
