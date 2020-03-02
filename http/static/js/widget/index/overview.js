import {HyperApi} from "../../api/hyper.js";


export class Overview {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.tasks = props.tasks;
    }

    loading() {
        return (``);
    }

    refresh(func) {
        //$(this.id).html(this.loading());
        new HyperApi({tasks: this.tasks}).get(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            if (func) {
                func({data, resp: e.resp});
            }
        });
    }

    render(data) {
        return template.compile(`
            <dl class="dl-horizontal">
                <dt>Version:</dt>
                <dd>{{version.version}}</dd>
                <dt>Built on:</dt>
                <dd>{{version.date}}</dd>
                <dt>Hypervisor:</dt>
                <dd>{{hyper.name}}</dd>
                <dt>Processor:</dt>
                <dd>{{hyper.cpuNum}} | {{hyper.cpuUtils | figureCpuUsed hyper.cpuNum}} | {{hyper.cpuVendor}}</dd>
                <dt>Memory</dt>
                <dd>
                    {{hyper.memTotal | prettyByte}} | {{hyper.memFree | prettyByte}} | {{hyper.memCached | prettyByte}}
                </dd>
            </dl>
        `)(data)
    }
}