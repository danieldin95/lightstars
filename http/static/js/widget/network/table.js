import {NetworkApi} from "../../api/network.js";
import {WidgetBase} from "../base.js";


export class NetworkTable extends WidgetBase {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
    }

    loading() {
        return `<tr><td colspan="7" style="text-align: center">Loading...</td></tr>`;
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        new NetworkApi({tasks: this.tasks}).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        return this.compile(`
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.uuid}}"></td>
                <td>{{i+1}}</td>
                <td><a href="#/network/{{v.uuid}}">{{v.uuid}}</a></td>
                <td>{{v.name}}</td>
                <td>{{if v.address == ""}}
                    --
                    {{else}}
                    {{v.address}}/{{v.netmask | netmask2prefix }}{{v.prefix ? v.prefix : ''}}
                    {{/if}}
                </td>
                <td>{{v.mode != '' ? v.mode : 'isolated'}}</td>
                <td><span class="{{v.state}}">{{v.state}}</span></td>
            </tr>
        {{/each}}
        `, data);
    }
}
