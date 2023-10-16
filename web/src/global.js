// 按钮恢复时间
const btRecoverTime = 3 * 1000

function getUrl() {
    return window.location.protocol + "//" + window.location.hostname + ":" + (10000 + parseInt(window.location.port))
}


export {
    getUrl,
    btRecoverTime
}