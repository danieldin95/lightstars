import {NetworkApi} from "../../api/network.js";
import {Widget} from "../widget.js";
import {Location} from "../../com/location.js";


export class NetworkTable extends Widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
    }

    loading() {
        return `<tr><td colspan="6" style="text-align: center">Loading...</td></tr>`;
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
        let query = Location.query();
        return this.compile(`
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.uuid}}"></td>
                <td>{{i+1}}</td>
                <td><a id="on-this" class="text-decoration-none" data="{{v.uuid}}" href="#/network/{{v.uuid}}?${query}">{{v.uuid}}</a></td>
                <td>{{v.name}}</td>
                <td>{{v.bridge}}</td>
                <td><span class="{{v.state}}">{{v.state}}</span></td>
            </tr>
        {{/each}}
        `, data);
    }
}
