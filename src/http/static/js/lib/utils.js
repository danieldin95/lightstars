export class Utils {
    static os() {
        let OSName = "Unknown";
        if (window.navigator.userAgent.indexOf("Windows NT 10.0")!== -1) OSName="Windows 10";
        if (window.navigator.userAgent.indexOf("Windows NT 6.2") !== -1) OSName="Windows 8";
        if (window.navigator.userAgent.indexOf("Windows NT 6.1") !== -1) OSName="Windows 7";
        if (window.navigator.userAgent.indexOf("Windows NT 6.0") !== -1) OSName="Windows Vista";
        if (window.navigator.userAgent.indexOf("Windows NT 5.1") !== -1) OSName="Windows XP";
        if (window.navigator.userAgent.indexOf("Windows NT 5.0") !== -1) OSName="Windows 2000";
        if (window.navigator.userAgent.indexOf("Mac")            !== -1) OSName="Mac/iOS";
        if (window.navigator.userAgent.indexOf("X11")            !== -1) OSName="UNIX";
        if (window.navigator.userAgent.indexOf("Linux")          !== -1) OSName="Linux";
        return OSName.toLocaleLowerCase();
    }

    static firefox() {
        return navigator.userAgent.match(/firefox/i);
    }

    static chrome() {
        return navigator.userAgent.match(/chrome/i);
    }

    // num: int
    static i2n (data, n) {
        return (Array(n).join(0) + data).slice(-n);
    }

    // num: string
    static a2n(data, n) {
        let num = "" + data;
        if (num.length > n) {
            return num;
        }
        let ret = "";
        for (let i = 0; i < n - num.length; i++) {
            ret += "0";
        }
        return ret + num;
    }

    static toJSON (data) {
        let ret = {};
        data.forEach(function (element, index) {
            let key = element['name'];
            let value = element['value'];
            if (key && value) {
                ret[key] = value;
            }
        });
        return ret;
    }

    static toKiB(size, unit) {
        if (unit === 'MiB') return size * 1024;
        if (unit === 'GiB') return size * 1024 * 1024;
    }

    static graphic(instance, type, name) {
        if (!instance.graphics) {
            return ""
        }
        for (let g of instance.graphics) {
            if (g.type === type) {
                return g[name];
            }
        }
        return ''
    }

    static basename(str) {
        let idx = str.lastIndexOf('/');
        idx = idx > -1 ? idx : str.lastIndexOf('\\');
        if (idx < 0) {
            return str;
        }
        return str.substring(idx + 1);
    }
}
