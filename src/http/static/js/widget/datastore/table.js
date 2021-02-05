import {Widget} from "../widget.js";
import {Api} from "../../api/api.js";
import {DataStoreApi} from "../../api/datastores.js";
import {Location} from "../../lib/location.js";

export class DataStoreTable extends Widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
    }

    loading() {
        return `<tr><td colspan="8" style="text-align: center">Loading...</td></tr>`;
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        new DataStoreApi({tasks: this.tasks}).list(this,function (e) {
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
                <td><a id="on-this" class="text-decoration-none" data="{{v.uuid}}" href="#/datastore/{{v.uuid}}?${query}">{{v.name}}</a></td>
                <td>{{v.source}}</td>
                <td>{{v.capacity | prettyByte}}</td>
                <td>{{v.allocation | prettyByte}}</td>
                <td><span class="{{v.state}}">{{v.state}}</span></td>
            </tr>
        {{/each}}
        `, data);
    }
}
