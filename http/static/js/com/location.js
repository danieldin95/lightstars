
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
        let path = this.href();
        let uri = path.split('#', 2)[0];

        console.log('Collapse.set', path, name);

        name = name || "";
        this.href(uri + '#' + name);
        for (let i = 0; i < this.on.length; i++) {
            let on = this.on[i];
            if (on && on.func) {
                on.func({data: on.data, name});
            }
        }
    }

    static get (name) {
        let path = this.href();
        let pages = path.split('#', 2);

        name = name || "";
        return (pages.length === 2 && pages[1] !== "") ? pages[1] : name;
    }
}