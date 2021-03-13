import {Widget} from "../widget.js";
import {SnapshotApi} from "../../api/snapshot.js";


export class SnapshotTable extends Widget {
    // {
    //   id: '#xx',
    //   inst: 'uuid',
    // }
    constructor(props) {
        super(props);
        this.inst = props.inst;
    }

    loading() {
        return `<tr><td colspan="5" style="text-align: center">Loading...</td></tr>`;
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        new SnapshotApi({
            tasks: this.tasks,
            inst: this.inst,
        }).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        return this.compile(`
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.name}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.name}}</td>
                <td>{{v.uptime | prettyTime}}</td>
                <td>{{v.state}}</td>
            </tr>
        {{/each}}
        `, data);
    }
}
