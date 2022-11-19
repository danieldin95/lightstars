import {Widget} from "../widget.js";
import {LeaseApi} from "../../api/lease.js";


export class LeaseTable extends Widget {
    // {
    //   id: '#xx',
    //   uuid: 'uuid',
    // }
    constructor(props) {
        super(props);
        this.uuid = props.uuid;
    }

    loading() {
        return `<tr><td colspan="5" class="text-center">Loading...</td></tr>`;
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
        {{if (items.length === 0)}}
            <tr>
                <td colspan="5" class="text-center">{{'no data to display' | i}}</td>
            </tr>
        {{/if}}
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.type}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.hostname}}</td>
                <td>{{v.mac}}</td>
                <td>{{v.ipAddr}}/{{v.prefix}}</td>
            </tr>
        {{/each}}
        `, data);
    }
}
