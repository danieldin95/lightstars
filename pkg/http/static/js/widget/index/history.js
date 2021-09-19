import {HistoryApi} from "../../api/history.js";
import {Widget} from "../widget.js";


export class History extends Widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
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
        new HistoryApi({tasks: this.tasks}).list(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            func({data, resp: e.resp});
        });
    }

    render(data) {
        if (data && data.items) {
            data.items.reverse();
        }
        return this.compile(`
        {{if (items.length === 0)}}
            <tr>
                <td colspan="5" class="text-center">{{'no data to display' | i}}</td>
            </tr>
        {{/if}}
        {{each items v i}}
            <tr>
                <td>{{v.user}}</td>
                <td>{{v.date}}</td>
                <td>{{v.client}}</a></td>
                <td>{{v.method}}</td>
                <td>{{v.url}}</td>
            </tr>
        {{/each}}
        `, data);
    }
}
