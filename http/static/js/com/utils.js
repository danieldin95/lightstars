export class Utils {
    // num: int
    static iton (data, n) {
        return (Array(n).join(0) + data).slice(-n);
    }

    // num: string
    static aton(data, n) {
        let num = "" + data;
        if (num.length > n) {
            return num
        }
        let ret = "";
        for (let i = 0; i < n - num.length; i++) {
            ret += "0"
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
}