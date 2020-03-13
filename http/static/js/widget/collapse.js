
export class Collapse {
    // {
    //   default: page,
    //   pages: [ {id: '#', name: ""} ],
    //   force: false,
    // }
    constructor(props) {
        // collapse
        this.collapse = props;
        this.pages = props.pages;
        this.force = props.force ? props.force : false;
        this.default = props.default;

        if (this.force) {
            this.set(this.default);
        }
        let page = this.get(this.default);
        for (let i = 0;  i < this.pages.length; i++) {
            let v = this.pages[i];
            if (page == v.name) {
                $(v.id).fadeIn('slow');
                $(v.id).collapse();
            }
            $(v.id).on('show.bs.collapse', this, function (e) {
                e.data.set(v.name);
                $(this).fadeIn('slow');
            });
            $(v.id).on('hide.bs.collapse', this, function (e) {
                $(this).fadeOut('fast');
            });
        }
    }

    set(name) {
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
    }

    get (name) {
        let path = window.location.href.split("?", 2)[0];
        let pages = path.split('#', 2);

        name = name || "";
        return (pages.length === 2 && pages[1] !== "") ? pages[1] : name;
    }
}