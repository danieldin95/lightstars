import {WidgetBase} from "../base.js";
import {LeaseApi} from "../../api/lease.js";


export class LeaseTable extends WidgetBase {
    // {
    //   id: '#xx',
    //   uuid: 'uuid',
    // }
    constructor(props) {
        super(props);
        this.uuid = props.uuid;
        console.log("LeaseTable", props);
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
        new LeaseApi({
            tasks: this.tasks,
            net: this.uuid,
        }).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        return this.compile(`
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.type}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.mac}}</td>
                <td>{{v.ipAddr}}/{{v.prefix}}</td>
                <td>{{v.hostname}}</td>
            </tr>
        {{/each}}
        `, data);
    }
}
