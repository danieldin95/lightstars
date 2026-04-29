import {Widget} from "../widget.js";
import {HyperApi} from "../../api/hyperapi.js";


export class SystemWid extends Widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
        console.log(props);
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        new HyperApi({tasks: this.tasks}).get(this,function (e) {
            let resp = e.resp || {};
            let hyper = resp.hyper || {};
            let toBytes = (v) => {
                if (v === null || v === undefined) {
                    return NaN;
                }
                if (typeof v === "number") {
                    return v;
                }
                let s = String(v).trim();
                if (s === "") {
                    return NaN;
                }
                let m = s.match(/^([0-9]+(?:\.[0-9]+)?)\s*([a-zA-Z]+)?$/);
                if (!m) {
                    return Number(s);
                }
                let n = Number(m[1]);
                let u = (m[2] || "").toLowerCase();
                if (u === "" || u === "b") return n;
                if (u === "kb" || u === "kib" || u === "k") return n * 1024;
                if (u === "mb" || u === "mib" || u === "m") return n * 1024 * 1024;
                if (u === "gb" || u === "gib" || u === "g") return n * 1024 * 1024 * 1024;
                if (u === "tb" || u === "tib" || u === "t") return n * 1024 * 1024 * 1024 * 1024;
                return n;
            };
            let pickNum = (...vals) => {
                for (let v of vals) {
                    let n = toBytes(v);
                    if (!Number.isNaN(n)) {
                        return n;
                    }
                }
                return 0;
            };
            hyper.memTotal = pickNum(hyper.memTotal, hyper.memoryTotal, hyper.memory, hyper.mem);
            hyper.memFree = pickNum(hyper.memFree, hyper.memoryFree, hyper.freeMem, hyper.memAvailable);
            hyper.memCached = pickNum(hyper.memCached, hyper.memoryCached, hyper.cachedMem);
            resp.hyper = hyper;
            $(e.data.id).html(e.data.render(resp));
            if (func) {
                func({data, resp});
            }
        });
    }

    render(data) {
        return this.compile(`
            <div class="dashboard-stats system-dashboard">
                <div class="dashboard-grid dashboard-grid-2">
                    <div class="dashboard-stat total">
                        <div class="label">{{'uptime' | i}}</div>
                        <div class="value">{{hyper.uptime | prettyTime}}</div>
                    </div>
                    <div class="dashboard-stat total">
                        <div class="label">{{'version' | i}}</div>
                        <div class="value">{{version.version}}</div>
                    </div>
                </div>
                <div class="dashboard-grid dashboard-grid-2">
                    <div class="dashboard-stat total">
                        <div class="label">{{'built on' | i}}</div>
                        <div class="value">{{version.date}}</div>
                    </div>
                    <div class="dashboard-stat total">
                        <div class="label">{{'hypervisor' | i}}</div>
                        <div class="value">{{hyper.name}}</div>
                    </div>
                </div>
                <div class="dashboard-grid dashboard-grid-1">
                    <div class="dashboard-stat up">
                        <div class="label">{{'processor' | i}}</div>
                        <div class="value">
                            {{hyper.cpuUtils | figureCpuFree hyper.cpuNum}} / {{hyper.cpuNum}}
                        </div>
                        <div class="label">{{hyper.cpuVendor}}</div>
                    </div>
                </div>
                <div class="dashboard-grid dashboard-grid-1">
                    <div class="dashboard-stat total">
                        <div class="label">{{'memory' | i}}</div>
                        <div class="value">
                            {{hyper.memFree | prettyByte}} / {{hyper.memTotal | prettyByte}}
                        </div>
                        <div class="label">cache {{hyper.memCached | prettyByte}}</div>
                    </div>
                </div>
            </div>`, data);
    }
}
