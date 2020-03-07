import {Utils} from "./com/utils.js";

export class Filters {
    constructor() {
        template.defaults.imports.aton = function (data, n) {
            return Utils.iton(data, n);
        };

        template.defaults.imports.prettyKiB = function (data) {
            let dec = data;
            if (dec < 1024) {
                return dec.toFixed(2)+"KiB";
            }
            dec /= 1024.0;
            if (dec < 1024) {
                return dec.toFixed(2)+"MiB";
            }
            dec /= 1024.0;
            if (dec < 1024) {
                return dec.toFixed(2)+"GiB";
            }
            dec /=  1024.0;
            return dec.toFixed(2) + "TiB"
        };

        template.defaults.imports.prettyByte = function (data) {
            let dec = data;
            if (dec < 1024) {
                return dec.toFixed(2)+"B";
            }
            dec = dec / 1024.0;
            if (dec < 1024) {
                return dec.toFixed(2)+"KiB";
            }
            dec /= 1024.0;
            if (dec < 1024) {
                return dec.toFixed(2)+"MiB";
            }
            dec = dec / 1024.0;
            if (dec < 1024) {
                return dec.toFixed(2)+"GiB";
            }
            dec = dec / 1024;
            return dec.toFixed(2) + "TiB"
        };

        template.defaults.imports.figureCpuUsed = function (free, total) {
            return ((1000 - free) / 1000 * total).toFixed(2)
        };

        template.defaults.imports.figureCpuFree = function (free, total) {
            return (free / 1000 * total).toFixed(2)
        };

        template.defaults.imports.netmask2prefix = function (netmask) {
            if (!netmask) return undefined;
            return netmask.split('.').map(Number)
                .map(part => (part >>> 0).toString(2))
                .join('').split('1').length - 1;
        };

        template.defaults.imports.prefix2netmask = function (prefix) {
            if (!prefix) return undefined;
            let mask = [];
            for(let i = 0;i < 4; i++) {
                let n = Math.min(prefix, 8);
                mask.push(256 - Math.pow(2, 8-n));
                prefix -= n;
            }
            return mask.join('.');
        }
    }
}