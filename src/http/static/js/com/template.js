import {Utils} from "./utils.js";


export class Template {
    constructor() {
        this.i18n = $.i18n();
        this.i18n.locale = navigator.language || navigator.userLanguage;

        this.i = this.imports();
    }

    imports() {
        let m = template.defaults.imports;
        m.a2n = function (data, n) {
            return Utils.i2n(data, n);
        };
        m.prettyKiB = function (data, fra) {
            let dec = data;
            fra = fra === undefined ? 2 : fra;
            if (dec < 1024) {
                return dec.toFixed(fra)+"KiB";
            }
            dec /= 1024.0;
            if (dec < 1024) {
                return dec.toFixed(fra)+"MiB";
            }
            dec /= 1024.0;
            if (dec < 1024) {
                return dec.toFixed(fra)+"GiB";
            }
            dec /=  1024.0;
            return dec.toFixed(fra) + "TiB"
        };
        m.prettyByte = function (data, fra) {
            let dec = data;
            fra = fra === undefined ? 2 : fra;
            if (dec < 1024) {
                return dec.toFixed(fra)+"B";
            }
            dec = dec / 1024.0;
            if (dec < 1024) {
                return dec.toFixed(fra)+"KiB";
            }
            dec /= 1024.0;
            if (dec < 1024) {
                return dec.toFixed(fra)+"MiB";
            }
            dec = dec / 1024.0;
            if (dec < 1024) {
                return dec.toFixed(fra)+"GiB";
            }
            dec = dec / 1024;
            return dec.toFixed(fra) + "TiB"
        };
        m.figureCpuUsed = function (free, total) {
            return ((1000 - free) / 1000 * total).toFixed(2)
        };
        m.figureCpuFree = function (free, total) {
            return (free / 1000 * total).toFixed(2)
        };
        m.netmask2prefix = function (netmask) {
            if (!netmask) return undefined;
            return netmask.split('.').map(Number)
                .map(part => (part >>> 0).toString(2))
                .join('').split('1').length - 1;
        };
        m.prefix2netmask = function (prefix) {
            if (!prefix) return undefined;
            let mask = [];
            for(let i = 0;i < 4; i++) {
                let n = Math.min(prefix, 8);
                mask.push(256 - Math.pow(2, 8-n));
                prefix -= n;
            }
            return mask.join('.');
        };
        m.vncPassword = function (inst) {
            return Utils.graphic(inst, 'vnc', 'password')
        };
        m.spicePassword = function (inst) {
            return Utils.graphic(inst, 'spice', 'password')
        };
        m.i = function (value) {
            return $.i18n(value);
        };
        return m;
    }

    promise() {
        console.log("Template.promise", this.i18n.locale);

        let i18n = this.i18n;
        return new Promise(function (resolve) {
            i18n.load(`/static/i18n/${i18n.locale}.json`, i18n.locale)
                .done(function() {
                    resolve()
                });
        })
    }
}
