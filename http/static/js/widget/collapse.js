
export class Collapse {
    // {
    //   default: page,
    //   pages: [ {id: '#', name: ""} ],
    // }
    constructor(props) {
        // collapse
        this.collapse = props;

        let page = this.get(this.collapse.default);
        for (let i = 0;  i < this.collapse.pages.length; i++){
            let v = this.collapse.pages[i];
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
        let querier = window.location.href.split("?", 2);
        let path = querier[0];
        let query = querier.length == 2 ? querier[1] : "";
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
        return (pages.length == 2 && pages[1] != "") ? pages[1] : name;
    }
}