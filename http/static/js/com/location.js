
export class Location {
    // {
    //   data: '',
    //   func: function (e) {}
    // }
    static listen(on) {
        if (!this.on) {
            this.on = [];
        }
        this.on.push(on);
    }

    static href(value) {
        if (value) {
            window.location.href = value;
        }
        return window.location.href;
    }

    static set(name) {
        name = name || "";
        let path = this.href();
        let uri = path.split('#', 2)[0];
        let query = this.query();

        console.log('Location.set', path, name);
        this.href(uri + '#' + name + "?" + query);
        for (let i = 0; i < this.on.length; i++) {
            let on = this.on[i];
            if (on && on.func) {
                on.func({data: on.data, name});
            }
        }
        console.log('Location.set', this.query('node'));
    }

    static get (name) {
        name = name || "";
        let path = this.href();
        let page = name;
        if (path.indexOf('#') > 0) {
            page = path.split('#', 2)[1];
        }
        return page.split('?')[0];
    }

    static url() {
        let path = this.href();
        return path.split('#')[0];
    }

    static query(name, value) {
        let url = this.url();
        let page = this.get();

        console.log("Location.query", name, value);
        if (name && typeof value !== 'undefined') {
            this.href(url + "#" + page + "?" + name + "=" + value);
        }
        let path = this.href();
        page = "";
        if (path.indexOf('#') > 0) {
            page = path.split('#', 2)[1];
        }
        let query = "";
        if (page.indexOf('?') > 0) {
            query = page.split("?", 2)[1];
        }
        if (name && typeof value === 'undefined') {
            let values = query.split('&');
            for (let i = 0; i < values.length; i++) {
                let ele = values[i];
                if (ele.indexOf('=') <= 0) {
                    continue
                }
                let [k, v] = ele.split('=', 2);
                if (k == name) {
                    return v;
                }
            }
        }
        return query;
    }
}