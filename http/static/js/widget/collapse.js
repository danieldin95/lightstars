import {Location} from "../com/location.js";

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
            Location.set(this.default);
        }
        let page = Location.get(this.default);
        for (let i = 0;  i < this.pages.length; i++) {
            let v = this.pages[i];
            if (page == v.name) {
                $(v.id).fadeIn('slow');
                $(v.id).collapse();
            }
            $(v.id).on('show.bs.collapse', this, function (e) {
                Location.set(v.name);
                $(this).fadeIn('slow');
            });
            $(v.id).on('hide.bs.collapse', this, function (e) {
                $(this).fadeOut('fast');
            });
        }
    }
}