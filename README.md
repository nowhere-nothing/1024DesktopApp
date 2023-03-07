# 1024DesktopApp

草榴webview客户端
只支持windows

***TODO:***

- [ ] 根据不同的网站加载不同的js
- [ ] store增加mongo或者seaweedFS
- [ ] 下载成功/失败链接入库 sqlite/mysql/mongo
- [ ] 日志持久化



### 不同字段的格式

document.baseURI 获取当前页面的url
pc端格式: 'https://t66y.com/htm_data/2301/8/5476460.html'
title: <h4 class="f16">
pic: <div class="tpc_content do_not_catch" id="conttpc">
<img data-link ess-data src>

document.baseURI
手机端格式: 'https://t66y.com/htm_mob/2301/8/5470991.html'
title: <div class="f18">
pic: <div class="tpc_cont" id="conttpc">
<img data-link ess-data src">

/*
interface Handler {
// 1. 获取页面的所有图片
PageImages();

// 2. 获取页面列表的所有链接 => 在go中解析所有链接 => 将所有链接中的图片下载
PageTitle();

// 3. 获取标题 这里的标题是网页的标题或者具体内容中的标题
PagePosts();

// 4. 杂项例如去掉点击图片会跳转的事件
PreMisc();

    //AfterMisc();
}

// ***这里要区分移动端和手机端***
class T66y implements Handler {
PageImages() {

    }

    PagePosts() {

    }

    PageTitle() {

    }

    PreMisc() {
        this.RemoveImageClickEvent();
    }

    RemoveImageClickEvent() {

    }
}

// ==================== 4kup =====================
class K4up implements Handler {
PageImages() {
}

    PageTitle() {
    }

    PagePosts() {
    }

    PreMisc() {
    }
}

// ==================== ghost =====================
class Ghost implements Handler {
PagePosts() {
}

    PageImages() {
    }

    PageTitle() {
    }

    PreMisc() {
    }
}

// ==================== 没啥意思的 ===================
class IKanIns {
}

class SocialGirls {
}

// =================== 通用 =====================

// handler通过模板注入 默认t66y
let handler: Handler = new T66y();
*/

// ==================== 事件 ====================
