import {Widget} from "../widget.js";
import {HyperApi} from "../../api/hyper.js";


export class Overview extends Widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
    }

    loading() {
        return (``);
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
            <dl class="dl-horizontal">
                <dt>{{'version' | i}}:</dt>
                <dd>{{version.version}}</dd>
                <dt>{{'built on' | i}}:</dt>
                <dd>{{version.date}}</dd>
                <dt>{{'hypervisor' | i}}:</dt>
                <dd>{{hyper.name}}</dd>
                <dt>{{'processor' | i}}:</dt>
                <dd title="{{'total|free|vendor' | i}}">
                    {{hyper.cpuNum}} | {{hyper.cpuUtils | figureCpuFree hyper.cpuNum}} | {{hyper.cpuVendor}}
                </dd>
                <dt>{{'memory' | i}}</dt>
                <dd title="{{'total|free|cache' | i}}">
                    {{hyper.memTotal | prettyByte}} | {{hyper.memFree | prettyByte}} | {{hyper.memCached | prettyByte}}
                </dd>
            </dl>
        `, data);
    }
}
