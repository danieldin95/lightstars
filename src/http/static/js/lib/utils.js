export class Utils {
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
