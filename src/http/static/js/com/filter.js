import {Utils} from "./utils.js";

const art = template.defaults.imports;

export class Filters {
    constructor() {
        art.aton = function (data, n) {
            return Utils.iton(data, n);
        };

        art.prettyKiB = function (data, fra) {
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

        art.prettyByte = function (data, fra) {
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

        art.figureCpuUsed = function (free, total) {
            return ((1000 - free) / 1000 * total).toFixed(2)
        };

        art.figureCpuFree = function (free, total) {
            return (free / 1000 * total).toFixed(2)
        };

        art.netmask2prefix = function (netmask) {
            if (!netmask) return undefined;
            return netmask.split('.').map(Number)
                .map(part => (part >>> 0).toString(2))
                .join('').split('1').length - 1;
        };

        art.prefix2netmask = function (prefix) {
            if (!prefix) return undefined;
            let mask = [];
            for(let i = 0;i < 4; i++) {
                let n = Math.min(prefix, 8);
                mask.push(256 - Math.pow(2, 8-n));
                prefix -= n;
            }
            return mask.join('.');
        };

        art.vncPassword = function (inst) {
            return Utils.graphic(inst, 'vnc', 'password')
        };

        art.spicePassword = function (inst) {
            return Utils.graphic(inst, 'spice', 'password')
        };
    }
}
