import {Location} from "../com/location.js";

export class Collapse {
    // {
    //   default: page,
    //   pages: [ {id: '#', name: ""} ],
    //   force: false,
    //   update: false, // update href.
    // }
    constructor(props) {
        // collapse
        this.collapse = props;
        this.pages = props.pages;
        this.default = props.default;

        let force = props.force !== undefined ? props.force : false;
        let update = props.update !== undefined ? props.update : true;
        let page = this.default;
        console.log('Collapse.update', update);
        if (update) {
            if (force) {
                Location.set(this.default);
            }
            page = Location.get(this.default);
        }
        for (let v of this.pages) {
            if (page === v.name) {
                $(v.id).fadeIn('slow');
                $(v.id).collapse();
            }
            $(v.id).on('show.bs.collapse', this, function (e) {
                if (update) {
                    Location.set(v.name);
                }
                $(this).fadeIn('slow');
            });
            $(v.id).on('hide.bs.collapse', this, function (e) {
                $(this).fadeOut('fast');
            });
        }
    }
}