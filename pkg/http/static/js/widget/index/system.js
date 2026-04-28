import {Widget} from "../widget.js";
import {HyperApi} from "../../api/hyper.js";


export class System extends Widget {
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
            $(e.data.id).html(e.data.render(e.resp));
            if (func) {
                func({data, resp: e.resp});
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
                            {{hyper.cpuNum}} / {{hyper.cpuUtils | figureCpuFree hyper.cpuNum}}
                        </div>
                        <div class="label">{{hyper.cpuVendor}}</div>
                    </div>
                </div>
                <div class="dashboard-grid dashboard-grid-1">
                    <div class="dashboard-stat total">
                        <div class="label">{{'memory' | i}}</div>
                        <div class="value">
                            {{hyper.memTotal | prettyByte}} / {{hyper.memFree | prettyByte}}
                        </div>
                        <div class="label">cache {{hyper.memCached | prettyByte}}</div>
                    </div>
                </div>
            </div>`, data);
    }
}
