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
    }
}