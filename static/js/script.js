// 페이지에 표시되는 월~금 날짜를 반환
function GetTimeTableDays() {
    var days = new Array();
    var now = new Date();
    days[0] = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0);

    if (now.getDay() == 0) {
        // 일요일이면 6일 뺌
        days[0] = new Date(days[0].getTime() - 86400000 * 6);
    } else if (now.getDay() == 1) {
        // 월요일이면 8시 20분 이후인지 검사
        if (now.getHours() <= 8) {
            if (now.getHours() == 8) {
                if (now.getMinutes() < 20) {
                    days[0] = new Date(days[0].getTime() - (86400000 * 6));
                }
            } else {
                // 8시 전이면 저번주
                days[0] = new Date(days[0].getTime() - (86400000 * 6));
            }
        }
    } else {
        // 월~토이면 n-1(월: 1, 토: 6)일 뺌
        days[0] = new Date(days[0].getTime() + (86400000 * (1 - days[0].getDay())));
    }
    days[1] = new Date(days[0].getTime() + 86400000);
    days[2] = new Date(days[1].getTime() + 86400000);
    days[3] = new Date(days[2].getTime() + 86400000);
    days[4] = new Date(days[3].getTime() + 86400000);

    return days;
}

// 문자열 자르기
function split(str, separator, limit) {
    str = str.split(separator);

    if (str.length > limit) {
        var ret = str.splice(0, limit);
        ret.push(str.join(separator));

        return ret;
    }

    return str;
}

// QuerySelector
function getParameterByName(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, '\\$&');
    var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, ' '));
}

// yyyy-MM-dd 형식으로 포맷하는 함수 추가
Date.prototype.yyyymmdd = function () {
    var yyyy = this.getFullYear().toString();
    var mm = (this.getMonth() + 1).toString();
    var dd = this.getDate().toString();

    return yyyy + "-" + (mm[1] ? mm : '0' + mm[0]) + "-" + (dd[1] ? dd : '0' + dd[0]);
}
