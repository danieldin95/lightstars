import {Widget} from "../widget.js";
import {VolumeApi} from "../../api/volume.js";
import {CheckBox} from "../checkbox/checkbox.js";
import {Location} from "../../com/location.js";
import {Utils} from "../../com/utils.js";


export default class VolumeTable extends Widget {

    constructor(props) {
        super(props);
        this.checkbox = new CheckBox(props);
        this.pool = props.pool;
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
        new VolumeApi({
            pool: this.pool
        }).list(this, function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        })
    }

    formatData(data) {
        let items = data.items;
        return Object.assign({}, data, {
            items: items.map((i) => {
                let name = "";
                if (i.type === "dir") {
                    name = '.'+Utils.basename(i.name);
                } else {
                    name = Utils.basename(i.name);
                }
                return Object.assign({}, i, {name});
            })
        })
    }

    render(data) {
        let query = Location.query();
        let prefix = window.location.pathname;

        return this.compile(`
            {{each items v i}}
                <tr class="sortable">
                    <td><input id="on-one" type="checkbox" data="{{v.name}}"></td>
                    <td>
                        {{if v.type == "dir"}}
                        <img src="/static/images/folder-icon.svg" style="opacity: 0.6; filter:alpha(opacity=60);"/>
                        {{else if v.type == "file"}}
                        <img src="/static/images/file-icon.svg" style="opacity: 0.6; filter:alpha(opacity=60);"/>
                        {{/if}}
                    </td>
                    <td>
                      {{if v.type == "dir"}}
                      <a id="on-this" data-name="{{v.name}}" data-type="{{v.type}}" href="javascript:void(0)">{{v.name}}</a>
                      {{else if v.type == "file"}}
                      <a href="/api/datastore/${this.pool}/volume/{{v.name}}">{{v.name}}</a>
                      {{else}}
                      <a data-name="{{v.name}}" href="#">{{v.name}}</a>
                      {{/if}}
                    </td>
                    <td>{{if v.type == "dir"}} - {{else}} {{v.capacity | prettyByte}} {{/if}}</td>
                    <td>{{if v.type == "dir"}} - {{else}} {{v.allocation | prettyByte}} {{/if}}</td>
                </tr>
            {{/each}}
            `, this.formatData(data));
    }
}
