import {Widget} from "../widget.js";
import {PortApi} from "../../api/port.js";
import {Location} from "../../lib/location.js";


export class PortTable extends Widget {
    // {
    //   id: '#xx',
    //   uuid: 'uuid',
    // }
    constructor(props) {
        super(props);
        console.log(props);
        this.uuid = props.uuid;
        this.name = props.name;
        this.bridge = props.bridge;
    }

    loading() {
        return `<tr><td colspan="7" class="text-center">Loading...</td></tr>`;
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        new PortApi({
            tasks: this.tasks,
            bridge: this.bridge,
            query: {uuid: this.uuid}
        }).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        let query = Location.query();
        return this.compile(`
        {{if (items.length === 0)}}
            <tr>
                <td colspan="7" class="text-center">{{'no data to display' | i}}</td>
            </tr>
        {{/if}}
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.domain.uuid}},{{v.address}}"></td>
                <td>{{i+1}}</td>
                <td><a id="on-this" class="text-decoration-none" data="{{v.domain.uuid}}" 
                        href="#/guest/{{v.domain.uuid}}?${query}">{{v.domain.name}}</a>
                </td>
                <td>{{v.device == "" ? '-' : v.device}}</td>
                <td>{{v.address}}</td>
                <td>{{v.ipaddr == "" ? '-' : v.ipaddr}}</td>
                <td>{{v.model}}</td>
            </tr>
        {{/each}}
        `, data);
    }
}
