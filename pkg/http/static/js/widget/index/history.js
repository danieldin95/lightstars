import {HistoryApi} from "../../api/history.js";
import {Widget} from "../widget.js";


export class History extends Widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
        this.filter = props.filter || null;
        this.search = "";
        this.pageSize = props.pageSize || 10;
        this.page = 1;
        this.items = [];
        this.pager = props.pager || null;
        window.addEventListener("lightstar:history:append", () => this.refresh());
        $(this.id).off("history:refresh").on("history:refresh", () => this.refresh());
        if (this.pager) {
            $(this.pager.prev).off("click").on("click", () => this.prevPage());
            $(this.pager.next).off("click").on("click", () => this.nextPage());
        }
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
            let items = (e.resp && e.resp.items) ? e.resp.items : [];
            if (typeof e.data.filter === "function") {
                items = items.filter((v) => e.data.filter(v));
            }
            if (e.data.search) {
                let q = e.data.search.toLowerCase();
                items = items.filter((v) => {
                    let text = [
                        v.user,
                        v.date,
                        v.client,
                        v.result,
                        v.method,
                        v.url,
                    ].map((x) => (x || "").toString().toLowerCase()).join(" ");
                    return text.indexOf(q) >= 0;
                });
            }
            items = items.map((v) => {
                let item = {...v};
                let method = (item.method || "").toUpperCase();
                let result = (item.result || "").trim();
                if ((method === "POST" || method === "PUT" || method === "DELETE") &&
                    (result === "" || result === "-")) {
                    item.result = "success";
                }
                return item;
            });
            items = items.slice().sort((a, b) => {
                let ta = Date.parse(a.date || "") || 0;
                let tb = Date.parse(b.date || "") || 0;
                return tb - ta;
            });
            e.data.items = items;
            e.data.page = 1;
            e.data.renderPage();
            if (func) {
                func({data, resp: {...e.resp, items}});
            }
        });
    }

    setSearch(q) {
        this.search = (q || "").trim();
        this.page = 1;
        this.refresh();
    }

    totalPages() {
        return Math.max(1, Math.ceil(this.items.length / this.pageSize));
    }

    pageItems() {
        let start = (this.page - 1) * this.pageSize;
        return this.items.slice(start, start + this.pageSize);
    }

    renderPage() {
        let pageTotal = this.totalPages();
        if (this.page > pageTotal) {
            this.page = pageTotal;
        }
        let data = {items: this.pageItems()};
        $(this.id).html(this.render(data));
        this.updatePager();
    }

    updatePager() {
        if (!this.pager) {
            return;
        }
        let pageTotal = this.totalPages();
        $(this.pager.info).text(`${this.page}/${pageTotal}`);
        $(this.pager.prev).prop("disabled", this.page <= 1);
        $(this.pager.next).prop("disabled", this.page >= pageTotal);
    }

    prevPage() {
        if (this.page > 1) {
            this.page -= 1;
            this.renderPage();
        }
    }

    nextPage() {
        if (this.page < this.totalPages()) {
            this.page += 1;
            this.renderPage();
        }
    }

    toCSVValue(v) {
        let s = v == null ? "" : String(v);
        return `"${s.replace(/"/g, '""')}"`;
    }

    downloadCSV() {
        new HistoryApi({tasks: this.tasks}).list(this, (e) => {
            let items = (e.resp && e.resp.items) ? e.resp.items.slice() : [];
            if (typeof this.filter === "function") {
                items = items.filter((v) => this.filter(v));
            }
            items.sort((a, b) => {
                let ta = Date.parse(a.date || "") || 0;
                let tb = Date.parse(b.date || "") || 0;
                return tb - ta;
            });
            let header = ["User", "Date", "Client", "Message", "URL"];
            let lines = [header.map((v) => this.toCSVValue(v)).join(",")];
            items.forEach((v) => {
                let urlDisplay = `${v.method || ""} ${v.url || ""}`.trim() || "-";
                let row = [
                    v.user || "",
                    v.date || "",
                    v.client || "",
                    v.result || "-",
                    urlDisplay,
                ];
                lines.push(row.map((x) => this.toCSVValue(x)).join(","));
            });
            let csv = "\uFEFF" + lines.join("\n");
            let blob = new Blob([csv], {type: "text/csv;charset=utf-8;"});
            let ts = new Date().toISOString().replace(/[:.]/g, "-");
            let filename = `history-${ts}.csv`;
            let a = document.createElement("a");
            a.href = URL.createObjectURL(blob);
            a.download = filename;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(a.href);
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
                <td>{{v.user}}</td>
                <td class="history-date">{{v.date}}</td>
                <td>{{v.client}}</a></td>
                <td class="history-result" title="{{v.result || '-'}}">{{v.result || '-'}}</td>
                <td class="history-url" title="{{((v.method || '') + ' ' + (v.url || '')).trim() || '-'}}">{{((v.method || '') + ' ' + (v.url || '')).trim() || '-'}}</td>
            </tr>
        {{/each}}
        `, data);
    }
}
