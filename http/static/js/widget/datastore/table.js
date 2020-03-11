import {DataStoreApi} from "../../api/datastore.js";


export class DataStoreTable {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.tasks = props.tasks;
    }

    loading() {
        return `<tr><td colspan="8" style="text-align: center">Loading...</td></tr>`
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
        return template.compile(`
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.uuid}}"></td>
                <td>{{i+1}}</td>
                <td><a href="/download/{{v.id}}">{{v.uuid}}</a></td>
                <td>{{v.name}}</td>
                <td>{{v.source}}</td>
                <td>{{v.capacity | prettyByte}}</td>
                <td>{{v.available | prettyByte}}</td>
                <td><span class="{{v.state}}">{{v.state}}</span></td>
            </tr>
        {{/each}}
        `)(data)
    }
}