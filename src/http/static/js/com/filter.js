import {Utils} from "./utils.js";


export class Filters {
    constructor() {
        this.i18n = $.i18n();

        this.i18n.locale = navigator.language || navigator.userLanguage;
        console.log(this.i18n.locale);

        this.i = this.imports();
    }

    imports() {
        let i = template.defaults.imports;

        i.a2n = function (data, n) {
            return Utils.iton(data, n);
        };
        i.prettyKiB = function (data, fra) {
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
        i.prettyByte = function (data, fra) {
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
        i.figureCpuUsed = function (free, total) {
            return ((1000 - free) / 1000 * total).toFixed(2)
        };
        i.figureCpuFree = function (free, total) {
            return (free / 1000 * total).toFixed(2)
        };
        i.netmask2prefix = function (netmask) {
            if (!netmask) return undefined;
            return netmask.split('.').map(Number)
                .map(part => (part >>> 0).toString(2))
                .join('').split('1').length - 1;
        };
        i.prefix2netmask = function (prefix) {
            if (!prefix) return undefined;
            let mask = [];
            for(let i = 0;i < 4; i++) {
                let n = Math.min(prefix, 8);
                mask.push(256 - Math.pow(2, 8-n));
                prefix -= n;
            }
            return mask.join('.');
        };
        i.vncPassword = function (inst) {
            return Utils.graphic(inst, 'vnc', 'password')
        };
        i.spicePassword = function (inst) {
            return Utils.graphic(inst, 'spice', 'password')
        };
        i.i = function (value) {
            return $.i18n(value);
        };
        return i;
    }

    promise() {
        let i18n = this.i18n;
        return new Promise(function (resolve, reject) {
            i18n.load(`/static/i18n/${i18n.locale}.json`, i18n.locale)
                .done(function() {
                    console.log("loading.i18n done", this);
                    resolve()
                });
        })
    }
}
