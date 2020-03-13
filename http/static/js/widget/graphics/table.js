import {GraphicsApi} from "../../api/graphics.js";


export class GraphicsTable {
    // {
    //   id: '#xx',
    //   inst: 'uuid',
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.tasks = props.tasks;
        this.inst = props.inst;
    }

    loading() {
        return `<tr><td colspan="5" style="text-align: center">Loading...</td></tr>`
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $(this.id).html(this.loading());
        new GraphicsApi({
            tasks: this.tasks,
            inst: this.inst,
        }).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        return template.compile(`
        {{each items v i}}
            <tr>
                <td><input id="on-one" type="checkbox" data="{{v.type}}"></td>
                <td>{{i+1}}</td>
                <td>{{v.type}}</td>
                <td>{{v.password}}</td>
                <td>{{v.listen}}:{{v.port != '' ? v.port : -1}}</td>
            </tr>
        {{/each}}
        `)(data);
    }
}