
function getQueryParams() {
    const params = {};
    const query = window.location.search.substring(1);
    const regex = /([^&=]+)=([^&]*)/g;
    let match;
    while ((match = regex.exec(query))) {
        params[decodeURIComponent(match[1])] = decodeURIComponent(match[2]);
    }
    return params;
}
const params = getQueryParams();
if (params.err) {
    let message = "";
    if (params.err === "successful_accept") {
        message = "Вы успешно выдали награду!";
    } else if (params.err === "already_completed") {
        message = "Вы уже выдали награду!";
    } else if (params.err === "wrong_claim_id") {
        message = "Некорректный ID подозреваемого!";
    } else {
        message = "Произошла ошибка!";
    }
    const notificationDiv = document.getElementById('notification');
    notificationDiv.innerText = message;
    notificationDiv.style.display = 'block';
    const newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
    window.history.replaceState({ path: newUrl }, '', newUrl);
    setTimeout(() => {
            notificationDiv.style.display = 'none';
        }, 5000);
}
