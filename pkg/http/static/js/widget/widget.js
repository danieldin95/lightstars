import {Location} from "../lib/location.js";

export class Widget {
    constructor(props) {
        this.id = props.id || "";
        this.props = props;
        this.tasks = props.tasks || "tasks";
    }

    compile(tmpl, data) {
        return template.compile(tmpl)(data);
    }

    loading() {
        return 'loading...';
    }

    title(name) {
        if (!this._alias) {
            this._alias = "LightStar"
        }
        $(document).attr("title", `${name} - ${this._alias}`);
    }

    url(page) {
        let query = Location.query();
        return page + "?" + query
    }
}
