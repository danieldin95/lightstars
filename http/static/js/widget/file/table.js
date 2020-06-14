import {Widget} from "../widget.js";
import {FileApi} from "../../api/file.js";
import {CheckBoxTab} from "../checkbox/checkbox.js";

class CheckBox extends CheckBoxTab {
}

export default class FileTable extends Widget {

    constructor(props) {
        super(props);
        this.api = new FileApi()
        this.checkbox = new CheckBox(props)

        this.refresh((e) => {
            this.checkbox.refresh();
        })
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

        this.api.list(this, function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        })

    }

    render(data) {
        let prefix = window.location.pathname
        return this.compile(`
            <div class="card-body">
                <table class="table table-hover table-striped">
                    <thead>
                        <tr class="sortable">
                            <th><input type="checkbox" value=""></th>
                            <th>Type</th>
                             <th>Name</th>
                            <th>Size</th>
                            <th>Time</th>
                        </tr>
                    </thead>
                    <tbody>
                    {{each items v i}}
                        <tr class="sortable">
                            <td><input id="on-one" type="checkbox" data="{{v.uuid}}"></td>
                            <td >
                                
<!--                                <svg class="{{v.IsDir ? 'bi bi-folder-fill' : 'bi bi-file-earmark'}}" width="1em" height="1em" viewBox="0 0 16 16"></svg>-->
                                {{if v.IsDir}}
                                <svg class="bi bi-folder-fill" width="1em" height="1em" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                                    <path fill-rule="evenodd" d="M9.828 3h3.982a2 2 0 0 1 1.992 2.181l-.637 7A2 2 0 0 1 13.174 14H2.826a2 2 0 0 1-1.991-1.819l-.637-7a1.99 1.99 0 0 1 .342-1.31L.5 3a2 2 0 0 1 2-2h3.672a2 2 0 0 1 1.414.586l.828.828A2 2 0 0 0 9.828 3zm-8.322.12C1.72 3.042 1.95 3 2.19 3h5.396l-.707-.707A1 1 0 0 0 6.172 2H2.5a1 1 0 0 0-1 .981l.006.139z"/>
                                </svg>
                                {{else}}
                                <svg class="bi bi-file-earmark" width="1em" height="1em" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M4 1h5v1H4a1 1 0 0 0-1 1v10a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V6h1v7a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2z"/>
                                    <path d="M9 4.5V1l5 5h-3.5A1.5 1.5 0 0 1 9 4.5z"/>
                                </svg>
                                {{/if}}
                                
                            </td>
                            <td><a href="${prefix}{{v.FileName}}">{{v.FileName}}</a></td>
                            <td>{{v.Size}}</td>
                            <td>{{v.ModTime}}</td>
                        </tr>
                    {{/each}}
                    </tbody>
                </table>
            </div>`, data);
    }
}
