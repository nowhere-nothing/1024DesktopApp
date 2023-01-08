function resetImage() {
    let out_images = document.getElementById('conttpc')
        ?.getElementsByTagName('img');
    if (out_images) {
        for (let i in out_images) {
            out_images[i].onclick = null;
        }
    }
}

/*
document.baseURI
'https://t66y.com/htm_data/2301/8/5476460.html'
title: <h4 class="f16">
pic: <div class="tpc_content do_not_catch" id="conttpc">
         <img data-link ess-data src>
 */

/*
document.baseURI
'https://t66y.com/htm_mob/2301/8/5470991.html'
title: <div class="f18">
pic: <div class="tpc_cont" id="conttpc">
        <img data-link ess-data src">
 */
function getImages() {
    let out_images = document.getElementById('conttpc')
        ?.getElementsByTagName("img");
    let list = [];
    if (out_images) {
        for (let i in out_images) {
            if (out_images[i]?.src) {
                list.push(out_images[i].src)
            }
        }
    }
    return list;
}

function getTitle() {
    let titleDiv = document.getElementsByClassName("f18")[0];
    if (titleDiv) {
        return titleDiv.innerHTML;
    }
    return document.title;
}

function addDownloadBtn() {
    let btn = document.createElement('button');
    btn.innerHTML = '下载';
    btn.id = "downloadBtn";
    // todo check is not mob html convert it
    btn.onclick = () => {
        let images = getImages();
        let title = getTitle();
        download(title, document.baseURI, images);
    }
    let bdy = document.getElementsByTagName("body")[0];
    bdy.appendChild(btn);
}

function addProgress() {
    let pb = document.createElement('div');
    pb.id = "progress"
    pb.innerHTML = `
<label for="progressBar">
<progress id="progressBar" max="0" value="0"/>
</label>
<span id="progressBarTip">0/0</span>`;
    document.getElementsByTagName("body")[0].appendChild(pb);
}

const save_folder_key = "save_folder"

function addSaveFolderBtn() {
    let tb = document.createElement("button");
    tb.innerHTML = "路径";
    tb.id = "saveFolderBtn";
    tb.onclick = function () {
        let v = prompt("输入保存路径")
        if (v) {
            setSaveFolder(v).then(() => {
                localStorage.setItem(save_folder_key, v)
            }).catch(err => {
                alert(`保存路径错误 ${err}`)
            })
        }
    }
    document.getElementsByTagName("body")[0].appendChild(tb);
}

function addViewer() {
    let node = document.getElementById('conttpc')
    if (node) {
        const gallery = new Viewer(node, {
            backdrop: 'static',
            movable: false,
            rotatable: false,
        });
    }
}

const max_progress_key = "max_progress"
const cur_progress_key = "cur_progress"

function setGlobalProgress(max, val) {
    if (max === 0 && val === 0) {
        localStorage.removeItem(max_progress_key)
        localStorage.removeItem(cur_progress_key)
    } else {
        localStorage.setItem(max_progress_key, max);
        localStorage.setItem(cur_progress_key, val);
    }
    setProgress(max, val);
}

function setProgress(max, val) {
    let pb = document.getElementById("progressBar")
    if (pb) {
        pb.max = max;
        pb.value = val;
    }
    let tip = document.getElementById("progressBarTip")
    if (tip) {
        tip.innerText = `${val}/${max}`
    }
}

window.addEventListener('DOMContentLoaded', e => {
    const style = document.createElement("style");
    style.innerHTML = `
    #downloadBtn {
      position: fixed;
      right: 0;
      bottom: 0;
      width: 50px;
    }
    
    #progress {
    position: fixed;
    left: 0;
    bottom: 0;
    width: 300px;
    }
    
    #saveFolderBtn {
    position: fixed;
    right: 50px;
    bottom: 0;
    width: 50px;
    }
    `;
    document.head.append(style);
    resetImage();
    addDownloadBtn();
    addProgress();
    addSaveFolderBtn();
    addViewer();
    let m = localStorage.getItem(max_progress_key)
    let c = localStorage.getItem(cur_progress_key)
    if (m && c) {
        setProgress(m, c)
    }
})
