import {Utils} from "./utils.js";
import {I18N} from "./i18n.js";


export class Template {
    constructor() {
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
        m.i = function (key, param1) {
            return I18N.i18n(key, param1);
        };
        m.fullTime = function(second) {
            let time = [];
            let min = Math.floor(second / 60);
            time.push((second % 60) + "s");
            if (min < 60) {
                time.push(min + "m");
                return time.reverse().join("");
            }
            let hour = Math.floor(min / 60);
            time.push((min % 60) + "m");
            if (hour < 24) {
                time.push(hour + "h");
                return time.reverse().join("");
            }
            time.push((hour % 24) + "h");
            time.push(Math.floor(hour / 24) + "d");
            return time.reverse().join("");
        };
        m.prettyTime = function(second) {
            let min = Math.floor(second / 60);
            let sec = (second % 60);
            if (min < 60) {
                return min + "m" + sec + "s";
            }
            let hour = Math.floor(min / 60);
            min = (min % 60);
            if (hour < 24) {
                return hour + "h" + min + "m";
            }
            let day = Math.floor(hour / 24);
            hour = (hour % 24);
            return day + "d" + hour + "h";
        };
        m.prettyCpuMode = function (mode) {
          if (mode === "custom") {
              return "custom";
          } else if (mode === "host-passthrough") {
              return "passthrough";
          } else {
              return "default";
          }
        };
        return m;
    }
}
