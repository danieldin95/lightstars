
export class Location {
    // {
    //   data: '',
    //   func: function (e) {}
    // }
    static listen = [];

    static set(name) {
        console.log('Collapse.set', name);
        let queries = window.location.href.split("?", 2);
        let path = queries[0];
        let query = queries.length == 2 ? queries[1] : "";
        let uri = path.split('#', 2)[0];

        name = name || "";
        if (query == "") {
            window.location.href = uri+'#'+name;
        } else {
            window.location.href = uri+'#'+name+"?"+query;
        }
        for (let i = 0; i < this.listen.length; i++) {
            let l = this.listen[i];
            if (l && l.func) {
                l.func({data: l.data, name});
            }
        }
    }

    static get (name) {
        let path = window.location.href.split("?", 2)[0];
        let pages = path.split('#', 2);

        name = name || "";
        return (pages.length === 2 && pages[1] !== "") ? pages[1] : name;
    }
}